package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"julive/components/logger"
	"julive/structs"
	"net/http"
	"smsService/models/system"
	"strconv"
	"strings"
)

func Interface2String(inter interface{}) string {

	switch inter.(type) {
	case string:
		return inter.(string)
	case int:
		return strconv.Itoa(inter.(int))
	case int64:
		return strconv.FormatInt(inter.(int64), 10)
	case float64:
		return strconv.FormatFloat(inter.(float64), 'f', -1, 64)
	}
	return ""
}

func FormateMsg(msg string, params interface{}) string {
	paramsArr := params.(map[string]interface{})
	for k, v := range paramsArr {
		if strings.Contains(msg, "@var("+k+")") {
			if v != "" {
				// fmt.Println(fmt.Sprintf("%T", v))
				msg = strings.Replace(msg, "@var("+k+")", Interface2String(v), -1)
			}
		}
	}
	return msg
}

func ParseHeader(header string) map[string]string {
	result := make(map[string]string)
	if strings.Contains(header, ",") {
		headerTmp := strings.Split(header, ",")
		for _, v := range headerTmp {
			if strings.Contains(v, ":") {
				headerMap := strings.Split(v, ":")
				result[headerMap[0]] = headerMap[1]
			}
		}
	} else {
		if strings.Contains(header, ":") {
			headerMapTmp := strings.Split(header, ":")
			result[headerMapTmp[0]] = headerMapTmp[1]
		}
	}
	return result
}

func ParseBody(data structs.M) map[string]string {
	result := make(map[string]string)
	for k, v := range data {
		result[k] = Interface2String(v)
	}
	return result
}

func Response(c *gin.Context, code int, obj gin.H) {
	c.JSON(code, obj)
}

func EncodeJson(value interface{}) string {
	carrier, err := json.Marshal(value)
	if err != nil {
		logger.Error("JSON编码异常：", err)
		return ""
	}
	return string(carrier)
}

func DecodeJson(str string) system.SmsCarrier {
	result := system.SmsCarrier{}
	if len(str) <= 0 {
		return result
	}
	d := json.NewDecoder(bytes.NewReader([]byte(str)))
	d.UseNumber()
	err := d.Decode(&result)
	if err != nil {
		logger.Error("JSON解码异常：", err)
	}
	return result
}

func DecodeBlackJson(str string) system.SmsBlackWhite {
	result := system.SmsBlackWhite{}
	if len(str) <= 0 {
		return result
	}
	d := json.NewDecoder(bytes.NewReader([]byte(str)))
	d.UseNumber()
	err := d.Decode(&result)
	if err != nil {
		logger.Error("JSON解码异常：", err)
	}
	return result
}

func ParseLimiter(data interface{}, id string) (int, string) {
	dataArr := data.(map[string]interface{})
	for k, v := range dataArr {
		vArr := v.(map[string]interface{})
		if _, ok := vArr["data"]; ok {
			vArrMap := vArr["data"].([]interface{})
			for _, vv := range vArrMap {
				if vv == id {
					i, err := strconv.Atoi(k)
					if err != nil {
						logger.Error("limiter解析异常：", err)
						return 0, "limiter解析异常"
					}
					return i, vArr["msg"].(string)
				}
			}
		}
	}
	return 0, "limiter解析异常"
}

func SendDingDing(phone string, tpl string, data string, carrier_name string, ip string) bool {
	formt := `### 基础短信服务服务监控报警 \n\n #### 发短信手机号:%s \n\n  #### 发短信运营商:%s \n\n  ### **错误信息**:<font color=#FF0000>%s</font> \n\n #### 短信模板:%s \n\n #### ip:%s \n\n `
	data = strings.Replace(data, `"`, `'`, -1)
	text := fmt.Sprintf(formt, phone, carrier_name, data, tpl, ip)
	content := `{"msgtype": "markdown",
					"markdown": {
            			"title":"基础短信服务监控报警",
            			"text": "` + text + `"
        			}
			}`
	req, err := http.NewRequest("POST", viper.GetString("server.dingding"), strings.NewReader(content))
	if err != nil {
		logger.Error("发送钉钉消息异常：", err)
		return false
	}
	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("发送钉钉消息异常：", err)
		return false
	}
	defer resp.Body.Close()
	return true
}

func AuthWhiteTable(ip string) bool {
	data := viper.Get("server.ipWhiteTable")
	for _, vv := range data.([]interface{}) {
		if strings.Contains(vv.(string), "/") {
			return ipMatch(ip, vv.(string))
		}
		if vv == ip {
			return true
		}
	}

	return false
}

func ip2binary(ip string) string {
	str := strings.Split(ip, ".")
	var ipstr string
	for _, s := range str {
		i, _ := strconv.ParseUint(s, 10, 8)
		ipstr = ipstr + fmt.Sprintf("%08b", i)
	}
	return ipstr
}

func ipMatch(ip string, iprange string) bool {
	ipb := ip2binary(ip)
	ipr := strings.Split(iprange, "/")
	masklen, err := strconv.ParseUint(ipr[1], 10, 32)
	if err != nil {
		logger.Error("ip地址段解析异常：", err)
		return false
	}
	iprb := ip2binary(ipr[0])
	return strings.EqualFold(ipb[0:masklen], iprb[0:masklen])
}
