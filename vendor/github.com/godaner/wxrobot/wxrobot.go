package wxrobot

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"bytes"
	"os"
	"github.com/546669204/golang-http-do"
	"github.com/tuotoo/qrcode"
	"github.com/mdp/qrterminal"
)

var (
	LoginUri                   = "https://login.weixin.qq.com"
	ErrUnknow                  = errors.New("Unknow Error")
	ErrUserNotExists           = errors.New("Error User Not Exist")
	ErrNotLogin                = errors.New("Not Login")
	ErrLoginTimeout            = errors.New("Login Timeout")
	ErrWaitingForConfirm       = errors.New("Waiting For Confirm")
	TEXT_MSG             int64 = 1
	IMG_MSG              int64 = 3
	VOICE_MSG            int64 = 34
	FACE_0_MSG           int64 = 43
	FACE_1_MSG           int64 = 47
	LINK_MSG             int64 = 49
	ENTER_CHAT_MSG       int64 = 51

	MSG_TYPE_MAP = map[int64]string{
		TEXT_MSG:"TEXT_MSG",
		IMG_MSG:"IMG_MSG",
		VOICE_MSG:"VOICE_MSG",
		FACE_0_MSG:"FACE_0_MSG",
		FACE_1_MSG:"FACE_1_MSG",
		LINK_MSG:"LINK_MSG",
		ENTER_CHAT_MSG:"ENTER_CHAT_MSG",
	}
)

type WXRobot struct {
	httpClient     *Client
	secret         *wxSecret
	baseRequest    *BaseRequest
	user           *User
	contacts       map[string]*User
	messageHandler *MessageHandler
}

func NewWeixin(messageHandler *MessageHandler) *WXRobot {
	return &WXRobot{
		httpClient:     NewClient(),
		secret:         &wxSecret{},
		baseRequest:    &BaseRequest{},
		user:           &User{},
		messageHandler: messageHandler,
	}
}

func (wxrobot *WXRobot) GetUser(userName string) (*User, error) {
	u, ok := wxrobot.contacts[userName]
	if ok {
		return u, nil
	} else {
		return nil, ErrUserNotExists
	}
}

func (wxrobot *WXRobot) GetUserName(userName string) string {
	u, err := wxrobot.GetUser(userName)
	if err != nil {
		return "[myself]"
	}
	if u.RemarkName != "" {
		return u.RemarkName
	} else {
		return u.NickName
	}
}

func (wxrobot *WXRobot) getUuid() (string, error) {
	values := &url.Values{}
	values.Set("appid", "wx782c26e4c19acffb")
	values.Set("fun", "new")
	values.Set("lang", "zh_CN")
	values.Set("_", TimestampStr())
	uri := fmt.Sprintf("%s/jslogin", LoginUri)
	b, err := wxrobot.httpClient.Get(uri, values)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`"([\S]+)"`)
	find := re.FindStringSubmatch(string(b))
	if len(find) > 1 {
		return find[1], nil
	} else {
		return "", fmt.Errorf("get uuid error, response: %s", b)
	}
}

func (wxrobot *WXRobot) ShowQRcodeUrl(uuid string) error {
	//qr url
	qrStr := fmt.Sprintf("%s/qrcode/%s", LoginUri, uuid)

	log.Println("Please open link in browser: " + qrStr," , or you can scan this qr :")
	//qr img
	qrStrP:=&qrStr
	op := httpdo.Default()
	op.Url = `https://login.weixin.qq.com/qrcode/` + uuid
	httpbyte, err := httpdo.HttpDo(op)
	if err != nil {
		log.Println(err)
		return errors.New("get qr byte err")
	}
	M, err := qrcode.Decode(bytes.NewReader(httpbyte))
	if qrStrP != nil {
		qrterminal.GenerateHalfBlock(M.Content, qrterminal.L, os.Stdout)
	} else {
		*qrStrP = M.Content
	}

	return nil
}

func (wxrobot *WXRobot) WaitingForLoginConfirm(uuid string) (string, error) {
	re := regexp.MustCompile(`window.code=([0-9]*);`)
	tip := "1"
	for {
		values := &url.Values{}
		values.Set("uuid", uuid)
		values.Set("tip", tip)
		values.Set("_", TimestampStr())
		b, err := wxrobot.httpClient.Get("https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login", values)
		if err != nil {
			log.Printf("HTTP GET err: %s", err.Error())
			return "", err
		}
		s := string(b)
		codes := re.FindStringSubmatch(s)
		if len(codes) == 0 {
			log.Printf("find window.code failed, origin response: %s\n", s)
			return "", ErrUnknow
		} else {
			code := codes[1]
			if code == "408" {
				log.Println("login timeout, reconnecting...")
				// }else if code == "400" {
				// 	log.Println("login timeout, need refresh qrcode")
			} else if code == "201" {
				log.Println("scan success, please confirm login on your phone")
				tip = "0"
			} else if code == "200" {
				log.Println("login success")
				re := regexp.MustCompile(`window\.redirect_uri="(.*?)";`)
				us := re.FindStringSubmatch(s)
				if len(us) == 0 {
					log.Println(s)
					return "", errors.New("find redirect uri failed")
				}
				return us[1], nil
			} else {
				log.Printf("unknow window.code %s\n", code)
				return "", ErrUnknow
			}
		}
	}
	return "", nil
}

