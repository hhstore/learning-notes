package chatbot

import (
	"fmt"
	"strings"
)

// 聊天机器人-英文版
type simpleEN struct {
	name string
	talk Talk
}

// 构造函数:
func NewSimpleEN(name string, talk Talk) Chatbot {
	return &simpleEN{
		name: name,
		talk: talk,
	}
}

/****************************************************
              Chatbot 接口 - 英文版实现部分:

****************************************************/

func (robot *simpleEN) Name() string {
	return robot.name
}

func (robot *simpleEN) Begin() (string, error) {
	return "Please input your name:", nil
}

func (robot *simpleEN) Hello(userName string) string {
	userName = strings.TrimSpace(userName)
	if robot.talk != nil {
		return robot.talk.Hello(userName)
	}
	return fmt.Sprintf("Hello, %s! What can I do for you?", userName)

}

func (robot *simpleEN) Talk(heard string) (saying string, end bool, err error) {
	heard = strings.TrimSpace(heard)
	if robot.talk != nil {
		return robot.talk.Talk(heard)
	}

	switch heard {
	case "":
		return
	case "noting", "bye":
		saying = "bye!"
		end = true
		return
	default:
		saying = "Sorry, I didn't catch you."
		return
	}
}

func (robot *simpleEN) ReportError(err error) string {
	return fmt.Sprintf("An error occurred: %s\n", err)
}

func (robot *simpleEN) End() error {
	return nil
}
