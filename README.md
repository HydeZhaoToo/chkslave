chkslave
========

用于检测 mysql从库健康 并发送邮件报警

依赖
========
安装golang 
详情http://www.alaiblog.com/program/step-by-step-install-golang-go.html

go get github.com/Unknwon/goconfig

go get github.com/go-sql-driver/mysql

安装
========
go build chkslave.go

执行
========
./chkslave -config=slave.ini

进行后台执行

nohup ./chkslave -config=slave.ini &



```sequence
participant 通知服务
participant 消费者
participant 平台
participant 代理
participant 生产者


消费者-平台:API请求
平台-代理:加入header
代理-生产者:代理请求
生产者--代理:应答数据
代理--平台:应答数据
平台--消费者:应答数据
生产者-代理:通知请求(鉴权限流统计)
代理-平台:代理通知请求
平台-->>通知服务:通知请求
通知服务--消费者:通知
```
