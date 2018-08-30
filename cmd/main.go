package main

import (
	"log"

	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/godaner/go-route/route"
	"github.com/godaner/go-util"
	"github.com/godaner/wxrobot"
	"github.com/larspensjo/config"
	"github.com/skip2/go-qrcode"
	"image/jpeg"
	"os"
	"time"
	"wxautoreplyrobot"
	"wxautoreplyrobot/handler"
	"strings"
	"net/smtp"
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

	//// web server ////
	//routes
	handler.Routes()

	//run
	go route.Start(wxautoreplyrobot.Addr)

	//// wxrobot ////
	// forever retry show login qr
	for {
		for {
			wxrobot.SetClientHandler(&wxrobot.Handler{
				TextHandler:   textHandler,
				ShowQRHandler: showQRHandler,
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

func textHandler(msg *wxrobot.Message) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	c, _ := config.ReadDefault(wxautoreplyrobot.TextReplyPath)
	//reply
	reply, err := c.String("msg", msg.Content)
	if err != nil {
		//log.Println("textHandler : get reply is err ! err is : ",err)
		return
	}
	if reply == "" {
		return
	}
	wxrobot.SendMsg(msg.FromUserName, reply)
}
func showQRHandler(qrStrP *string) {
	if wxautoreplyrobot.Email == ""{
		return
	}
	////qr page////
	base64qr, _ := generateQRBase64(*qrStrP, 120)
	content := `
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Reply List</title>
</head>
<body>
	Please scan this qr to login wx:
	<img src="` + base64qr + `"/>
</body>
</html>
`
	////qr email////

	err := SendToMail("godanermail@gmail.com", "ZK19951217.", "smtp.gmail.com:465", wxautoreplyrobot.Email, "wxrobot login", content, "html")


	if  err != nil {
		log.Printf("Send emial to %s err , err is : %s ! ",wxautoreplyrobot.Email,err.Error())
	}else {
		log.Printf("Send emial to %s success ! ",wxautoreplyrobot.Email)
	}
}
func generateQRBase64(qrCode string, size int) (string, error) {
	if size <= 0 {
		size = 250
	}
	q, err := qrcode.New(qrCode, qrcode.Medium)
	if err != nil {
		return "", err
	}
	image := q.Image(size)
	if err != nil {
		return "", err
	}
	//image to base 64
	emptyBuff := bytes.NewBuffer(nil)
	jpeg.Encode(emptyBuff, image, nil)

	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(emptyBuff.Bytes()), nil
}
func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}