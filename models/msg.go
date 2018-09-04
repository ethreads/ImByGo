package models

import (
	"fmt"
	"gopkg.in/olahol/melody.v1"
	"encoding/json"
	"fpdxIm/config"
)

const (
	// 顺序不可改变
	ConnectType = iota 	// 连接消息
	DisconnectType		// 断开连接
	SystemType			// 系统消息
	MsgType				// 发送消息
	ReplyConn			// 连接回执
	ReplyMsg			// 消息回执
	ReplyError			// 异常回执
	ReplyWarn			// 警告回执
	ReplyInfo			// 提示回执
)

// 打包消息模型
type Msg struct {
	Ack string `json:"ask"`
	Type int `json:"type"`
	Msg *Message
	Session *melody.Session
}

// 客户端推送消息模型
type Message struct {
	Ack string `json:"ack"`
	FromUser string `json:"from_user"`
	ToUser string `json:"to_user"`
	MsgType int `json:"msg_type"`
	Data string `json:"data"`
}

// 发送队列任务模型
type SendMsg struct {
	Uuid string `json:"uuid"`
	Msg *Message
}

// 实例化一个回执消息
func (msg * Msg) NewReply(data string, Type int) (*Msg, error) {
	var response Msg
	var message Message
	message.Ack = msg.Ack
	message.FromUser = config.Config.Address
	message.ToUser = msg.Msg.FromUser
	message.Data = data
	response.Ack = msg.Ack
	response.Msg = &message
	switch Type {
	case ReplyConn:
		response.Type = ReplyConn
	case ReplyMsg:
		response.Type = ReplyMsg
	case ReplyError:
		response.Type = ReplyError
	case ReplyWarn:
		response.Type = ReplyWarn
	case ReplyInfo:
		response.Type = ReplyInfo
	default:
		return nil, fmt.Errorf("收到未知的回执类型:%d,消息Ack:%s", msg.Type, msg.Ack)
	}
	response.Session = msg.Session
	return &response, nil
}

// 实例化消息体
func NewMsg(data []byte, s *melody.Session, Type int) (*Msg, error) {
	var msg Msg
	var message Message
	switch Type {
	case ConnectType:
		msg.Type = ConnectType
		msg.Ack = "ack_connection"
		message.FromUser = config.Config.Address
		message.ToUser = s.Request.RemoteAddr
		message.Data = string(data)
	case DisconnectType:
		msg.Type = DisconnectType
		msg.Ack = "ack_disconnection"
		message.FromUser = config.Config.Address
		message.ToUser = s.Request.RemoteAddr
		message.Data = string(data)
	case SystemType:
		msg.Type = SystemType
		msg.Ack = "ack_system"
		message.FromUser = config.Config.Address
		message.ToUser = s.Request.RemoteAddr
		message.Data = string(data)
	case MsgType:
		if err := json.Unmarshal(data, &message); err != nil {
			return nil, fmt.Errorf("消息体反序列化:%s", err.Error())
		}
		if message.Ack == "" {
			return nil, fmt.Errorf("ack异常")
		}
		msg.Ack = message.Ack
		msg.Type = MsgType
	default:
		return nil, fmt.Errorf("未知的消息类型:%d", Type)
	}
	msg.Msg = &message
	msg.Session = s
	return &msg, nil
}