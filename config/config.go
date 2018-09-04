package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"net"
	"log"
)

type config struct {
	DB map[string]interface{} `yaml:db`
	Redis map[string]interface{} `yaml:redis`
	App map[string]interface{} `yaml:app`
	Address string
}

var Config = config{
	DB: make(map[string]interface{}),
	Redis: make(map[string]interface{}),
	App: make(map[string]interface{}),
	Address: "127.0.0.1",
}

func parseYAML() {
	data, err := ioutil.ReadFile("./conf.yaml")
	if err != nil {
		log.Fatal("Conf file ReadFile:", err.Error())
	}
	if err := yaml.Unmarshal(data, &Config); err != nil {
		log.Fatal("Conf file ReadFile:", err.Error())
	}
}

// 获取服务器IP
func GetLocalIp() string {
	addrSlice, err := net.InterfaceAddrs()
	if  err != nil {
		log.Println("Get local IP addr failed!")
		return "localhost"
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}

func init()  {
	parseYAML()
	Config.Address = GetLocalIp()
}