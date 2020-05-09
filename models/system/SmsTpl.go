package system

import (
	"github.com/spf13/viper"
	"julive/components/db"
	"julive/components/logger"
	"smsService/models"
	"strings"
)

type SmsTpl struct {
	CarrierId         string `gorm:"column:carrier_id" json:"carrier_id"`
	CarrierTplId      string `gorm:"column:carrier_tpl_id" json:"carrier_tpl_id"`
	CarrierTplContent string `gorm:"column:carrier_tpl_content" json:"carrier_tpl_content"`
	CarrierTplType    int64  `gorm:"column:carrier_tpl_type" json:"carrier_tpl_type"`
	CarrierTplHeader  string `gorm:"column:carrier_tpl_header" json:"carrier_tpl_header"`
	CarrierTplSuffix  string `gorm:"column:carrier_tpl_suffix" json:"carrier_tpl_suffix"`
	CarrierTpl        string `gorm:"column:carrier_tpl" json:"carrier_tpl"`
	TplStatus         int64  `gorm:"column:tpl_status" json:"tpl_status"`
	models.BaseModel
}

//获取表名
func (m *SmsTpl) TableName() string {
	return "sms_tpl"
}

//获取运营商模板
func (m *SmsTpl) GetCarrierByTpl(tpl string) ([]*SmsTpl, error) {
	tpls := []*SmsTpl{}
	if strings.ToLower(viper.GetString("server.debug")) == "true" {
		db.Get("system").LogMode(true)
	}
	err := db.Get("system").Table(m.TableName()).Where("carrier_tpl_id = ? AND tpl_status = 1", tpl).Limit(10).Find(&tpls).Error
	if err != nil {
		logger.Error("GetCarrierByTpl执行查询异常：", err)
		return nil, err
	}
	return tpls, nil
}

//添加运营商模板
func (m *SmsTpl) Create() error {
	if strings.ToLower(viper.GetString("server.debug")) == "true" {
		db.Get("system").LogMode(true)
	}
	err := db.Get("system").Create(m).Error
	if err != nil {
		logger.Error("运营商添加异常：", err)
	}
	return nil
}

//更新运营商模板
func (m *SmsTpl) Update() error {
	if strings.ToLower(viper.GetString("server.debug")) == "true" {
		db.Get("system").LogMode(true)
	}
	err := db.Get("system").Save(m).Error
	if err != nil {
		logger.Error("运营商更新异常：", err)
	}
	return nil
}
