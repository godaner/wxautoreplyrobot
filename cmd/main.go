package main

import (
	"log"

	"github.com/godaner/go-route/route"
	"github.com/godaner/go-util"
	"github.com/godaner/wxrobot"
	"github.com/larspensjo/config"
	"os"
	"time"
	"wxautoreplyrobot"
	"wxautoreplyrobot/handler/webhandler"
	"wxautoreplyrobot/handler/wxhandler"
)

const (
	TIME_LAYOUT = "2006-01-02 15:04:05"
)

func init() {
	//check cfg
	if !go_util.FileExists(wxautoreplyrobot.TextReplyPath) {
		os.Create(wxautoreplyrobot.TextReplyPath)
	}
	c, err := config.ReadDefault(wxautoreplyrobot.TextReplyPath)
	if err != nil {
		log.Println("init err : ", err)
	}
	if !c.HasSection("msg") {
		c.AddOption("msg", "hello", "i am wxautoreplyrobot! birthday is : "+time.Now().Format(TIME_LAYOUT))
	}
	c.WriteFile(wxautoreplyrobot.TextReplyPath, os.ModeDevice, "")
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//// webhandler server ////
	//routes
	webhandler.Routes()

	//run
	go route.Start(wxautoreplyrobot.Addr)

	//// wxrobot ////
	// forever retry show login qr
	for {
		for {
			wxrobot.SetClientHandler(&wxrobot.Handler{
				TextHandler:   wxhandler.TextHandler,
				ShowQRHandler: wxhandler.ShowQRHandler,
			})
			err := wxrobot.StartClient() //will be block
			if err != nil {
				log.Printf("wxrobot client err , err is : %s , time is %s !", err.Error(), time.Now().Format(TIME_LAYOUT))
				continue
			}
			break
		}
	}

}
