package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
说明:
	1. 交互式响应用户输入指令
	2. 根据用户输入指令, 反馈不同的操作.
	3. 涉及 for / switch 用法.

*/
func main() {
	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("please input your name:")

	// 读取字符串, 遇到换行符结束.
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Printf("an error ocuured: %s\n", err)
		os.Exit(1)
	} else {
		// 过滤掉换行符
		name := input[:len(input)-1]
		fmt.Printf("hello, %s! what can I do for you?\n", name)

	}

	// 死循环, 等待用户输入指令:
	for {
		// 读取字符串, 遇到换行符结束.
		input, err = inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("an error occurred: %s\n", err)
			continue
		}

		// 过滤换行符
		input = input[:len(input)-1]
		// 转换成全小写
		input = strings.ToLower(input)

		// 响应用户输入指令:
		switch input {
		case "":
			continue
		case "nothing", "bye":
			fmt.Println("bye!")
			os.Exit(0)
		default:
			fmt.Println("sorry, I did't catch you.")
		}
	}

}
