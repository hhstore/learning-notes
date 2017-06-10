package chatbot

import "errors"

//
// 接口: 聊天功能
//
type Talk interface {
	Hello(userName string) string
	Talk(heard string) (saying string, end bool, err error)
}

//
// 接口: 聊天机器人
//	- 复用 Talk 接口
//
type Chatbot interface {
	Talk

	Name() string
	Begin() (string, error)
	ReportError(err error) string
	End() error
}

// 出错提示信息:
var (
	ErrInvalidChatbotName = errors.New("Invalid chatbot name")
	ErrInvalidChatbot     = errors.New("Invalid chatbot")
	ErrExistingChatbot    = errors.New("Existing chatbot")
)

// 全局字典: 存储机器人
var chatbotMap = map[string]Chatbot{}

//
// 注册:
//
func Register(chatbot Chatbot) error {
	if chatbot == nil {
		return ErrInvalidChatbot
	}

	name := chatbot.Name()
	if name == "" {
		return ErrInvalidChatbotName
	}
	if _, ok := chatbotMap[name]; ok {
		return ErrExistingChatbot
	}

	// 成功注册
	chatbotMap[name] = chatbot
	return nil
}

//
// 获取指定名称的机器人
//	- 依赖全局字典: chatbotMap
//
func Get(name string) Chatbot {
	return chatbotMap[name]
}
