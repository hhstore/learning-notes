package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"time"
)

/****************************************************

        多进程编程 - IPC通信方式 - 管道

- 基于管道方式:
	- 基于 命名管道 和 内存管道

说明:
	1. 基于反射reflect, 实现函数的多态(行为相同,参数类型不同)
	2. 基于reflect各种接口方法, 获取参数, 执行函数调用
	3. 注意: 通过反射调用函数, 返回结果类型与常规函数调用不同.
	4. 理解反射的重要作用, 实现对 鸭子类型 变量的支持.
	5. 运行时, 获取变量类型.
	6. 当然 API接口的参数类型, 都要定义成 reflect.Value 类型
	7. reflect包重要接口:
		- reflect.ValueOf
		- v.Call(args)
		- v.MethodByName().Call()

****************************************************/

func main() {
	//fmt.Println(reflect.TypeOf(os.Pipe))
	//fmt.Println(reflect.TypeOf(io.Pipe))

	//fmt.Println(reflect.ValueOf(os.Pipe))
	//fmt.Println(reflect.ValueOf(io.Pipe))

	//v := reflect.ValueOf(os.Pipe)
	//args := []reflect.Value{}
	//ret := v.Call(args)
	//fmt.Println("[kind]", v.Kind(), "\t[Type] ", v.Type())
	//fmt.Println("[result]", ret)

	// 基于文件io实现:
	osPipe := reflect.ValueOf(os.Pipe)
	doPipe(osPipe)

	// 基于内存io实现:
	ioPipe := reflect.ValueOf(io.Pipe)
	doPipe(ioPipe)

}

func doPipe(pipeFunc reflect.Value) {
	var reader, writer reflect.Value

	pipeArgs := []reflect.Value{}      // 空参数, 因为 Pipe()参数为空
	pipeRet := pipeFunc.Call(pipeArgs) // 执行 Pipe() 方法

	fmt.Println("\n====================\n[pipeResult]:", pipeRet)
	if len(pipeRet) == 3 { // 基于os.Pipe(), 命名管道
		reader, writer, _ = pipeRet[0], pipeRet[1], pipeRet[2]
	} else { // 基于io.Pipe(), 内存管道
		reader, writer = pipeRet[0], pipeRet[1]
	}

	go func() {
		output := make([]byte, 100)

		read := reader.MethodByName("Read")

		read.Call([]reflect.Value{
			reflect.ValueOf(output),
		})

		fmt.Println("[output]:", string(output))

	}()

	//----------------------------------------------
	// 写入操作:
	//----------------------------------------------

	input := make([]byte, 26)

	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}

	fmt.Println("\t wait data...")
	time.Sleep(2 * time.Second) // 加入睡眠, 让 读取go例程, 等待数据.

	write := writer.MethodByName("Write")

	result := write.Call([]reflect.Value{
		reflect.ValueOf(input),
	})

	_ = result

	//fmt.Println(result)

	time.Sleep(200 * time.Millisecond)

}
