package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"./chatbot" // 自定义包
)

/****************************************************

说明:
	1. 聊天机器人的双语言版本
	2. 注意 interface 接口实现的写法, 通常是定义指针类型
	3. flag包使用注意:
		- 命令行传参格式, 命名参数是key-value方式传参.
	4. 自定义包的导包路径:
		- 使用了相对路径, 避免路径不一致.
		- 正常发布项目, 应推荐使用绝对路径.

****************************************************/

//
// 全局变量:
// 	- init() 中给默认值
//  - 命令行解析, 传参修改
//
var chatbotName string

func init() {
	// 命令行参数, 给默认值:
	flag.StringVar(
		&chatbotName,
		"chatbot",
		"simple.en",
		"The chatbot's name for dialogue.",
	)

}

//
// 命令行传参获取方式:
// 2. 根据传参顺序位置获取
// 2. 根据参数名, key-value 方式.
//
func main() {
	// 通过命令行解析参数
	// 正确传参格式: [go run main.go -chatbot=simple.cn ]
	flag.Parse()
	//fmt.Println(flag.Args())

	// 注册2个子模块
	chatbot.Register(chatbot.NewSimpleEN("simple.en", nil))
	chatbot.Register(chatbot.NewSimpleCN("simple.cn", nil))

	myChatbot := chatbot.Get(chatbotName)
	if myChatbot == nil {
		err := fmt.Errorf("Fatal error: Unsupported chatbot named %s\n", chatbotName)
		checkError(nil, err, true)
	}

	// 创建读取器:
	inputReader := bufio.NewReader(os.Stdin)
	begin, err := myChatbot.Begin()
	checkError(myChatbot, err, true)
	fmt.Println(begin)

	// 提取输入内容, 过滤换行符
	input, err := inputReader.ReadString('\n')
	checkError(myChatbot, err, true)
	fmt.Println(myChatbot.Hello(input[:len(input)-1]))

	for {
		// 提取读取内容:
		input, err = inputReader.ReadString('\n')
		if checkError(myChatbot, err, false) {
			continue
		}

		output, end, err := myChatbot.Talk(input)
		if checkError(myChatbot, err, false) {
			continue
		}

		if output != "" {
			fmt.Println(output)
		}

		// 结束判断:
		if end {
			err = myChatbot.End()
			checkError(myChatbot, err, false)
			os.Exit(0)
		}

	}

}

// 校验:
func checkError(cb chatbot.Chatbot, err error, exit bool) bool {
	if err == nil {
		return false
	}

	if cb != nil {
		fmt.Println(cb.ReportError(err))
	} else {
		fmt.Println(err)
	}

	// 退出判断:
	if exit {
		debug.PrintStack()
		os.Exit(1)
	}

	return true
}
