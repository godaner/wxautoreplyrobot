package wxhandler

import (
	"fmt"
	"github.com/larspensjo/config"
	"wxautoreplyrobot"
	"log"
	"github.com/skip2/go-qrcode"
	"bytes"
	"image/jpeg"
	"encoding/base64"
	"strings"
	"net/smtp"
	"github.com/godaner/wxrobot"
	"gopkg.in/gomail.v2"
)

func TextHandler(msg *wxrobot.Message) {
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
func ShowQRHandler(qrStrP *string) {
	if wxautoreplyrobot.Email == ""{
		return
	}
	////qr page////
	base64qr, _ := generateQRBase64(*qrStrP, 256)
	content := `
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>wxrobot login</title>
</head>
<body>
	Please scan this qr to login wxhandler:<br/>
	<img src="` + base64qr + `"/>
</body>
</html>
`
	////qr email////

	log.Printf("Sending qr emial to %s ...... ",wxautoreplyrobot.Email)

	m := gomail.NewMessage()
	// 发件人
	m.SetAddressHeader("From", "1138829222@QQ.COM", "wxrobot")
	// 收件人
	m.SetHeader("To", m.FormatAddress(wxautoreplyrobot.Email, "w"))
	// 主题
	m.SetHeader("Subject", "wxrobot login")
	// 发送的body体
	m.SetBody("text/html", content)

	// 发送邮件服务器、端口、发件人账号、发件人密码
	d := gomail.NewPlainDialer("smtp.qq.com", 465, "1138829222@QQ.COM", "nofuhedsnzduibeb")
	err := d.DialAndSend(m)
	if  err != nil {
		log.Printf("Send qr emial to %s err , err is : %s ! ",wxautoreplyrobot.Email,err.Error())
	}else {
		log.Printf("Send qr emial to %s success ! ",wxautoreplyrobot.Email)
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
func sendToMail(user, password, host, to, subject, body, mailtype string) error {
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
