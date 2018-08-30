package wxautoreplyrobot

import "flag"

var TextReplyPath string
var Addr string
var Email string
func init(){
	textReplyPath :=flag.String("textReplyPath","","")
	addr :=flag.String("addr",":80","")
	email :=flag.String("email","1138829222@QQ.COM","")

	flag.Parse()

	TextReplyPath=*textReplyPath
	Addr=*addr
	Email=*email
}