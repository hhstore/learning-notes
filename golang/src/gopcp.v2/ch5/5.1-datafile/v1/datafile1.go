package v1

import (
	"errors"
	"io"
	"os"
	"sync"
)

//***************************************************************
// 数据类型
type Data []byte

// 数据文件的接口声明:
type DataFile interface {
	Read() (rsn int64, d Data, err error)
	Write(d Data) (wsn int64, err error)

	RSN() int64
	WSN() int64

	DataLen() uint32
	Close() error
}

//***************************************************************
//                    数据文件的接口实现
// 说明:
//		- 基于 互斥锁 + 读写锁 实现
//
//***************************************************************

type myDataFile struct {
	f       *os.File
	fmutex  sync.RWMutex // 读写锁: 文件锁
	woffset int64        // 偏移量: 写入
	roffset int64        // 偏移量: 读入
	wmutex  sync.Mutex   // 互斥锁: 写入锁
	rmutex  sync.Mutex   // 互斥锁: 读入锁

	dataLen uint32
}

//***************************************************************

// 创建数据文件:
func NewDataFile(path string, dataLen uint32) (DataFile, error) {
	f, err := os.Create(path) // 创建文件
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
	return df, nil

}

// 读取:
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

	//
	for {
		df.fmutex.RLock() // 加锁
		{
			_, err = df.f.ReadAt(bytes, offset)
			if err != nil {
				if err == io.EOF {
					df.fmutex.RUnlock()
					continue
				}
				//
				df.fmutex.RUnlock()
				return
			}
			//
			d = bytes

		}
		df.fmutex.RUnlock() // 解锁
		return
	}

}

// 写入:
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

	_, err = df.f.Write(bytes)
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

func (df *myDataFile) Close() error {
	if df.f == nil {
		return nil
	}

	return df.f.Close()
}
