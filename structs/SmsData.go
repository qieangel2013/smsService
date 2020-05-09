package structs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"julive/components/logger"
	"julive/structs"
	"julive/tools/coding"
	"julive/tools/http"
	"smsService/caches/redis"
	"smsService/models/system"
	"smsService/models/www"
	"smsService/tools/helper"
	"strconv"
	"strings"
	"time"
)

type SmsData struct {
	Phone   string
	Tpl     string
	Params  interface{}
	SmsTpl  system.SmsTpl
	Ip      string
	Gct     *gin.Context
	Data    []interface{}
	Result  string
	IntFlag int
	Sync    bool
}

var SmsDataChan = make(chan SmsData, 1000)

func (s *SmsData) LoopInsertData(data chan SmsData) {
	phoneMap := make(map[string]int, 0)
	timerActive := time.NewTimer(60 * time.Second)
	for {
		select {
		case <-timerActive.C:
			//一分钟清空
			phoneMap = map[string]int{}
			timerActive.Reset(60 * time.Second)
		//SmsData
		case smsdata := <-data:
			phoneMap[smsdata.Phone]++
			if phoneMap[smsdata.Phone] > 3 {
				//发送钉钉报警
				helper.SendDingDing(smsdata.Phone, smsdata.Tpl, "短信每分钟发送超过3条", "服务监控", smsdata.Ip)
			} else {
				go s.SendSmsResponse(smsdata)
			}
		}
	}

}

func (s *SmsData) SendSmsResponse(data SmsData) {
	var status bool
	flag := viper.GetString("server.sms." + strconv.Itoa(data.IntFlag))
	carrier := s.getCarrierByName(flag)

	//验证白名单
	whiteState, whiteResp := s.checkWhite(data.Phone, data.Tpl, data.Params, data.SmsTpl.CarrierTplContent, carrier.CarrierName, data.Ip)
	if !whiteState {
		if data.Sync {
			helper.Response(data.Gct, 200, gin.H{
				"errCode": 1,
				"msg":     whiteResp,
				"data":    data.Data,
			})
		} else {
			logger.Error(data.Phone + whiteResp + ",模板id：" + data.Tpl + ",运营商:" + carrier.CarrierName + ",ip:" + data.Ip)
		}
		return
	}

	//验证黑名单
	state, resp := s.checkBlack(data.Phone, data.Tpl, data.Params, data.SmsTpl.CarrierTplContent, carrier.CarrierName, data.Ip)
	if !state {
		if data.Sync {
			helper.Response(data.Gct, 200, gin.H{
				"errCode": 1,
				"msg":     resp,
				"data":    data.Data,
			})
		} else {
			logger.Error(data.Phone + resp + ",模板id：" + data.Tpl + ",运营商:" + carrier.CarrierName + ",ip:" + data.Ip)
		}
		return
	}

	//验证发送限制
	LimitState, limitResp := s.checkLimiter(data.Phone, data.Tpl, carrier.CarrierName, data.Ip)
	if !LimitState {
		if data.Sync {
			helper.Response(data.Gct, 200, gin.H{
				"errCode": 1,
				"msg":     limitResp,
				"data":    data.Data,
			})
		} else {
			logger.Error(data.Phone + limitResp + ",模板id：" + data.Tpl + ",运营商:" + carrier.CarrierName + ",ip:" + data.Ip)
		}
		return
	}

	//优先发送flag为1的
	if carrier.CarrierName == "changzhuo" {
		status, data.Result = s.sendSmsMsg(carrier, data.Tpl, data.Phone, data.Params, &data.SmsTpl, data.Ip, true)
	} else {
		status, data.Result = s.sendSmsMsg(carrier, data.Tpl, data.Phone, data.Params, &data.SmsTpl, data.Ip, false)
	}
	if status {
		if data.Sync {
			helper.Response(data.Gct, 200, gin.H{
				"errCode": 1,
				"msg":     data.Result,
				"data":    data.Data,
			})
		} else {
			logger.Info(data.Phone + data.Result + ",模板id：" + data.Tpl + ",运营商:" + carrier.CarrierName + ",ip:" + data.Ip)
		}
	} else {
		data.IntFlag++
		//发送钉钉报警
		helper.SendDingDing(data.Phone, data.Tpl, data.Result, carrier.CarrierName, data.Ip)
		if viper.GetString("server.sms."+strconv.Itoa(data.IntFlag)) != "" {
			s.SendSmsResponse(data)
		} else {
			if data.Sync {
				helper.Response(data.Gct, 200, gin.H{
					"errCode": 1,
					"msg":     data.Result,
					"data":    data.Data,
				})
			} else {
				logger.Error(data.Phone + data.Result + ",模板id：" + data.Tpl + ",运营商:" + carrier.CarrierName + ",ip:" + data.Ip)
			}
		}
	}
}