func findTicket(s string) (*ticket, error) {
	re := regexp.MustCompile(`window\.redirect_uri="(.*?)";`)
	us := re.FindStringSubmatch(s)
	if len(us) == 0 {
		log.Println(s)
		return nil, errors.New("find redirect_uri error")
	}
	u, err := url.Parse(us[1])
	if err != nil {
		return nil, err
	}
	q := u.Query()
	return &ticket{
		Ticket: q.Get("ticket"),
		Scan:   q.Get("scan"),
	}, nil
}

func (wxrobot *WXRobot) NewLoginPage(newLoginUri string) error {
	b, err := wxrobot.httpClient.Get(newLoginUri+"&fun=new", nil)
	if err != nil {
		log.Printf("HTTP GET err: %s", err.Error())
		return err
	}
	err = xml.Unmarshal(b, wxrobot.secret)
	if err != nil {
		log.Printf("parse wxSecret from xml failed: %v", err)
		return err
	}
	if wxrobot.secret.Code == "0" {
		u, _ := url.Parse(newLoginUri)
		wxrobot.secret.BaseUri = newLoginUri[:strings.LastIndex(newLoginUri, "/")]
		wxrobot.secret.Host = u.Host
		wxrobot.secret.DeviceID = "e" + RandNumbers(15)
		return nil
	} else {
		return errors.New("Get wxSecret Error")
	}

}

func (wxrobot *WXRobot) Init() error {
	values := &url.Values{}
	values.Set("r", TimestampStr())
	values.Set("lang", "en_US")
	values.Set("pass_ticket", wxrobot.secret.PassTicket)
	url := fmt.Sprintf("%s/webwxinit?%s", wxrobot.secret.BaseUri, values.Encode())
	wxrobot.baseRequest = &BaseRequest{
		Uin:      wxrobot.secret.Uin,
		Sid:      wxrobot.secret.Sid,
		Skey:     wxrobot.secret.Skey,
		DeviceID: wxrobot.secret.DeviceID,
	}
	b, err := wxrobot.httpClient.PostJson(url, map[string]interface{}{
		"BaseRequest": wxrobot.baseRequest,
	})
	if err != nil {
		log.Printf("HTTP GET err: %s", err.Error())
		return err
	}
	var r InitResponse
	err = json.Unmarshal(b, &r)
	if err != nil {
		return err
	}
	if r.BaseResponse.Ret == 0 {
		wxrobot.user = r.User
		wxrobot.updateSyncKey(r.SyncKey)
		return nil
	}
	return fmt.Errorf("Init error: %+v", r.BaseResponse)
}

func (wxrobot *WXRobot) updateSyncKey(s *SyncKey) {
	wxrobot.secret.SyncKey = s
	syncKeys := make([]string, s.Count)
	for n, k := range s.List {
		syncKeys[n] = fmt.Sprintf("%d_%d", k.Key, k.Val)
	}
	wxrobot.secret.SyncKeyStr = strings.Join(syncKeys, "|")
}

func (wxrobot *WXRobot) GetNewLoginUrl() (string, error) {
	uuid, err := wxrobot.getUuid()
	if err != nil {
		return "", err
	}
	err = wxrobot.ShowQRcodeUrl(uuid)
	if err != nil {
		return "", err
	}
	newLoginUri, err := wxrobot.WaitingForLoginConfirm(uuid)
	if err != nil {
		return "", err
	}
	return newLoginUri, nil
}

type syncStatus struct {
	Retcode  string
	Selector string
}

func (wxrobot *WXRobot) StatusNotify() error {
	values := &url.Values{}
	values.Set("lang", "zh_CN")
	values.Set("pass_ticket", wxrobot.secret.PassTicket)
	url := fmt.Sprintf("%s/webwxstatusnotify?%s", wxrobot.secret.BaseUri, values.Encode())
	b, err := wxrobot.httpClient.PostJson(url, map[string]interface{}{
		"BaseRequest":  wxrobot.baseRequest,
		"code":         3,
		"FromUserName": wxrobot.user.UserName,
		"ToUserName":   wxrobot.user.UserName,
		"ClientMsgId":  TimestampMicroSecond(),
	})
	if err != nil {
		return err
	}
	return wxrobot.CheckCode(b, "Status Notify error")
}

