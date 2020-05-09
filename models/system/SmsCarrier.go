package system

import (
	"github.com/spf13/viper"
	"julive/components/db"
	"julive/components/logger"
	"smsService/models"
	"strings"
)

type SmsCarrier struct {
	CarrierName      string `gorm:"column:carrier_name" json:"carrier_name"`
	CarrierTplHeader string `gorm:"column:carrier_tpl_header" json:"carrier_tpl_header"`
	CarrierTplSuffix string `gorm:"column:carrier_tpl_suffix" json:"carrier_tpl_suffix"`
	CarrierTplParams string `gorm:"column:carrier_tpl_params" json:"carrier_tpl_params"`
	CarrierGateway   string `gorm:"column:carrier_gateway" json:"carrier_gateway"`
	CarrierAppid     string `gorm:"column:carrier_appid" json:"carrier_appid"`
	CarrierSid       string `gorm:"column:carrier_sid" json:"carrier_sid"`
	CarrierToken     string `gorm:"column:carrier_token" json:"carrier_token"`
	CarrierStatus    int64  `gorm:"column:carrier_status" json:"carrier_status"`
	models.BaseModel
}

//获取表名
func (m *SmsCarrier) TableName() string {
	return "sms_carrier"
}

//获取运营商
func (m *SmsCarrier) GetCarriers(id int64) (SmsCarrier, error) {
	carriers := SmsCarrier{}
	if strings.ToLower(viper.GetString("server.debug")) == "true" {
		db.Get("system").LogMode(true)
	}
	err := db.Get("system").Table(m.TableName()).Where("id = ? AND carrier_status = 1", id).Limit(1).Find(&carriers).Error
	if err != nil {
		logger.Error("GetCarriers执行查询异常：", err)
	}
	return carriers, nil
}

//获取运营商
func (m *SmsCarrier) GetCarriersByName(name string) (SmsCarrier, error) {
	carriers := SmsCarrier{}
	if strings.ToLower(viper.GetString("server.debug")) == "true" {
		db.Get("system").LogMode(true)
	}
	err := db.Get("system").Table(m.TableName()).Where("carrier_name = ? AND carrier_status = 1", name).Limit(1).Find(&carriers).Error
	if err != nil {
		logger.Error("GetCarriersByName执行查询异常：", err)
	}
	return carriers, nil
}

//添加运营商
func (m *SmsCarrier) Create() error {
	if strings.ToLower(viper.GetString("server.debug")) == "true" {
		db.Get("system").LogMode(true)
	}
	err := db.Get("system").Create(m).Error
	if err != nil {
		logger.Error("运营商添加异常：", err)
	}
	return nil
}

//更新运营商
func (m *SmsCarrier) Update() error {
	if strings.ToLower(viper.GetString("server.debug")) == "true" {
		db.Get("system").LogMode(true)
	}
	err := db.Get("system").Save(m).Error
	if err != nil {
		logger.Error("运营商更新异常：", err)
	}
	return nil
}
