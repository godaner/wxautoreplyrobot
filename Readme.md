## 微信自动回复机器人-wxautoreplyrobot
## Usage：

### 	download：

​		

```
cd ${GOPATH}/src
git clone https://github.com/godaner/wxautoreplyrobot.git
```

### 	run demo：

#### linux：

```
cd ${GOPATH}/src/wxautoreplyrobot/cmd
go run main.go -textReplyPath ${GOPATH}/src/wxautoreplyrobot/textreply.cfg -addr :8887 -email  1138829***@.com
if you wanna run wxrobot in background:
	nohup go run main.go -textReplyPath ${GOPATH}/src/wxautoreplyrobot/textreply.cfg -addr :8887 -email 1138829***@.com>wxautoreplyrobot.log 2>&1 & 
if you wanna see log:
	tail -f wxautoreplyrobot.log
```

#### docker：

```
docker pull godaber/wxautoreplyrobot
docker run -p 8887:8887 -e email="38829***@.com" godaner/wxautoreplyrobot
if you wanna run wxrobot in background:
	docker run -d -p 8887:8887 -e email="1138829***@.com" --name wxautoreplyrobot godaner/wxautoreplyrobot
if you wanna see log:
	docker logs -f wxautoreplyrobot
```



### 	demo result：

#### console:

```
2018/08/29 09:44:24 wx.go:104: Please open link in browser: https://login.weixin.qq.com/qrcode/IesWCyGxZg==
2018/08/29 09:44:49 wx.go:129: login timeout, reconnecting...
2018/08/29 09:45:09 wx.go:133: scan success, please confirm login on your phone
2018/08/29 09:45:12 wx.go:136: login success
2018/08/29 09:45:19 wx.go:305: update 141 contacts
2018/08/29 09:45:19 wx.go:313: @c458675e3c522f5f0bc436f0a861ca16 => 微信安全中心
2018/08/29 09:45:19 wx.go:313: @c8e81be227ce9428490833eb837a5ee81f3b9b6beed19809cdd4a53976fd4104 => 快乐人生
......
```

#### web:

​	visit http://127.0.0.1:8887/reply/list to manage your reply msg.

### about textreply.cfg：

#### content：

```
[msg]
hello: hello , i am godaner !
```

#### explain：

```
*.you don't need create this file , program will create it auto.
*.it is a demo reply config
*.don't change node [msg].
*.if anyone send "hello" to you , your wxrobot will response "hello , i am godaner !" to him auto.
```

