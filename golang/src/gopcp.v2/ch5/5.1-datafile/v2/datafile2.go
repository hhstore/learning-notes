package v2

import (
	"errors"
	"io"
	"os"
	"sync"
)

//***************************************************************

type Data []byte // 数据类型

// 数据文件的接口声明:
type DataFile interface {
	Read() (rsn int64, d Data, err error) // 读取一个数据块
	Write(d Data) (wsn int64, err error)  // 写入一个数据块
	RSN() int64                           // 获取最后读取的数据块的序列号
	WSN() int64                           // 获取最后写入的数据块的序列号

	DataLen() uint32 // 获取数据块长度
	Close() error    // 关闭数据文件
}

//***************************************************************
//                    数据文件的接口实现
// 说明:
//		- 基于 互斥锁 + 读写锁 + 条件变量实现
//
//***************************************************************

// 数据文件的接口实现:
type myDataFile struct {
	f       *os.File     // 文件
	fmutex  sync.RWMutex // 读写锁: 针对文件
	rcond   *sync.Cond   // 条件变量: 针对读操作 [关键]
	woffset int64        // 偏移量: 针对写操作
	roffset int64        // 偏移量: 针对读操作
	wmutex  sync.Mutex   // 互斥锁: 写操作
	rmutex  sync.Mutex   // 互斥锁: 读操作
	dataLen uint32       // 数据块长度
}

//***************************************************************

//
func NewDataFile(path string, dataLen uint32) (DataFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	if dataLen == 0 {
		return nil, errors.New("Invalid data length!")
	}

	df := &myDataFile{
		f:       f,
		dataLen: dataLen,
	}

	//
	// 注意:
	//	- 创建条件变量
	//
	df.rcond = sync.NewCond(df.fmutex.RLocker())
	return df, nil
}

//
//  读操作:
//		- 条件变量的应用
//
func (df *myDataFile) Read() (rsn int64, d Data, err error) {
	var offset int64

	df.rmutex.Lock() // 加锁
	{
		offset = df.roffset
		df.roffset += int64(df.dataLen)
	}
	df.rmutex.Unlock()

	rsn = offset / int64(df.dataLen)
	bytes := make([]byte, df.dataLen)

	df.fmutex.RLock()         // 加锁
	defer df.fmutex.RUnlock() // 解锁

	// 注意:
	// 	- 条件变量的应用
	//
	for {
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			if err == io.EOF {
				df.rcond.Wait() // 条件变量
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
// 写操作:
//		- 条件变量的应用
//
func (df *myDataFile) Write(d Data) (wsn int64, err error) {
	var offset int64

	df.wmutex.Lock() // 加锁
	{
		offset = df.woffset
		df.woffset += int64(df.dataLen)
	}
	df.wmutex.Unlock()

	//
	wsn = offset / int64(df.dataLen)
	var bytes []byte

	//
	if len(d) > int(df.dataLen) {
		bytes = d[0:df.dataLen]
	} else {
		bytes = d
	}

	df.fmutex.Lock()         // 加锁
	defer df.fmutex.Unlock() // 解锁

	_, err = df.f.Write(bytes) // 文件写
	//注意:
	// 	- 条件变量的应用
	//
	df.rcond.Signal() // 条件变量, 发出信号
	return
}

//
func (df *myDataFile) RSN() int64 {
	df.rmutex.Lock()
	defer df.rmutex.Unlock()

	return df.roffset / int64(df.dataLen)
}

//
func (df *myDataFile) WSN() int64 {
	df.wmutex.Lock()
	defer df.wmutex.Unlock()

	return df.woffset / int64(df.dataLen)
}

//
func (df *myDataFile) DataLen() uint32 {
	return df.dataLen
}

//
func (df *myDataFile) Close() error {
	if df.f == nil {
		return nil
	}
	return df.f.Close()
}
