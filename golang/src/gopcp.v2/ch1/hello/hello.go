package main

import (
	"bufio"
	"os"
	"fmt"
)

func main() {
	// 带缓冲读取器 + 标准输入
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("please input your name:")

	// 换行符为分界符
	input, err := inputReader.ReadString('\n')

	if err != nil {
		fmt.Printf("Found an error: %s\n", err)

	} else {
		// 去掉换行符
		input = input[:len(input)-1]
		fmt.Printf("hello, %s!\n", input)
	}

}
