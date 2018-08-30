package handler

import (
	"fmt"
	"github.com/godaner/go-route/route"
	"github.com/larspensjo/config"
	"os"
	"wxautoreplyrobot"
	"regexp"
	"strings"
)

func ReplyListHandler(response route.RouteResponse, request route.RouteRequest) {

	response.WriteString(buildListPage(nil))
}
func ReplyAddHandler(response route.RouteResponse, request route.RouteRequest) {

	//param
	key := request.Params["key"].(string)
	val := request.Params["val"].(string)

	if key == "" || val == "" {

		response.WriteString(buildListPage(nil))
		return
	}

	c, _ := config.ReadDefault(wxautoreplyrobot.TextReplyPath)
	c.AddOption("msg", key, val)
	c.WriteFile(wxautoreplyrobot.TextReplyPath, os.ModeDevice, "")

	response.WriteString(buildListPage(nil))
}
func ReplyDeleteHandler(response route.RouteResponse, request route.RouteRequest) {

	//param
	key := request.Params["key"].(string)

	c, _ := config.ReadDefault(wxautoreplyrobot.TextReplyPath)
	c.RemoveOption("msg", key)
	c.WriteFile(wxautoreplyrobot.TextReplyPath, os.ModeDevice, "")

	response.WriteString(buildListPage(nil))
}
func ReplyEditViewHandler(response route.RouteResponse, request route.RouteRequest) {

	//param
	key := request.Params["key"].(string)

	c, _ := config.ReadDefault(wxautoreplyrobot.TextReplyPath)
	val, _ := c.String("msg", key)

	response.WriteString(buildListPage(map[string]string{
		"key": key,
		"val": val,
	}))
}
func buildListPage(params map[string]string) string {

	c, _ := config.ReadDefault(wxautoreplyrobot.TextReplyPath)
	options, _ := c.Options("msg")

	var content string
	for _, opt := range options {
		val, _ := c.String("msg", opt)
		content += fmt.Sprintf("<tr><td> %s</td>  <td> %s </td><td> <a href='/reply/delete?key="+opt+"'>删除</a> # <a href='/reply/editview?key="+opt+"'>编辑</a></td></tr>", opt, val)
	}

	listPage := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Reply List</title>
</head>
<body>
	<form action="/reply/add">
		<input name="key" value='${key}'/>
		<input name="val" value='${val}'/>
		<input type="submit" value='添加/更新'/>
	</form>
	<table border='1px'>
		<tr>
			<td>
			需回复字段
			</td>
			<td>
			回复内容
			</td>
			<td>
			操作
			</td>
		</tr>
			` + content + `
	</table>
</body>
</html>
`

	return parsePageParam(listPage, params)
}
func parsePageParam(page string, params map[string]string) string {
	for key := range params {
		page = strings.Replace(page,"${"+key+"}",params[key],-1)
	}
	reg, _ := regexp.Compile(`\$\{[\w\d]*\}`)
	page = reg.ReplaceAllString(page, "")
	return page

}
