package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

/****************************************************

说明:
	1. 区别 os.pipe() 和 io.pipe()
		- 2个lib的对外操作接口一致, 但数据类型不同.
		- 理解 reader 和 writer 操作行为
		- 读取操作 和 写入操作
			- 读取操作, 新启一个go例程, 阻塞等数据
			- 为何是阻塞的? 加入 sleep()作了验证.

	2. 理解静态类型语言和动态类型语言, 本质差异:
		- 动态类型:
			- 如Python, 这里2个操作函数, 完全可以合并成一个.
			- 对于2个lib中 reader 和 writer 参数类型不同.
			- Python是可以等到运行时, 再确定参数类型.
			- 本质上是变量名的复用, 2个不同类型的变量, 用同一个参数名而已.
		- 静态类型:
			- go这里,就不可以合并2个子函数
			- 因为 reader 和 writer 的数据类型不同
			- go在编译时, 就要确定参数类型, 传入的参数类型要匹配, 不匹配, 就要写成2个.
			- go并能定义一个 temp, 把 A 和 B 不同类型, 强制转换赋值给 temp.
		- 数据类型不同 vs 行为一致:
			- python 是真鸭子类型, 可以合并API操作逻辑, 动态类型的特点.
			- go 无法在数据类型不同的前提下, 合并API操作逻辑.
			- go 并不是真的鸭子类型. (虽然面向接口编程, 很像)

	3. 理解面向接口编程:
		- 拆开 数据 和 行为
		- 本质上是面向行为编程,
		- 实现函数多态:
			- 表面上, 统一行为操作.
			- 形参设计: 利用 interface {} 扩展
			- 行为优先, 若数据类型不同, 扩展新的 interface{} 用以支持 "一致的行为"
			- 数据 为 行为 服务.
			- 若行为统一, 根源数据类型不同, 则相似逻辑代码不可直接合并.
			- 解决办法:
				- 扩展根源数据类型, 即: 统一 A, B 类型的 interface{}源
				- 这样 A, B 的 数据源(interface) 和 操作行为, 都得到统一.
				- 实现函数多态.
				- - 理解go作为静态类型语言本身的限制.
	4. 函数多态设计:
		- 基于 interface{} 实现
		- 基于反射 reflect 实现

****************************************************/

func main() {
	fileBasePipe() // 基于文件io
	fmt.Println("---------")
	inMemorySyncPipe() // 基于内存io
}

func fileBasePipe() {
	reader, writer, err := os.Pipe() // 区别: os包
	doByOSPipe(reader, writer, err)
}

func inMemorySyncPipe() {
	reader, writer := io.Pipe() // 区别: io包
	doByIOPipe(reader, writer)
}

func doByOSPipe(reader *os.File, writer *os.File, err error) {
	if err != nil {
		fmt.Printf("Error: Coundn't create pipe: %s\n", err)
	}

	//----------------------------------------------
	// 读取操作:
	//		- 新起一个go例程, 阻塞等.
	//		- 写入操作 sleep 多久, 则等多久.
	//----------------------------------------------

	go func() {
		output := make([]byte, 100) // 输出缓冲

		// 读取操作:
		//		- 从下面的 input读来, 写到 output中.
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Couldn't read data from pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s).\nresult: %s\n", n, output)
	}()

	//----------------------------------------------
	// 写入操作:
	//----------------------------------------------

	input := make([]byte, 26) // 写入缓冲

	// 准备数据:
	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}

	time.Sleep(2 * time.Second) // 加入睡眠, 让 读取go例程, 等待数据.

	// 写入操作:
	n, err := writer.Write(input) // 通过writer()写入数据.
	if err != nil {
		fmt.Printf("Error: Couldn't write data to pipe: %s\n", err)
	}

	fmt.Printf("Written %d byte(s).\n", n)
	time.Sleep(200 * time.Millisecond)

}

func doByIOPipe(reader *io.PipeReader, writer *io.PipeWriter) {
	//----------------------------------------------
	// 读取操作:
	//		- 新起一个go例程, 阻塞等.
	//		- 写入操作 sleep 多久, 则等多久.
	//----------------------------------------------

	go func() {
		output := make([]byte, 100)
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Couldn't read data from pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s).\n", n)
	}()

	//----------------------------------------------
	// 写入操作:
	//----------------------------------------------

	input := make([]byte, 26)

	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}

	n, err := writer.Write(input)
	if err != nil {
		fmt.Printf("Error: Couldn't write data to pipe: %s\n", err)
	}

	fmt.Printf("Written %d byte(s).\n", n)
	time.Sleep(200 * time.Millisecond)
}
