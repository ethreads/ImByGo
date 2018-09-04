package services

import (
	"gopkg.in/olahol/melody.v1"
	"fpdxIm/models"
	"fpdxIm/config"
)

// 连接任务队列
var Connqueue = make(chan *models.Msg, config.Config.App["Connqueue"].(int))
// 断开连接任务队列
var DisConnqueue = make(chan *models.Msg, config.Config.App["DisConnqueue"].(int))
// 消息接收队列
var Msgqueue = make(chan *models.Msg, config.Config.App["Msgqueue"].(int))
// 消息发送队列
var Sendqueue = make(chan []byte, config.Config.App["Sendqueue"].(int))
// 连接池[uuid=>session]
var Conns = make(map[string]*melody.Session)
// 连接池的反向映射[session=>uuid]
var ConnsMap = make(map[*melody.Session]string)

func Init() {
	go connWorking()
	go sendWorking()
	go replyWorking()
}

// 移除slice中值为value的第一项
func remove(stack []string, value string) []string {
	var i int
	for k,v := range stack {
		i = k
		if v == value {
			break
		}
	}
	if i == len(stack) {
		return stack
	}
	copy(stack[i:], stack[i+1:])
	return stack[:len(stack)-1]
}