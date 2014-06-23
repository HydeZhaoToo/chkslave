chkslave
========

用于检测 mysql从库健康 并发送邮件报警

依赖
========
安装golang 
详情http://www.alaiblog.com/program/step-by-step-install-golang-go.html

安装 github.com/Unknwon/goconfig

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
