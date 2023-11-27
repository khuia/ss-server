package main

import (
	"encoding/json"

	"log"
	"os"
)

type Config struct {
	Type string

	LocalAddr string

	RemoteAddr string

	Socks5Addr string

	Key string
}

func getConfig(configPath string, config *Config) {

	// 读取配置文件
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// 解析配置文件

	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatal(err)
	}

}
