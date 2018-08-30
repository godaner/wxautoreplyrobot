package handler

import (
	"github.com/godaner/go-route/route"
)

func Routes() {
	route.RegistRoutes(
		route.MakeAnyRoute("/reply/list",ReplyListHandler),
		route.MakeAnyRoute("/reply/delete",ReplyDeleteHandler),
		route.MakeAnyRoute("/reply/editview",ReplyEditViewHandler),
		route.MakeAnyRoute("/reply/add",ReplyAddHandler),)
}

