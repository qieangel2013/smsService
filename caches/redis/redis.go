package redis

import (
	rcache "julive/components/cache"
	"julive/components/logger"
	"julive/tools/coding"
	"time"
)

type Cache struct {
	Key    string
	Value  string
	Client string
}

//生成key
func (c *Cache) GetKeyStr(k string) string {
	md5str, err := coding.MD5(k)
	if err != nil {
		logger.Error("生成"+k+"md5失败", err)
	}
	return "sms:" + k + ":" + md5str[8:24]
}

//设置key
func (c *Cache) SetValue(k string, v string) error {
	key := c.GetKeyStr(k)
	err := rcache.Get("local").Set(key, v, time.Hour).Err()
	if err != nil {
		logger.Error("设置缓存"+key+"失败", err)
	}
	return nil
}

//获取缓存
func (c *Cache) GetValue(k string) string {
	key := c.GetKeyStr(k)
	result, err := rcache.Get("local").Get(key).Result()
	if err != nil {
		logger.Error("获取缓存"+key+"失败", err)
	}
	return result
}

//删除缓存
func (c *Cache) DelValue(k string) error {
	key := c.GetKeyStr(k)
	err := rcache.Get("local").Del(key).Err()
	if err != nil {
		logger.Error("删除缓存"+key+"失败", err)
	}
	return nil
}