func (s *SmsData) getCarrierByName(name string) system.SmsCarrier {
	//生成redis缓存
	cache := &redis.Cache{}
	value := cache.GetValue(name)
	carrier := system.SmsCarrier{}
	if value == "" {
		carriers := &system.SmsCarrier{}
		carrier, _ = carriers.GetCarriersByName(name)
		cache.SetValue(name, helper.EncodeJson(carrier))
	} else {
		carrier = helper.DecodeJson(value)
	}

	return carrier
}

func (s *SmsData) getBlackByPhone(phone string, isBlack bool) system.SmsBlackWhite {
	black := system.SmsBlackWhite{}
	cache := &redis.Cache{}
	var bkey string
	if isBlack {
		bkey = "black" + ":" + phone

	} else {
		bkey = "white" + ":" + phone
	}
	value := cache.GetValue(bkey)
	if value == "" {
		blacks := &system.SmsBlackWhite{}
		black, _ = blacks.GetBlackByPhone(phone, isBlack)
		if black.ID != 0 {
			cache.SetValue(bkey, helper.EncodeJson(black))
		} else {
			cache.SetValue(bkey, "no")
		}
	} else {
		if value == "no" {
			black = system.SmsBlackWhite{}
		} else {
			black = helper.DecodeBlackJson(value)
		}

	}

	return black
}

func (s *SmsData) checkWhite(phone string, tpl string, params interface{}, msg string, name string, ip string) (bool, string) {
	if strings.ToLower(viper.GetString("server.env")) == "test" {
		black := s.getBlackByPhone(phone, false)
		if black.ID > 0 {
		} else {
			//测试环境没有在白名单里写进数据库
			inserData := www.CjSmsLog{}
			inserData.Phone = phone
			inserData.Content = helper.FormateMsg(msg, params)
			inserData.Provider = name
			inserData.Tpl = tpl
			inserData.ReturnStr = "{}"
			inserData.Status = 0
			inserData.ProvSmsid = "test_send_id"
			inserData.Ip = ip
			dataTmp := InsertData{inserData, true}
			InsertDataChan <- dataTmp
			return false, "测试环境下phone为：" + phone + "没有在白名单里，不能发短信"
		}
	}
	return true, "ok"
}

func (s *SmsData) checkBlack(phone string, tpl string, params interface{}, msg string, name string, ip string) (bool, string) {
	black := s.getBlackByPhone(phone, true)
	if black.ID > 0 {
		//黑名单里写进数据库
		inserData := www.CjSmsLog{}
		inserData.Phone = phone
		inserData.Content = helper.FormateMsg(msg, params)
		inserData.Provider = name
		inserData.Tpl = tpl
		inserData.ReturnStr = "phone is in black list"
		inserData.Status = 0
		inserData.ProvSmsid = "black_send_id"
		inserData.Ip = ip
		dataTmp := InsertData{inserData, true}
		InsertDataChan <- dataTmp
		return false, "phone为：" + phone + "在黑名单里"
	}
	return true, "ok"
}

func (s *SmsData) checkLimiter(phone string, tpl string, name string, ip string) (bool, string) {
	limiter := viper.Get("server.limiter")
	cache := &redis.Cache{}
	limit, msg := helper.ParseLimiter(limiter, tpl)
	phoneLimit := &www.CjSmsLog{}
	lkey := "limit" + ":" + phone
	value := cache.GetValue(lkey)
	if value == "" {
		counter := phoneLimit.GetCountByPhone(phone, tpl)
		if counter > limit {
			helper.SendDingDing(phone, tpl, msg, name, ip)
			cache.SetValue(lkey, strconv.Itoa(counter))
			return false, msg
		}
	} else {
		return false, msg
	}

	return true, "ok"
}