func (wxrobot *WXRobot) CheckCode(b []byte, errmsg string) error {
	var r InitResponse
	err := json.Unmarshal(b, &r)
	if err != nil {
		return err
	}
	if r.BaseResponse.Ret != 0 {
		return errors.New("Status Notify error")
	}
	return nil
}

func (wxrobot *WXRobot) GetContacts() error {
	values := &url.Values{}
	values.Set("seq", "0")
	values.Set("pass_ticket", wxrobot.secret.PassTicket)
	values.Set("skey", wxrobot.secret.Skey)
	values.Set("r", TimestampStr())
	url := fmt.Sprintf("%s/webwxgetcontact?%s", wxrobot.secret.BaseUri, values.Encode())
	b, err := wxrobot.httpClient.PostJson(url, map[string]interface{}{})
	if err != nil {
		return err
	}
	var r ContactResponse
	err = json.Unmarshal(b, &r)
	if err != nil {
		return err
	}
	if r.BaseResponse.Ret != 0 {
		return errors.New("Get Contacts error")
	}
	log.Printf("update %d contacts", r.MemberCount)
	wxrobot.contacts = make(map[string]*User, r.MemberCount)
	return wxrobot.updateContacts(r.MemberList)
}

func (wxrobot *WXRobot) updateContacts(us []*User) error {
	for _, u := range us {
		wxrobot.contacts[u.UserName] = u
		log.Printf("%s => %s", u.UserName, u.NickName)
	}
	return nil
}

func (wxrobot *WXRobot) TestSyncCheck() error {
	for _, h := range []string{"webpush.", "webpush2."} {
		wxrobot.secret.PushHost = h + wxrobot.secret.Host
		syncStatus, err := wxrobot.SyncCheck()
		if err == nil {
			if syncStatus.Retcode == "0" {
				return nil
			}
		}
	}
	return errors.New("Test SyncCheck error")
}

func (wxrobot *WXRobot) SyncCheck() (*syncStatus, error) {
	uri := fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin/synccheck", wxrobot.secret.PushHost)
	values := &url.Values{}
	values.Set("r", TimestampStr())
	values.Set("sid", wxrobot.secret.Sid)
	values.Set("uin", strconv.FormatInt(wxrobot.secret.Uin, 10))
	values.Set("skey", wxrobot.secret.Skey)
	values.Set("deviceid", wxrobot.secret.DeviceID)
	values.Set("synckey", wxrobot.secret.SyncKeyStr)
	values.Set("_", TimestampStr())

	b, err := wxrobot.httpClient.Get(uri, values)
	if err != nil {
		return nil, err
	}
	s := string(b)
	re := regexp.MustCompile(`window.synccheck=\{retcode:"(\d+)",selector:"(\d+)"\}`)
	matchs := re.FindStringSubmatch(s)
	if len(matchs) == 0 {
		log.Println(s)
		return nil, errors.New("find Sync check code error")
	}
	syncStatus := &syncStatus{Retcode: matchs[1], Selector: matchs[2]}
	return syncStatus, nil
}

func (wxrobot *WXRobot) Sync() ([]*Message, error) {
	values := &url.Values{}
	values.Set("sid", wxrobot.secret.Sid)
	values.Set("skey", wxrobot.secret.Skey)
	values.Set("lang", "en_US")
	values.Set("pass_ticket", wxrobot.secret.PassTicket)
	url := fmt.Sprintf("%s/webwxsync?%s", wxrobot.secret.BaseUri, values.Encode())
	b, err := wxrobot.httpClient.PostJson(url, map[string]interface{}{
		"BaseRequest": wxrobot.baseRequest,
		"SyncKey":     wxrobot.secret.SyncKey,
		"rr":          ^int(time.Now().Unix()) + 1,
	})
	if err != nil {
		return nil, err
	}

	var r SyncResponse
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	if r.BaseResponse.Ret != 0 {
		log.Println(string(b))
		// log.Printf("%+v", r.BaseResponse)
		return nil, errors.New("sync error")
	}
	wxrobot.updateSyncKey(r.SyncKey)
	return r.MsgList, nil
}

func (wxrobot *WXRobot) HandleMsgs(ms []*Message) {
	for _, m := range ms {
		wxrobot.HandleMsg(m)
	}
}

