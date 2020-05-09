package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"smsService/services/sms"
	"smsService/tools/helper"
)

var configPath string

func main() {
	defer func() {
		if err := recover(); err != nil {
			//发送钉钉报警
			helper.SendDingDing("", "", fmt.Sprintf("%s", err), "服务监控", "127.0.0.1")
			fmt.Println("panic", err)
			return
		}
	}()

	flag.StringVar(&configPath, "config", "", "配置文件路劲")
	flag.Parse()
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	if configPath == "" {
		viper.AddConfigPath("./conf")
	} else {
		viper.AddConfigPath(configPath[:len(configPath)-9])
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = sms.InitServer()
	if err != nil {
		panic(fmt.Errorf("server init error!: %s \n", err))
	}
	sms.Run()
}
