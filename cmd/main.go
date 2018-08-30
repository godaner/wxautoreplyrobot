package main

import (
	"log"

	"github.com/godaner/wxrobot"
	"github.com/larspensjo/config"
	"fmt"
	"wxautoreplyrobot"
	"github.com/godaner/go-util"
	"os"
	"github.com/godaner/go-route/route"
	"wxautoreplyrobot/handler"
	"time"
)

const (
	TIME_LAYOUT="2006-01-02 15:04:05"
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
		c.AddOption("msg","hello","i am wxautoreplyrobot! birthday is : "+time.Now().Format(TIME_LAYOUT))
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

	isLogin := false
	//// wxrobot ////
	for{
		isLogin = false
		for !isLogin{//refresh qr
			err0 := wxrobot.Init(&wxrobot.MessageHandler{
				TextHandler: textHandler,
			})
			if err0 != nil {
				log.Printf(err0.Error())
				time.Sleep(time.Minute*time.Duration(10))
				continue
			}
			isLogin=true
		}

		//have login success
		err:=wxrobot.Listening()

		if err!=nil{
			log.Printf("wxrobot litening err , err is : %s , time is %s !",err.Error(),time.Now().Format(TIME_LAYOUT))
		}
	}

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
