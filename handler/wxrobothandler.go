package handler

import (
	"github.com/godaner/go-route/route"
	"github.com/larspensjo/config"
	"wxautoreplyrobot"
	"fmt"
	"os"
)

func ReplyListHandler(response route.RouteResponse, request route.RouteRequest) {

	response.WriteString(returnListPage())
}
func ReplyAddHandler(response route.RouteResponse, request route.RouteRequest) {

	//param
	key:=request.Params["key"].(string)
	val:=request.Params["val"].(string)

	c, _ := config.ReadDefault(wxautoreplyrobot.TextReplyPath)
	c.AddOption("msg",key,val)
	c.WriteFile(wxautoreplyrobot.TextReplyPath,os.ModeDevice,"")

	response.WriteString(returnListPage())
}


func returnListPage() string{

	c, _ := config.ReadDefault(wxautoreplyrobot.TextReplyPath)
	options,_:=c.Options("msg")

	var content string
	for _,opt:=range options{
		val,_:=c.String("msg",opt)
		content+=fmt.Sprintf("<tr><td> %s</td>  <td> %s </td></tr>",opt,val)
	}


	listPage:=`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Reply List</title>
</head>
<body>
	<form action="/reply/add">
		<input name="key"/>
		<input name="val"/>
		<input type="submit"/>
	</form>
	<table border='1px'>
		<tr>
			<td>
			需回复字段
			</td>
			<td>
			回复内容
			</td>
		</tr>
			`+content+`
	</table>
</body>
</html>
`
return listPage
}