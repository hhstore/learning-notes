package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

/****************************************************
	  优化版 v4 - 基于 interface 实现函数多态.

说明:
	- 支持多态的函数的形参:
		- 要使用 接口类型作参数, 而非真实类型
		- 基于自定义 interface{}, 实现代码复用
	- 注意, 与利用反射实现的版本, 差异.
		- 相比反射, 更简单.
		- 反射实际上是运行时获取变量类型.
		- 实参变量名, 表面上得到复用.
		- 而本版本, 变量名不同

****************************************************/

//----------------------------------------------
//                自定义接口:
// 说明:
// 		- 用于扩展支持: os.Pipe() 和 io.Pipe()
//		- 俩包, 数据类型不同, 但行为一致.
//----------------------------------------------

type Reader interface {
	Read(data []byte) (n int, err error)
}

type Writer interface {
	Write(data []byte) (n int, err error)
}

//----------------------------------------------

func main() {
	// 基于文件io:
	reader, writer, _ := os.Pipe()
	do(reader, writer)
	fmt.Println("---------")

	// 基于内存io:
	r, w := io.Pipe() // 注意: 实参变量名, 与上不同
	do(r, w)

}

//----------------------------------------------
//
// 读写操作:
//		- 统一相同操作逻辑的代码段
// 形参类型:
// 		- 为自定义接口类型.
//		- 用于复用代码, 支持 os.Pipe() 和 io.Pipe()
//
//----------------------------------------------
func do(r Reader, w Writer) {
	//----------------------------------------------
	// 读取操作:
	//		- 新起一个go例程, 阻塞等.
	//		- 写入操作 sleep 多久, 则等多久.
	//----------------------------------------------

	go func() {
		output := make([]byte, 100) // 输出缓冲

		// 读取操作:
		//		- 从下面的 input读来, 写到 output中.
		n, err := r.Read(output)
		if err != nil {
			fmt.Printf("Error: Couldn't read data from pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s).\nresult: %s\n", n, output)
	}()

	//----------------------------------------------
	// 写入操作:
	//----------------------------------------------

	// write:
	input := make([]byte, 26)

	// 准备数据:
	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}

	time.Sleep(2 * time.Second) // 加入睡眠, 让 读取go例程, 等待数据.

	// 写入操作:
	n, err := w.Write(input) // 通过writer()写入数据.
	if err != nil {
		fmt.Printf("Error: Couldn't write data to pipe: %s\n", err)
	}

	fmt.Printf("Written %d byte(s).\n", n)
	time.Sleep(200 * time.Millisecond)

}
