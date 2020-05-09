package dingding

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

/**
 * 发送钉钉的底层方法
 * 钉钉文档:https://ding-doc.dingtalk.com/doc#/serverapi2/ye8tup/7ae8ebf3
 */
func SendDingDing(webhook string, content string) (code int, err error) {
	//创建一个请求
	req, err := http.NewRequest("POST", webhook, strings.NewReader(content))
	if err != nil {
		return 1, err
	}

	client := &http.Client{}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		return 1, err
	}
	//最后关闭请求
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 1, err
	}
	//json str 转map
	var ding map[string]interface{}
	unerr := json.Unmarshal(body, &ding)
	if unerr != nil {
		return 1, unerr
	}
	// interface类型转float64 和 interface转string
	return int(ding["errcode"].(float64)), fmt.Errorf(ding["errmsg"].(string))
}

/**
 * 发送信息
 */
func SendMessage(webhook string, msg string) (code int, err error) {

	var content = make(map[string]interface{})
	content["msgtype"] = "text"
	content["text"] = map[string]string{"content": msg}
	json, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	return SendDingDing(webhook, string(json))
}

/**
 * 发送信息并at部分人
 */
func SendMessageAtMobiles(webhook string, msg string, atMobiles []string) (code int, err error) {
	var content = make(map[string]interface{})
	content["msgtype"] = "text"
	content["text"] = map[string]string{"content": msg}
	if len(atMobiles) > 0 {
		var at = make(map[string]interface{})
		at["atMobiles"] = atMobiles
		content["at"] = at
	}
	json, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	return SendDingDing(webhook, string(json))
}

/**
 * 发送信息并at所有人
 */
func SendMessageAtAll(webhook string, msg string, isAtAll bool) (code int, err error) {
	/*
		content := `{
			"msgtype": "text",
			"text": {
				"content": "` + msg + `"
			},
			"at" : {
					"atMobiles": "",
					"isAtAll": isAtAll
				}
		}`
	*/
	var content = make(map[string]interface{})
	content["msgtype"] = "text"
	content["text"] = map[string]string{"content": msg}
	if isAtAll {
		var at = make(map[string]interface{})
		at["isAtAll"] = isAtAll
		content["at"] = at
	}

	json, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	return SendDingDing(webhook, string(json))
}

/**
 * 发送信息
 */
func SendMarkdown(webhook string, title string, msg string) (code int, err error) {

	var content = make(map[string]interface{})
	content["msgtype"] = "markdown"
	content["markdown"] = map[string]string{"title": title, "text": msg}
	json, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	return SendDingDing(webhook, string(json))
}
