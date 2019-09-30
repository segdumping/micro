#### 微服务框架， 基于[go-micro](https://github.com/micro/go-micro "go-micro")，做了部分封装

##### configcenter：配置中心，存储微服务全局、通用的一些配置，支持watch
##### gateway：网关，封装了服务类型动态路由和部分包装器
##### monitor：监控，准备加入通知和自动拉起功能
##### registry：注册，封装了配置逻辑
##### log：日志，用logrus替换默认log
##### server：服务，添加通用配置
##### demo：echo、hello