package models

import (
// "julive/components/db"
// "julive/components/logger"
)

type BaseModel struct {
	ID       int64 `gorm:"column:id;primary_key" json:"id"`
	CreateAt int64 `gorm:"column:create_datetime" json:"create_at"`
	UpdateAt int64 `gorm:"column:update_datetime" json:"update_at"`
	CreateBy int64 `gorm:"column:creator" json:"create_by"`
	UpdateBy int64 `gorm:"column:updator" json:"update_by"`
}
