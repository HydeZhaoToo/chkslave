chkslave
========

用于检测 mysql从库健康 并发送邮件报警

安装
========
go build chkslave.go

执行
========
./chkslave -config=slave.ini

进行后台执行

nohup ./chkslave -config=slave.ini &
