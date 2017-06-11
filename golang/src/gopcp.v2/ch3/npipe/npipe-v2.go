package main

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"time"
)

func main() {
	//var reader, writer reflect.Value
	//var err error

	// 基于文件io实现:
	osp := reflect.ValueOf(os.Pipe)
	pipeRet := osp.Call([]reflect.Value{})
	fmt.Println("\n====================\n[pipeResult]:", pipeRet)
	//fmt.Println(pipeRet)
	reader, writer, _ := pipeRet[0], pipeRet[1], pipeRet[2]
	do(reader, writer)

	// 基于内存io实现:
	iop := reflect.ValueOf(io.Pipe)
	pipeRet = iop.Call([]reflect.Value{})
	fmt.Println("\n====================\n[pipeResult]:", pipeRet)
	reader, writer = pipeRet[0], pipeRet[1]
	do(reader, writer)

}

// 利用反射reflect, 实现函数多态:
// 		- go原生不支持多态, 只能通过反射实现
//		- 比较繁琐, 不过能够完美实现需求
//		- 实现统一API: 当数据类型不同, 且行为接口一致时.
//
func do(reader reflect.Value, writer reflect.Value) {

	//----------------------------------------------
	// 读取操作:
	//		- 新起一个go例程, 阻塞等.
	//		- 写入操作 sleep 多久, 则等多久.
	//----------------------------------------------

	go func() {
		output := make([]byte, 100)

		read := reader.MethodByName("Read")
		result := read.Call([]reflect.Value{
			reflect.ValueOf(output),
		})
		_ = result
		fmt.Println("[output]:", string(output))

	}()

	//----------------------------------------------
	// 写入操作:
	//----------------------------------------------

	input := make([]byte, 26)

	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}
	fmt.Println("wait data...")
	time.Sleep(2 * time.Second) // 加入睡眠, 让 读取go例程, 等待数据.

	write := writer.MethodByName("Write")

	result := write.Call([]reflect.Value{
		reflect.ValueOf(input),
	})

	_ = result
	//fmt.Println(result)
	time.Sleep(200 * time.Millisecond)
}
