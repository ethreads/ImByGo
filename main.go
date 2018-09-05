package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"fpdxIm/services"
	"fpdxIm/models"
	"log"
	"github.com/DeanThompson/ginpprof"
)

func main() {
	r := gin.Default()
	ginpprof.Wrap(r)
	m := melody.New()
	models.Init()
	services.Init()

	// 升级http请求
	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	// 监听连接事件
	m.HandleConnect(func(s *melody.Session) {
		// 1. 实例化连接消息
		msg, err := models.NewMsg(nil, s, models.ConnectType)
		if err != nil {
			log.Printf("实例化连接消息:%s", err.Error())
		}
		// 2. 消息打包
		err = services.Package(msg)
		// 3. 递交连接任务
		if err != nil {
			log.Printf("连接消息打包:%s", err.Error())
			s.Close()
			return
		}
		services.Connqueue <- msg
	})

	// 监听连接断开事件
	m.HandleDisconnect(func(s *melody.Session) {
		// 1. 实例化断开连接消息
		msg, err := models.NewMsg(nil, s, models.DisconnectType)
		if err != nil {
			log.Printf("实例化断开连接消息:%s", err.Error())
			return
		}
		// 2. 消息打包
		err = services.Package(msg)
		if err != nil {
			log.Printf("断开连接消息打包:%s", err.Error())
			return
		}
		// 3. 递交断开连接任务
		services.DisConnqueue <- msg
	})

	// TODO 监听连接错误
	m.HandleError(func(s *melody.Session, e error) {
		log.Println("发生错误", e.Error())
	})

	// 监听接收事件
	m.HandleMessage(func(s *melody.Session, bytes []byte) {
		//1. 实例化消息
		msg, err := models.NewMsg(bytes, s, models.MsgType)
		if err != nil {
			log.Printf("实例化消息:%s", err.Error())
			return
		}
		//2. 消息打包
		err = services.Package(msg)
		if err != nil {
			log.Printf("消息打包:%s", err.Error())
			return
		}
		//3. 递交消息任务
		services.Msgqueue <- msg
	})

	r.Run(":9090")
}