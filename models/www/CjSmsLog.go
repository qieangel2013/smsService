package www

import (
	"github.com/spf13/viper"
	"julive/components/db"
	"julive/components/logger"
	"strconv"
	"strings"
	"time"
)

type CjSmsLog struct {
	ID             int64  `gorm:"column:id;primary_key" json:"id"`
	Phone          string `gorm:"column:phone" json:"phone"`
	Content        string `gorm:"column:content" json:"content"`
	ProvSmsid      string `gorm:"column:prov_smsid" json:"prov_smsid"`
	Status         int64  `gorm:"column:status" json:"status"`
	Provider       string `gorm:"column:provider" json:"provider"`
	CreateDatetime int64  `gorm:"column:create_datetime" json:"create_datetime"`
	ProvTime       string `gorm:"column:prov_time" json:"prov_time"`
	ReturnStr      string `gorm:"column:return_str" json:"return_str"`
	Tpl            string `gorm:"column:tpl" json:"tpl"`
	Ip             string `gorm:"column:ip" json:"ip"`
}

//获取表名
func (m *CjSmsLog) TableName() string {
	return "cj_sms_log_" + strconv.Itoa(time.Now().Year()) + time.Now().Format("01")
}

//获取当前手机号发送短信总数
func (m *CjSmsLog) GetCountByPhone(phone string, tpl string) int {
	var count int
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	timeNumber := t.Unix()
	start_time := timeNumber - 8*3600
	end_time := start_time + 24*3600
	if strings.ToLower(viper.GetString("server.debug")) == "true" {
		db.Get("www").LogMode(true)
	}
	err := db.Get("www").Table(m.TableName()).Where("phone = ? and tpl = ? and create_datetime between ? and ?", phone, tpl, start_time, end_time).Count(&count).Error
	if err != nil {
		logger.Error("GetCountByPhone执行查询总量异常：", err)
		return 0
	}
	return count
}

//添加短信日志
func (m *CjSmsLog) Create() error {
	if strings.ToLower(viper.GetString("server.debug")) == "true" {
		db.Get("www").LogMode(true)
	}
	tx := db.Get("www").Begin()
	m.CreateDatetime = time.Now().Unix()
	err := tx.Table(m.TableName()).Create(m).Error
	if err != nil {
		logger.Error("短信日志添加异常：", err)
		tx.Rollback()
	}
	tx.Commit()
	return nil
}

//更新短信日志
func (m *CjSmsLog) Update() error {
	if strings.ToLower(viper.GetString("server.debug")) == "true" {
		db.Get("www").LogMode(true)
	}
	err := db.Get("www").Save(m).Error
	if err != nil {
		logger.Error("短信日志更新异常：", err)
	}
	return nil
}
