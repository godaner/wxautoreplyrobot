## 微信自动回复机器人-wxautoreplyrobot
## 用法：

### 	下载：

​		

```
cd ${GOPATH}/src
git clone https://github.com/godaner/wxautoreplyrobot.git
```

### 	运行：

#### 	参数：

```
textReplyPath：回复词条储存位置

email：邮箱地址，如1138829***@qq.com

emailPassword：邮箱密码/授权码

emailHost：邮箱服务器主机，如smtp.qq.com

emailPort：邮箱主机端口，如465=5

refreshWhenError：当程序出现错误时，重新打印登录二维码的间隔时间

注：其中参数email、emailPassword、emailHost、emailPort在你不需要邮件通知的情况下可以不用填写，refreshWhenError默认120分钟，也可不用填写
```



#### 在linux上运行：

```
cd ${GOPATH}/src/wxautoreplyrobot/cmd
go run main.go -textReplyPath ${GOPATH}/src/wxautoreplyrobot/textreply.cfg -addr :8887 -email  1138829***@qq.com -emailPassword 123 -emailHost smtp.qq.com -emailPort 465 -refreshWhenError 120 godaner/wxautoreplyrobot
如果你想后台运行:
	nohup go -textReplyPath ${GOPATH}/src/wxautoreplyrobot/textreply.cfg -addr :8887 -email  1138829***@qq.com -emailPassword 123 -emailHost smtp.qq.com -emailPort 465 -refreshWhenError 120>wxautoreplyrobot.log 2>&1 & 
如果你想看日志:
	tail -f wxautoreplyrobot.log
```

#### 在docker上运行：

```
docker pull godaber/wxautoreplyrobot
docker run -p 8887:8887 -e email="1138829222@qq.com" -e emailPassword="nofuhedsnzduibeb" -e emailHost="smtp.qq.com" -e emailPort=465 -e refreshWhenError=120 --name wxautoreplyrobot godaner/wxautoreplyrobot
如果你想后台运行:
	docker run -d -p 8887:8887 -e email="1138829222@qq.com" -e emailPassword="nofuhedsnzduibeb" -e emailHost="smtp.qq.com" -e emailPort=465 -e refreshWhenError=120 --name wxautoreplyrobot godaner/wxautoreplyrobot
如果你想看日志:
	docker logs -f wxautoreplyrobot
```



### 	运行结果：

#### 控制台输出:

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

#### 访问程序管理页面:

​	http://127.0.0.1:8887/reply/list

用于管理恢复词条

### 关于 textreply.cfg：

#### 内容：

```
[msg]
hello: hello , i am godaner !
```

#### 解释：

```
*.你不需要手动创建此文件，系统将自己创建。
*.这是一个类似于数据库的文件，用于储存恢复词条键值对。
*.不要尝试手动修改该文件.
*.例如：如果在微信中有任何人发送“hello”给你，机器人将自动回复“hello , i am godaner !”
```

