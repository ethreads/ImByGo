package services

import (
	"fpdxIm/models"
	"encoding/json"
	"log"
	"fmt"
	"strings"
)

// 消息打包
func Package(msg *models.Msg) error {
	return nil
}

// 接收消息处理
func message(msg *models.Msg)  {
	// 1. TODO 存档
	log.Printf("消息已存档\n")
	// 2. 回执
	replyMsg, err := msg.NewReply("发送成功", models.ReplyMsg)
	if err != nil {
		log.Printf("回执消息异常:%s", err.Error())
		return
	}
	Replyqueue <- replyMsg
	// 3. 获取to_user
	data, err  := models.Redis.HGet("ws:connpool", msg.Msg.ToUser).Result()
	if err != nil {
		log.Println("获取映射关系:", err.Error())
		return
	}
	// 4. 获取uuid列表
	var v []string
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		log.Println("反序列号映射关系:", err.Error())
		return
	}
	// 5. 推送至发送队列
	for _, uuid := range v {
		var sendMsg models.SendMsg
		sendMsg.Msg = msg.Msg
		sendMsg.Uuid = uuid
		strUuid := strings.Split(uuid, "@")
		strMsg, err := json.Marshal(sendMsg)
		if err != nil {
			log.Printf("发送消息序列化:%s", err.Error())
			continue
		}
		if err = models.Redis.Publish("ws:msgpool:" + strUuid[1], strMsg).Err(); err != nil {
			log.Printf("推送至UUid:%s发送队列:%s", uuid, err.Error())
			continue
		}
	}
}

// 消息分发中心
func connWorking() {
	for {
		select {
		case msg := <- Connqueue:
			connection(msg)
		case msg := <- DisConnqueue:
			disconnection(msg)
		case msg := <- Msgqueue:
			message(msg)
		}
	}
}
// 消息发送中心
func sendWorking()  {
	Sendqueue := models.Pubsub.Channel()
	for {
		select {
		case send := <- Sendqueue:
			var sendMsg models.SendMsg
			err := json.Unmarshal([]byte(send.Payload), &sendMsg)
			if err != nil {
				log.Printf("消息发送中心:反序列化发送任务包:%s", err.Error())
			}
			s, ok := Conns[sendMsg.Uuid]
			if !ok {
				// TODO 提交清道夫,清理redis残留的垃圾连接
				continue
			}
			data, err := json.Marshal(sendMsg.Msg)
			if err != nil {
				log.Printf("消息发送中心:序列化发送数据:%s", err.Error())
				continue
			}
			if s.IsClosed() {
				// TODO 提交清道夫,清理redis残留的垃圾连接
				fmt.Println("会话已关闭")
				continue
			}
			if err = s.Write(data); err != nil {
				log.Printf("消息发送中心:发送消息:%s", err.Error())
			}
		}
	}
}