package v3

import (
	"errors"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

//***************************************************************

type Data []byte

// 数据文件的接口声明:
type DataFile interface {
	Read() (rsn int64, d Data, err error) // 读取一个数据块
	Write(d Data) (wsn int64, err error)  // 写入一个数据块
	RSN() int64                           // 获取最后读取的数据块的序列号
	WSN() int64                           // 获取最后写入的数据块的序列号
	DataLen() uint32                      // 获取数据块长度
	Close() error                         // 关闭数据文件
}

//***************************************************************
//                    数据文件的接口实现
// 说明:
//		- 基于 读写锁 + 条件变量 + 原子操作实现
//
//***************************************************************

// 数据文件的接口实现:
type myDataFile struct {
	f       *os.File     // 文件
	fmutex  sync.RWMutex // 读写锁: 针对文件
	rcond   *sync.Cond   // 条件变量: 针对读操作 [关键]
	woffset int64        // 偏移量: 针对写操作
	roffset int64        // 偏移量: 针对读操作
	dataLen uint32       // 数据块长度
}

//***************************************************************

func NewDataFile(path string, dataLen uint32) (DataFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	if dataLen == 0 {
		return nil, errors.New("Invalid data length!")
	}

	//
	df := &myDataFile{
		f:       f,
		dataLen: dataLen,
	}
	//
	df.rcond = sync.NewCond(df.fmutex.RLocker())

	return df, nil
}

//
// 基于原子操作实现:
//
func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	var offset int64

	//
	// 使用原子操作:
	//
	for {
		offset = atomic.LoadInt64(&df.roffset) // 原子操作:

		if atomic.CompareAndSwapInt64(&df.roffset, offset, (offset + int64(df.dataLen))) {
			break
		}
	}

	//
	rsn = offset / int64(df.dataLen)
	bytes := make([]byte, df.dataLen)

	//
	df.fmutex.RLock()
	defer df.fmutex.RUnlock()

	//
	// 使用条件变量:
	//
	for {
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			if err == io.EOF {
				// 条件变量
				df.rcond.Wait()
				continue
			}
			return
		}
		//
		d = bytes
		return
	}

}

//
// 基于原子操作实现:
//
func (df *myDataFile) Write(d Data) (wsn int64, err error) {
	var offset int64

	//
	// 原子操作:
	//
	for {
		offset = atomic.LoadInt64(&df.woffset)

		if atomic.CompareAndSwapInt64(&df.woffset, offset, (offset + int64(df.dataLen))) {
			break
		}
	}

	//
	wsn = offset / int64(df.dataLen)
	var bytes []byte

	if len(d) > int(df.dataLen) {
		bytes = d[0:df.dataLen]
	} else {
		bytes = d
	}

	//
	df.fmutex.Lock()
	defer df.fmutex.Unlock()

	//
	_, err = df.f.Write(bytes)
	//
	// 条件变量:
	//
	df.rcond.Signal()
	return
}

//
// 基于原子操作实现:
//
func (df *myDataFile) RSN() int64 {
	offset := atomic.LoadInt64(&df.roffset)

	return offset / int64(df.dataLen)
}

//
// 基于原子操作实现:
//
func (df *myDataFile) WSN() int64 {
	offset := atomic.LoadInt64(&df.woffset)

	return offset / int64(df.dataLen)
}

func (df *myDataFile) DataLen() uint32 {
	return df.dataLen
}

func (df *myDataFile) Close() error {
	if df.f == nil {
		return nil
	}
	return df.f.Close()
}
