package chatbot

import (
	"fmt"
	"strings"
)

// 聊天机器人-中文版
type simpleCN struct {
	name string
	talk Talk
}

// 构造函数:
func NewSimpleCN(name string, talk Talk) Chatbot {
	return &simpleCN{
		name: name,
		talk: talk,
	}
}

/****************************************************
              Chatbot 接口 - 中文版实现部分:

****************************************************/

func (robot *simpleCN) Name() string {
	return robot.name
}

func (robot *simpleCN) Begin() (string, error) {
	return "请输入你的名字:", nil
}

func (robot *simpleCN) Hello(userName string) string {
	userName = strings.TrimSpace(userName)
	if robot.talk != nil {
		return robot.talk.Hello(userName)
	}
	return fmt.Sprintf("你好, %s! 我可以为你做点什么?", userName)
}

func (robot *simpleCN) Talk(heard string) (saying string, end bool, err error) {
	heard = strings.TrimSpace(heard)
	if robot.talk != nil {
		return robot.talk.Talk(heard)
	}

	switch heard {
	case "":
		return
	case "没有", "再见":
		saying = "再见"
		end = true
		return
	default:
		saying = "对不起, 我没听懂你说什么."
		return

	}
}

func (robot *simpleCN) ReportError(err error) string {
	return fmt.Sprintf("发生了一个错误: %s\n", err)
}

func (robot *simpleCN) End() error {
	return nil

}
