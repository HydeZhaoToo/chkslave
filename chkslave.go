/*
用于检测mysql从库健康状态 邮件报警
*/
package main

import (
	"database/sql"
	"errors"
	
	"flag"
	"fmt"
	"github.com/Unknwon/goconfig"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"
)

var (
	Slave_IO_State                string
	Master_Host                   string
	Master_User                   string
	Master_Port                   string
	Connect_Retry                 string
	Master_Log_File               string
	Read_Master_Log_Pos           string
	Relay_Log_File                string
	Relay_Log_Pos                 string
	Relay_Master_Log_File         string
	Slave_IO_Running              string
	Slave_SQL_Running             string
	Replicate_Do_DB               string
	Replicate_Ignore_DB           string
	Replicate_Do_Table            string
	Replicate_Ignore_Table        string
	Replicate_Wild_Do_Table       string
	Replicate_Wild_Ignore_Table   string
	Last_Errno                    string
	Last_Error                    string
	Skip_Counter                  string
	Exec_Master_Log_Pos           string
	Relay_Log_Space               string
	Until_Condition               string
	Until_Log_File                string
	Until_Log_Pos                 string
	Master_SSL_Allowed            string
	Master_SSL_CA_File            string
	Master_SSL_CA_Path            string
	Master_SSL_Cert               string
	Master_SSL_Cipher             string
	Master_SSL_Key                string
	Seconds_Behind_Master         string
	Master_SSL_Verify_Server_Cert string
	Last_IO_Errno                 string
	Last_IO_Error                 string
	Last_SQL_Errno                string
	Last_SQL_Error                string
	Replicate_Ignore_Server_Ids   string
	Master_Server_Id              string
)

//默认设置
type chkset struct {
	rate   int
	user   string
	passwd string
	port   int
	host   string
}

//mail配置
type mailini struct {
	user        string
	passwd      string
	smtpaddress string
	maillist    string
	smtpport    int
	content     string
}

var (
	configfile string
	loger      *log.Logger
	logurl     string
)

//参数摄取
func init() {
	flag.StringVar(&configfile, "config", "", "配置文件路径")
	flag.StringVar(&logurl, "log", "check_slave.log", "配置文件路径")
	flag.Parse()
}

//初始化mail结构体配置
func (m mailini) newmailini(g *goconfig.ConfigFile) *mailini {
	var err error
	m.maillist, err = g.GetValue("mail", "receive")
	chkerr(err)
	m.user, err = g.GetValue("mail", "mailuser")
	chkerr(err)
	m.passwd, err = g.GetValue("mail", "mailpasswd")
	chkerr(err)
	m.smtpaddress, err = g.GetValue("mail", "smtpaddress")
	chkerr(err)
	m.smtpport = g.MustInt("mail", "smtpport", 25)
	m.content = ""
	return &m
}

//初始化检测参数
func (c chkset) newchkset(g *goconfig.ConfigFile) *chkset {
	var err error
	c.rate = g.MustInt("", "rate", 60)
	c.host, err = g.GetValue("", "host")
	chkerr(err)
	c.port = g.MustInt("", "port", 3306)
	c.user, err = g.GetValue("", "user")
	chkerr(err)
	c.passwd, err = g.GetValue("", "passwd")
	chkerr(err)
	return &c
}

func main() {
	logfile, err := os.OpenFile(logurl, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	chkerr(err)
	defer logfile.Close()
	loger = log.New(logfile, "", log.Ldate|log.Ltime)
	config, err := goconfig.LoadConfigFile(configfile)
	chkerr(err)
	var c chkset
	var m mailini
	s := c.newchkset(config)
	n := m.newmailini(config)

	for {
		if err := chkslave(s); err != nil {
			fmt.Println(err)
			n.content = fmt.Sprintf("IP地址:%s的mysql slave检测报错:%v", s.host, err)
			sendmail(n)
		}
		time.Sleep(time.Duration(s.rate) * time.Second)
	}
}

//检测从库错误
func chkslave(c *chkset) error {
	Dsource := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql", c.user, c.passwd, c.host, c.port)
	db, err := sql.Open("mysql", Dsource)
	defer db.Close()
	chkerr(err)
	result, err2 := db.Query("show slave status")
	chkerr(err2)
	if result.Next() {
		result.Scan(&Slave_IO_State, &Master_Host, &Master_User, &Master_Port, &Connect_Retry, &Master_Log_File, &Read_Master_Log_Pos, &Relay_Log_File, &Relay_Log_Pos, &Relay_Master_Log_File, &Slave_IO_Running, &Slave_SQL_Running, &Replicate_Do_DB, &Replicate_Ignore_DB, &Replicate_Do_Table, &Replicate_Ignore_Table, &Replicate_Wild_Do_Table, &Replicate_Wild_Ignore_Table, &Last_Errno, &Last_Error, &Skip_Counter, &Exec_Master_Log_Pos, &Relay_Log_Space, &Until_Condition, &Until_Log_File, &Until_Log_Pos, &Master_SSL_Allowed, &Master_SSL_CA_File, &Master_SSL_CA_Path, &Master_SSL_Cert, &Master_SSL_Cipher, &Master_SSL_Key, &Seconds_Behind_Master, &Master_SSL_Verify_Server_Cert, &Last_IO_Errno, &Last_IO_Error, &Last_SQL_Errno, &Last_SQL_Error, &Replicate_Ignore_Server_Ids, &Master_Server_Id)
		if Seconds_Behind_Master != "0" || Slave_IO_Running != "Yes" || Slave_SQL_Running != "Yes" {
			return errors.New(fmt.Sprintf("%s %s %v", Last_SQL_Error, Last_Error,Seconds_Behind_Master))
		} else {
			return nil
		}
	} else {
		return errors.New("未查找到slave状态")
	}
}

//发送邮件
func sendmail(m *mailini) error {
	sub := "subject:  mysql从库报警通知邮件  \r\n\r\n"
	mailList := strings.Split(m.maillist, ",")
	auth := smtp.PlainAuth("", m.user, m.passwd, m.smtpaddress)
	au := fmt.Sprintf("%s:%d", m.smtpaddress, m.smtpport)
	err := smtp.SendMail(au, auth, m.user, mailList, []byte(sub+m.content))
	if err != nil {
		return err
	} else {
		return nil
	}
}

//错误输出
func chkerr(err error) {
	if err != nil {
		fmt.Println(err)
		loger.Fatalln(err)
	}
}
