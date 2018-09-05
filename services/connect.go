package services

import (
	"fpdxIm/models"
	"log"
	"strconv"
	"fpdxIm/config"
	"time"
	"encoding/json"
	"strings"
)

// 建立连接处理
func connection(msg *models.Msg)  {
	// 1. 认证
	uid, err := Auth(msg.Session.Request)
	if err != nil {
		replyMsg, err := msg.NewReply("建立连接处理:认证失败" + err.Error(), models.ReplyError)
		if err != nil {
			log.Printf("建立连接处理:回执消息异常:%s", err.Error())
		}
		Replyqueue <- replyMsg
		return
	}
	// 2. 生成uuid
	uuid := uid + "@" + config.Config.Address + "@" + strconv.FormatInt(time.Now().UnixNano(), 10)
	// 3. 建立uid-uuid映射
	data, _  := models.Redis.HGet("ws:connpool", uid).Result()
	var v []string
	if data != "" {
		if err := json.Unmarshal([]byte(data), &v); err != nil {
			log.Println("建立连接处理:反序列号映射关系:", err.Error())
			replyMsg, err := msg.NewReply("连接失败:服务器异常", models.ReplyError)
			if err != nil {
				log.Printf("建立连接处理:回执消息异常:%s", err.Error())
			}
			Replyqueue <- replyMsg
			return
		}
	}
	v = append(v, uuid)
	str, err := json.Marshal(v)
	if err != nil {
		log.Println("建立连接处理:序列号映射关系:", err.Error())
		replyMsg, err := msg.NewReply("连接失败:服务器异常", models.ReplyError)
		if err != nil {
			log.Printf("建立连接处理:回执消息异常:%s", err.Error())
		}
		Replyqueue <- replyMsg
		return
	}
	if _, err = models.Redis.HSet("ws:connpool", uid, str).Result(); err != nil {
		log.Println("建立连接处理:Redis加入映射关系:", err.Error())
		replyMsg, err := msg.NewReply("连接失败:服务器异常", models.ReplyError)
		if err != nil {
			log.Printf("建立连接处理:回执消息异常:%s", err.Error())
		}
		Replyqueue <- replyMsg
		return
	}
	// 4. 注册map
	Conns[uuid] = msg.Session
	ConnsMap[msg.Session] = uuid
	// 5. 回执消息
	replyMsg, err := msg.NewReply("连接成功", models.ReplyConn)
	if err != nil {
		log.Printf("建立连接处理:回执消息异常:%s", err.Error())
		return
	}
	Replyqueue <- replyMsg
	return
}

// 断开连接处理
func disconnection(msg *models.Msg)  {
	// 1. 清理redis[uid<->uuid]的映射
	uuid, ok := ConnsMap[msg.Session]
	if !ok {
		log.Printf("断开连接处理:ack:%v,未能找到连接映射", msg.Ack)
		return
	}
	strUuid := strings.Split(uuid, "@")
	uid := strUuid[0]
	data, err  := models.Redis.HGet("ws:connpool", uid).Result()
	if err != nil {
		log.Printf("断开连接处理:ack:%v,获取映射关系:%s", msg.Ack, err.Error())
		return
	}
	var v []string
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		log.Println("断开连接处理:反序列号映射关系:", err.Error())
		return
	}
	v = remove(v, uuid)
	if len(v) > 0 {
		str, err := json.Marshal(v)
		if err != nil {
			log.Println("断开连接处理:序列号映射关系:", err.Error())
			return
		}
		if _, err = models.Redis.HSet("ws:connpool", uid, str).Result(); err != nil {
			log.Println("断开连接处理:Redis加入映射关系:", err.Error())
		}
	} else {
		if err = models.Redis.HDel("ws:connpool", uid).Err(); err != nil {
			log.Println("断开连接处理:Redis映射关系置空:", err.Error())
		}
	}
	// 2. 清理注册表[uuid<->session]的map
	delete(Conns, uuid)
	delete(ConnsMap, msg.Session)
	return
}
