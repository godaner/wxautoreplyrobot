package main

import (
	"log"

	"github.com/godaner/wxrobot"
	"github.com/larspensjo/config"
	"fmt"
	"wxautoreplyrobot"
	"go-util"
	"os"
	"github.com/godaner/go-route/route"
	"wxautoreplyrobot/handler"
	"time"
)

func init(){
	//check
	if !go_util.FileExists(wxautoreplyrobot.TextReplyPath) {
		os.Create(wxautoreplyrobot.TextReplyPath)
	}
	c, err := config.ReadDefault(wxautoreplyrobot.TextReplyPath)
	if err!=nil {
		log.Println("init err : ",err)
	}
	if !c.HasSection("msg"){
		timeStr := time.Now().Format("2006-01-02 15:04:05")
		c.AddOption("msg","hello","i am wxautoreplyrobot! birthday is : "+timeStr)
	}
	c.WriteFile(wxautoreplyrobot.TextReplyPath,os.ModeDevice,"")
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//// web server ////
	//routes
	handler.Routes()

	//run
	go route.Start(wxautoreplyrobot.Addr)

	//// wxrobot ////
	err0 := wxrobot.Init(&wxrobot.MessageHandler{
		TextHandler: textHandler,
	})
	if err0 != nil {
		log.Fatal(err0.Error())
	}

	wxrobot.Listening()


}


func textHandler(msg *wxrobot.Message){
	defer func() {
		if err := recover();err != nil {
			fmt.Println(err)
		}
	}()
	c, _ := config.ReadDefault(wxautoreplyrobot.TextReplyPath)
	//reply
	reply,err:=c.String("msg", msg.Content)
	if err!=nil {
		//log.Println("textHandler : get reply is err ! err is : ",err)
		return
	}
	if reply==""{
		return
	}
	wxrobot.SendMsg(msg.FromUserName,reply)
}
