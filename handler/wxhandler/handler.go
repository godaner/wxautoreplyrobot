package wxhandler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/godaner/wxrobot"
	"github.com/larspensjo/config"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
	"log"
	"wxautoreplyrobot"
)

func TextHandler(msg *wxrobot.Message) error {
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
		return nil
	}
	if reply == "" {
		return nil
	}

	return wxrobot.SendMsg(msg.FromUserName, reply)
}
func ShowQRHandler(qrbyte []byte) error {
	if wxautoreplyrobot.Email == "" || wxautoreplyrobot.EmailPassword == "" || wxautoreplyrobot.EmailHost == "" || wxautoreplyrobot.EmailPort == 0 {
		return nil
	}
	////qr page////
	base64qr, _ := generateQRBase64(qrbyte, 256)
	content := `
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>wxrobot login</title>
</head>
<body>
	Please scan this qr to login wxrobot:<br/>
	<img src="` + base64qr + `"/>
</body>
</html>
`
	////qr email////

	log.Printf("Sending qr emial to %s ...... ", wxautoreplyrobot.Email)

	m := gomail.NewMessage()
	// 发件人
	m.SetAddressHeader("From", wxautoreplyrobot.Email, "wxrobot")
	// 收件人
	m.SetHeader("To", m.FormatAddress(wxautoreplyrobot.Email, "w"))
	// 主题
	m.SetHeader("Subject", "wxrobot login")
	// 发送的body体
	m.SetBody("text/html", content)

	// 发送邮件服务器、端口、发件人账号、发件人密码
	d := gomail.NewPlainDialer(wxautoreplyrobot.EmailHost, 465, wxautoreplyrobot.Email, wxautoreplyrobot.EmailPassword)
	err := d.DialAndSend(m)
	if err != nil {
		log.Printf("Send qr emial to %s err , err is : %s ! ", wxautoreplyrobot.Email, err.Error())
		return errors.Errorf("send email error")
	} else {
		log.Printf("Send qr emial to %s success ! ", wxautoreplyrobot.Email)
	}
	return nil
}
func generateQRBase64(qrbyte []byte, size int) (string, error) {
	//image to base 64
	buf := bytes.NewBuffer(qrbyte)
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
