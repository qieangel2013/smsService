# 基础类库Go基础类库

## 目录结构说明

``` bash
├── components          //基础组件(需要注册才能使用)
│   ├── cache
│   ├── db
│   ├── logger
│   └── queue
├── go.mod              //mod文件(声明依赖）
├── go.sum
├── middlewares         //中间件(Gin框架专用)
│   ├── logger.go
│   └── recover.go
├── models              //模型代码(跨项目共用的model才放这里）
│   └── base.go
├── services            //业务逻辑代码(跨项目共用的services才放这里)
├── structs             //通用结构体定义
│   └── base.go
├── tests               //类库单元测试文件
├── thirdpartys         //第三方服务封装
│   ├── aliyun
│   └── dingding
└── tools               //工具函数
    ├── coding
    ├── debuger
    ├── file
    └── http
```
