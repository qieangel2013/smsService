# 说明文档
- 基于gin和gorm实现的短信基础服务，支持swagger看restful的api文档
## 部署说明
- 导入database.sql 添加一些模板
- 如果使用swagger的话，需要go get -u github.com/go-swagger/go-swagger/cmd/swagger export PATH=$GOPATH/bin:$PATH，然后cd docs 再进行 swag init
- 启动直接运行**main**或者main -config /etc/conf.yaml即可，生产环境建议使用**supervisor**部署。
- 终止运行需使用**kill -9 pid**,该指令可以实现平滑关闭，如果有运行中的任务会等待任务运行结束。

## 配置文件说明

``` yaml
server:
    listenHost: 0.0.0.0                            #http设置监听ip
    listenPort: 9510                               #http设置监听端口
    authorization: "test"                        #header验证
    env: "test"                                    # 可以设置test 、sandbox、prod 
    ipWhiteTable: ["127.0.0.1","192.168.234.1"]    #白名单绕过auth认证
    limiter:                                       
        1:                                         #设置具体类别限制发送条数
            data : ["2"]                      #103内部系统的模板编码 ShortSmsConstant::SMS_TYPE_QCODE
            msg : ""
        5: 
            data: ["x","h"]                #100,105
            msg: "您今天发送的验证码已经超过5次，不能再发送了"
    sms:                                          #设置运营商优先级发送短信，数字越小，优先级越高          
        1: "test1"
        2: "test2"
    debug: true                                   #打开debug模式，可以用于有问题定位到问题，优先级大于环境env配置
    dingding: "https://oapi.dingtalk.com/robot/send?access_token=51123" #短信运营商发送失败发送钉钉报警
component:                                      #依赖组件配置
    db:
        system:
            host: "127.0.0.1"
            port: 33060
            user: "test"
            password: "123456"
            database: "test1"
        www:
            host: "127.0.0.1"
            port: 33060
            user: "test1"
            password: "123456"
            database: "test2"
    redis:
        local:
            host: "127.0.0.1"
            port: 6379
            password: ""
            db: 0
    log:
        console: true    #配置是否控制台输出，如果设置false会重定向到日志文件里
        level: "debug"   #日志级别，建议线上设置error
        dir : "/data/log/smsService/"
```
