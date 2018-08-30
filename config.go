package wxautoreplyrobot

import "flag"

var TextReplyPath string
var Addr string
func init(){
	textReplyPath :=flag.String("textReplyPath","","")
	addr :=flag.String("addr",":80","")

	flag.Parse()

	TextReplyPath=*textReplyPath
	Addr=*addr
}