func (wxrobot *WXRobot) SendMsgToMyself(msg string) error {
	return wxrobot.SendMsg(wxrobot.user.UserName, msg)
}

func (wxrobot *WXRobot) SendMsg(userName, msg string) error {
	values := &url.Values{}
	values.Set("pass_ticket", wxrobot.secret.PassTicket)
	url := fmt.Sprintf("%s/webwxsendmsg?%s", wxrobot.secret.BaseUri, values.Encode())
	msgId := fmt.Sprintf("%d%s", Timestamp()*1000, RandNumbers(4))
	b, err := wxrobot.httpClient.PostJson(url, map[string]interface{}{
		"BaseRequest": wxrobot.baseRequest,
		"Msg": map[string]interface{}{
			"Type":         1,
			"Content":      msg,
			"FromUserName": wxrobot.user.UserName,
			"ToUserName":   userName,
			"LocalID":      msgId,
			"ClientMsgId":  msgId,
		},
		"Scene": 0,
	})
	if err != nil {
		return err
	}
	return wxrobot.CheckCode(b, "发送消息失败")
}
func (wxrobot *WXRobot) HandleMsg(m *Message) {
	log.Printf("[%s] from %s to %s : %s", MSG_TYPE_MAP[m.MsgType], wxrobot.GetUserName(m.FromUserName), wxrobot.GetUserName(m.ToUserName), m.Content)
	switch m.MsgType {
	case TEXT_MSG: // 文本消息
		if wxrobot.messageHandler.TextHandler != nil {
			wxrobot.messageHandler.TextHandler(m)
		}
	case IMG_MSG:// 图片消息
		if wxrobot.messageHandler.ImgHandler != nil {
			wxrobot.messageHandler.ImgHandler(m)
		}
	case VOICE_MSG:// 语音消息
		if wxrobot.messageHandler.VoiceHandler != nil {
			wxrobot.messageHandler.VoiceHandler(m)
		}
	case FACE_0_MSG:// 表情消息
		if wxrobot.messageHandler.FaceHandler != nil {
			wxrobot.messageHandler.FaceHandler(m)
		}
	case FACE_1_MSG:// 表情消息
		if wxrobot.messageHandler.FaceHandler != nil {
			wxrobot.messageHandler.FaceHandler(m)
		}
	case LINK_MSG:// 链接消息
		if wxrobot.messageHandler.LinkHandler != nil {
			wxrobot.messageHandler.LinkHandler(m)
		}
	case ENTER_CHAT_MSG:// 用户在手机进入某个联系人聊天界面时收到的消息
		if wxrobot.messageHandler.EnterChatHandler != nil {
			wxrobot.messageHandler.EnterChatHandler(m)
		}
	default:

		if wxrobot.messageHandler.UnKnowHandler != nil {
			wxrobot.messageHandler.UnKnowHandler(m)
		}
	}

}

func (wxrobot *WXRobot) Listening() error {
	err := wxrobot.TestSyncCheck()
	if err != nil {
		return err
	}
	for {
		syncStatus, err := wxrobot.SyncCheck()
		if err != nil {
			log.Printf("sync check error: %s", err.Error())
			time.Sleep(3 * time.Second)
			continue
		}
		if syncStatus.Retcode == "1100" {
			return errors.New("从微信客户端上登出")
		} else if syncStatus.Retcode == "1101" {
			return errors.New("从其它设备上登了网页微信")
		} else if syncStatus.Retcode == "0" {
			if syncStatus.Selector == "0" { // 无更新
				continue
			} else if syncStatus.Selector == "2" { // 有新消息
				ms, err := wxrobot.Sync()
				if err != nil {
					log.Printf("sync err: %s", err.Error())
				}
				wxrobot.HandleMsgs(ms)
			} else { // 可能有其他类型的消息，直接丢弃
				log.Printf("New Message, Unknow type: %+v", syncStatus)
				_, err := wxrobot.Sync()
				if err != nil {

				}
			}
		} else if syncStatus.Retcode == "1102" {
			return fmt.Errorf("Sync Error %+v", syncStatus)
		} else {
			log.Printf("sync check Unknow Code: %+v", syncStatus)
		}
	}
}

func (wxrobot *WXRobot) Start() error {
	newLoginUri, err := wxrobot.GetNewLoginUrl()
	if err != nil {
		return err
	}

	err = wxrobot.NewLoginPage(newLoginUri)
	if err != nil {
		return err
	}

	err = wxrobot.Init()
	if err != nil {
		return err
	}

	// err = wxRobot.StatusNotify()
	// if err != nil {
	// 	return err
	// }

	err = wxrobot.GetContacts()
	if err != nil {
		return err
	}
	return wxrobot.Listening()
}
