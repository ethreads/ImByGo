package services

import (
	"fpdxIm/models"
	"encoding/json"
	"log"
)

var Replyqueue = make(chan *models.Msg, 10000)

// 回执处理中心
func replyWorking()  {
	for {
		select {
		case reply := <- Replyqueue:
			data, err := json.Marshal(reply.Msg)
			if err != nil {
				log.Printf("回执处理中心:回执消息序列化:%s", err.Error())
			}
			if err := reply.Session.Write(data); err != nil {
				log.Printf("回执处理中心:发送回执消息:%s", err.Error())
			}
			if reply.Type == models.ReplyError {
				reply.Session.Close()
			}
		}
	}
}

