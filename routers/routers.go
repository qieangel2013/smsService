package routers

import (
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"julive/middlewares"
	"smsService/controllers"
	_ "smsService/docs"
	"smsService/routers/v1"
	"strings"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {

	//初始化环境
	if strings.ToLower(viper.GetString("server.env")) == "prod" && strings.ToLower(viper.GetString("server.debug")) == "false" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	setUpConfig(router)
	setUpRouter(router)

	return router
}

// 初始化应用设置
func setUpConfig(router *gin.Engine) {

	// 使用swagger自动生成接口文档
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 使用Authorization中间件
	router.Use(controllers.Authorization)

	router.Use(middlewares.Logger)
	router.Use(middlewares.Recover)

	viewConfig := goview.DefaultConfig
	viewConfig.Root = "views"

	router.HTMLRender = ginview.New(viewConfig)

	// 未知路由处理
	router.NoRoute(controllers.Response404)
	// 未知调用方式
	router.NoMethod(controllers.Response404)

	router.Static("/images", "assets/images")

}

// 设置路由
func setUpRouter(router *gin.Engine) {
	api := router.Group("/api")
	{
		v1.RegisterRouter(api)
	}
}
