package api

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"julive/tools/coding"
	"smsService/models/system"
	lstructs "smsService/structs"
	"smsService/tools/helper"
	"strings"
)

// @Summary 发送短信
// @Schemes: http, https
// @Accept  json
// @Produce  json
// @Param  params body string  true "发短信参数" default({"tpl":"xKGgM","phone":"18614064093","params":{"code":6547,"minute":1}})
// @Success 200  "{"errCode":"0","msg":"","data":"","Action":"SendSms"}"
// @Failure 500
// @Router /api/v1/smsApi/SendSms [post]
func SendSms(c *gin.Context) {
	var result string
	var smsdataModel *lstructs.SmsData
	var data []interface{}
	var JsonParams lstructs.JsonParams

	// 解析JSON
	if c.BindJSON(&JsonParams) != nil {
		helper.Response(c, 200, gin.H{
			"errCode": 1,
			"msg":     "解析JSON参数失败",
			"data":    data,
		})
		return
	}

	tpls := &system.SmsTpl{}

	if JsonParams.TplId == "" {
		result = "模板id错误"
		helper.Response(c, 200, gin.H{
			"errCode": 1,
			"msg":     result,
			"data":    data,
		})
		return
	}
	if JsonParams.Phone == "" {
		result = "手机号参数错误"
		helper.Response(c, 200, gin.H{
			"errCode": 1,
			"msg":     result,
			"data":    data,
		})
		return
	}
	if JsonParams.Params == nil {
		result = "模板参数错误"
		helper.Response(c, 200, gin.H{
			"errCode": 1,
			"msg":     result,
			"data":    data,
		})
		return
	}

	if !c.GetBool(c.ClientIP()) {
		//验证Authorization
		auth := viper.GetString("server.authorization")
		authTmp, _ := coding.MD5(auth + ":" + JsonParams.Phone)
		if c.GetHeader("Authorization") != authTmp[8:24] {
			result = "Authorization认证失败"
			helper.Response(c, 403, gin.H{
				"errCode": 1,
				"msg":     result,
				"data":    data,
			})
			return
		}
	}

	tpl, err := tpls.GetCarrierByTpl(JsonParams.TplId)
	if err != nil {
		helper.Response(c, 500, gin.H{
			"errCode": 1,
			"msg":     err.Error(),
			"data":    data,
		})
	} else {
		//遍历运营商模板处理
		for _, v := range tpl {
			if strings.Contains(v.CarrierTplId, JsonParams.TplId) {
				//处理运营商发送短信
				smsdata := lstructs.SmsData{JsonParams.Phone, JsonParams.TplId, JsonParams.Params, *v, c.ClientIP(), c, data, result, 1, true}
				if c.GetBool(c.ClientIP()) {
					//异步模式
					smsdata.Sync = false
					lstructs.SmsDataChan <- smsdata
					helper.Response(c, 200, gin.H{
						"errCode": 1,
						"msg":     "异步成功!",
						"data":    data,
					})
				} else {
					//同步模式
					smsdataModel.SendSmsResponse(smsdata)
				}

			}
		}
	}
}
