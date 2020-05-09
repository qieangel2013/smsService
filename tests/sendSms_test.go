package tests

import (
	"fmt"
	"julive/structs"
	"julive/tools/coding"
	"julive/tools/http"
	"testing"
)

func TestSendSms(t *testing.T) {
	authTmp, _ := coding.MD5("julive:18614064093")
	client := http.New()
	client.AddHeader("Authorization", authTmp[8:24])
	data, err := client.PostJson("http://192.168.234.131:9510/api/v1/smsApi/SendSms", structs.M{
		"tpl":   "xKGgM",
		"phone": "18614064093",
		"params": structs.M{
			"code":   6547,
			"minute": 1,
		},
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)
}
