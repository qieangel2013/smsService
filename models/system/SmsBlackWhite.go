package system

import (
	"github.com/spf13/viper"
	"julive/components/db"
	"julive/components/logger"
	"smsService/models"
	"strings"
)

type SmsBlackWhite struct {
	Phone    string `gorm:"column:phone" json:"phone"`
	BwType   int64  `gorm:"column:bw_type" json:"bw_type"`
	BwStatus int64  `gorm:"column:bw_status" json:"bw_status"`
	models.BaseModel
}

//获取表名
func (m *SmsBlackWhite) TableName() string {
	return "sms_black_white"
}

//获取黑白名单
func (m *SmsBlackWhite) GetBlackByPhone(phone string, isBlack bool) (SmsBlackWhite, error) {
	blackWhites := SmsBlackWhite{}
	if strings.ToLower(viper.GetString("server.debug")) == "true" {
		db.Get("system").LogMode(true)
	}
	var err error
	if isBlack {
		err = db.Get("system").Table(m.TableName()).Where("phone = ? AND bw_type = 2 AND bw_status = 1 ", phone).Limit(1).Find(&blackWhites).Error
	} else {
		err = db.Get("system").Table(m.TableName()).Where("phone = ? AND bw_type = 1 AND bw_status = 1 ", phone).Limit(1).Find(&blackWhites).Error
	}
	if err != nil {
		logger.Error("GetBlackByPhone执行查询异常：", err)
	}
	return blackWhites, nil
}

//添加运营商
func (m *SmsBlackWhite) Create() error {
	if strings.ToLower(viper.GetString("server.debug")) == "true" {
		db.Get("system").LogMode(true)
	}
	err := db.Get("system").Create(m).Error
	if err != nil {
		logger.Error("黑白名单添加异常：", err)
	}
	return nil
}

//更新运营商
func (m *SmsBlackWhite) Update() error {
	if strings.ToLower(viper.GetString("server.debug")) == "true" {
		db.Get("system").LogMode(true)
	}
	err := db.Get("system").Save(m).Error
	if err != nil {
		logger.Error("黑白名单更新异常：", err)
	}
	return nil
}
