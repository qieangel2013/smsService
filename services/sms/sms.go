package sms

import (
	"fmt"
	"github.com/spf13/viper"
	"julive/components/cache"
	"julive/components/db"
	"julive/components/logger"
	"smsService/routers"
	"smsService/structs"
	"smsService/tools/helper"
)

var registers = [...]func() error{
	db.Register,
	logger.Register,
	cache.Register,
}
var insertData *structs.InsertData
var smsData *structs.SmsData

func InitServer() (err error) {
	//注册组件
	for _, register := range registers {
		err = register()
		if err != nil {
			return
		}
	}
	//初始化插入数据管道
	go insertData.LoopInsertData(structs.InsertDataChan)
	//初始化接受短信的管道
	go smsData.LoopInsertData(structs.SmsDataChan)

	return nil
}
func Run() {
	addr := fmt.Sprintf("%s:%s", viper.GetString("server.ListenHost"), viper.GetString("server.ListenPort"))
	r := routers.InitRouter()
	err := r.Run(addr)
	if err != nil {
		logger.Error("短信基础服务启动失败!")
		//发送钉钉报警
		helper.SendDingDing("", "", "短信基础服务启动失败!", "服务监控", "127.0.0.1")
		panic("短信基础服务启动失败!")
	}
}
