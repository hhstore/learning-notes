package main

import (
	"fmt"
	"runtime"
)

/*
代码说明:
	1. init() 在main() 之前执行
	2. 若全局变量定义时未初始化, 可在 init()中初始化. 保证main()执行前.

*/

func init() {
	fmt.Printf("map: %v\n", m)

	// 初始化全局变量
	info = fmt.Sprintf("Os: [%s], Arch: [%s]", runtime.GOOS, runtime.GOARCH)
}

// 全局变量+初始化
var m = map[int]string{
	1: "a",
	2: "b",
	3: "c",
}

// 全局变量+未初始化, 在init()里作初始化操作
var info string

func main() {
	fmt.Println(info)

}
