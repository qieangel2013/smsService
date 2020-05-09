package v1

import (
	"github.com/gin-gonic/gin"
	"smsService/routers/v1/smsApi"
)

// RegisterRouter 注册路由
func RegisterRouter(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	{
		// smsApi路由
		smsApi.RegisterRouter(v1.Group("/smsApi"))
	}
}
