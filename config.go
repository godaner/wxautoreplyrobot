package wxautoreplyrobot

import "flag"

var TextReplyPath string
var Addr string
var Email string
var EmailPassword string
var EmailHost string
var EmailPort int
func init(){
	textReplyPath :=flag.String("textReplyPath","","")
	addr :=flag.String("addr",":80","")
	email :=flag.String("email","","")
	emailPassword :=flag.String("emailPassword","","")
	emailHost :=flag.String("emailHost","","")
	emailPort :=flag.Int64("emailPort",0,"")

	flag.Parse()

	TextReplyPath=*textReplyPath
	Addr=*addr
	Email=*email
	EmailPassword=*emailPassword
	EmailHost=*emailHost
	EmailPort=(int)(*emailPort)
}