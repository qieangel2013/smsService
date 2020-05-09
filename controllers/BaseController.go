package controllers

import (
	"github.com/gin-gonic/gin"
	"smsService/tools/helper"
)

// @Summary 无路由或者无方法返回404
// @Produce  html
// @Success 200
// @Failure 500
// @Router /404 [get]
func Response404(c *gin.Context) {
	//异常模板
	c.HTML(404, "layouts/404.html", gin.H{
		"title": "404-对不起！您访问的页面不存在-居理新房",
	})
}

// @Summary 访问403
// @Produce  html
// @Success 200
// @Failure 500
// @Router /403 [get]
func Response403(c *gin.Context) {
	//异常模板
	c.HTML(403, "layouts/403.html", gin.H{
		"title": "403-网页访问错误-居理新房",
	})
}

//Authorization认证
func Authorization(c *gin.Context) {

	if c.GetHeader("Authorization") != "" {
	} else if helper.AuthWhiteTable(c.ClientIP()) {
		c.Set(c.ClientIP(), true)
	} else {
		c.HTML(403, "layouts/403.html", gin.H{
			"title": "403-网页访问错误-居理新房",
		})
		c.Abort()
		return
	}
	c.Next()
}
