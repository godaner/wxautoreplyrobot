package wxautoreplyrobot

import "flag"

var TextReplyPath string
var Addr string
var Email string
var EmailPassword string
var EmailHost string
var EmailPort int
var RefreshWhenError int64
func init(){
	textReplyPath :=flag.String("textReplyPath","","save reply word file path , like e:/textreply.cfg . ")
	addr :=flag.String("addr",":80","web server address . ")
	email :=flag.String("email","","if you wanna use email to send qr , please write your email , like 1138829***@qq.com . ")
	emailPassword :=flag.String("emailPassword","","if you wanna use email to send qr , please write your email password , like shfiuawhojfjha . ")
	emailHost :=flag.String("emailHost","","if you wanna use email to send qr , please write email server host , like smtp.qq.com . ")
	emailPort :=flag.Int64("emailPort",0,"if you wanna use email to send qr , please write email server port , like 465 . ")
	refreshWhenError :=flag.Int64("refreshWhenError",120,"when program occur error , program will wait some time , then print qr again . ")

	flag.Parse()

	TextReplyPath=*textReplyPath
	Addr=*addr
	Email=*email
	EmailPassword=*emailPassword
	EmailHost=*emailHost
	EmailPort=(int)(*emailPort)
	RefreshWhenError=*refreshWhenError
}