package smsApi

import (
	"github.com/gin-gonic/gin"
	"smsService/controllers"
	"smsService/controllers/api"
)

// RegisterRouter 注册路由
func RegisterRouter(r *gin.RouterGroup) {
	r.POST("/SendSms", api.SendSms)
	r.GET("/SendSms", controllers.Response403)
}