func (s *SmsData) sendSmsMsg(carrier system.SmsCarrier, tpl string, phone string, params interface{}, tplData *system.SmsTpl, ip string, ctype bool) (bool, string) {
	if ctype {
		content := helper.FormateMsg(tplData.CarrierTplContent, params)
		if tpl == "XGeUK1" {
			content = content + "【居理旗舰店】"
		} else if tpl == "XGeUK1" {
			content = content + "【居理旗舰店】"
		} else {
			content = content + carrier.CarrierTplSuffix
		}
		sendData := structs.M{
			"account":  carrier.CarrierAppid,
			"password": carrier.CarrierToken,
			"mobile":   phone,
			"content":  content,
		}
		return s.sendSubSmsExt(carrier.CarrierGateway, carrier.CarrierTplHeader, carrier.CarrierTplSuffix, sendData, tplData, carrier, ip)
	} else {
		sendData := structs.M{
			"appid":     carrier.CarrierAppid,
			"signature": carrier.CarrierToken,
			"project":   tpl,
			"to":        phone,
			"vars":      params,
		}
		return s.sendSubSms(carrier.CarrierGateway, carrier.CarrierTplHeader, carrier.CarrierTplSuffix, sendData, tplData, carrier, ip)
	}
}

func (s *SmsData) sendSubSmsExt(url string, header string, suffix string, data structs.M, tplData *system.SmsTpl, carrier system.SmsCarrier, ip string) (bool, string) {
	httpClient := http.New()
	if header != "" {
		for k, v := range helper.ParseHeader(header) {
			httpClient.AddHeader(k, v)
		}
	}
	inserData := www.CjSmsLog{}
	for k, v := range helper.ParseBody(data) {
		httpClient.AddQuery(k, v)
	}
	result, err := httpClient.PostRaw(url, "")
	if err != nil {
		logger.Error("发送的url为："+url+"，post请求发送短信异常：", err)
		return false, fmt.Sprintf("%s", err)
	}

	if coding.JSONGetInt(result, "status") >= 0 {
		inserData.Phone = data["mobile"].(string)
		inserData.Content = data["content"].(string)
		inserData.Provider = carrier.CarrierName
		inserData.Tpl = tplData.CarrierTplId
		inserData.ReturnStr = result
		inserData.Status = 1
		inserData.ProvSmsid = coding.JSONGet(result, "taskId").(string)
		inserData.Ip = ip
		dataTmp := InsertData{inserData, true}
		InsertDataChan <- dataTmp
		return true, "ok"
	} else {
		inserData.Phone = data["mobile"].(string)
		inserData.Content = data["content"].(string)
		inserData.Provider = carrier.CarrierName
		inserData.Tpl = tplData.CarrierTplId
		inserData.ReturnStr = result
		inserData.Status = 0
		inserData.ProvSmsid = ""
		inserData.Ip = ip
	}
	dataTmp := InsertData{inserData, false}
	InsertDataChan <- dataTmp
	return false, result
}

func (s *SmsData) sendSubSms(url string, header string, suffix string, data structs.M, tplData *system.SmsTpl, carrier system.SmsCarrier, ip string) (bool, string) {
	httpClient := http.New()
	if header != "" {
		for k, v := range helper.ParseHeader(header) {
			httpClient.AddHeader(k, v)
		}
	}
	inserData := www.CjSmsLog{}
	//对于https设置证书忽略
	httpClient.SetSkipVerify()
	result, err := httpClient.PostJson(url, data)
	if err != nil {
		logger.Error("发送的url为："+url+"，post请求发送短信异常：", err)
		return false, fmt.Sprintf("%s", err)
	}

	if coding.JSONGet(result, "status") == "success" {
		inserData.Phone = data["to"].(string)
		inserData.Content = helper.FormateMsg(tplData.CarrierTplContent, data["vars"])
		inserData.Provider = carrier.CarrierName
		inserData.Tpl = tplData.CarrierTplId
		inserData.ReturnStr = result
		inserData.Status = 1
		inserData.ProvSmsid = coding.JSONGet(result, "send_id").(string)
		inserData.Ip = ip
		dataTmp := InsertData{inserData, true}
		InsertDataChan <- dataTmp
		return true, "ok"
	} else {
		inserData.Phone = data["to"].(string)
		inserData.Content = helper.FormateMsg(tplData.CarrierTplContent, data["vars"])
		inserData.Provider = carrier.CarrierName
		inserData.Tpl = tplData.CarrierTplId
		inserData.ReturnStr = result
		inserData.Status = 0
		inserData.ProvSmsid = ""
		inserData.Ip = ip
	}
	dataTmp := InsertData{inserData, false}
	InsertDataChan <- dataTmp
	return false, result
}
