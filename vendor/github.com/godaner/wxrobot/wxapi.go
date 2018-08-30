package wxrobot

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/546669204/golang-http-do"
	"github.com/mdp/qrterminal"
	"github.com/tuotoo/qrcode"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	LOGIN_URL = "https://login.weixin.qq.com"

)

type WXApi struct {
	httpClient  *Client
	secret      *wxSecret
	baseRequest *BaseRequest
	user        *User
	contacts    map[string]*User
	handler     *Handler
}

func NewWXApi() *WXApi {
	return &WXApi{
		httpClient:  NewClient(),
		secret:      &wxSecret{},
		baseRequest: &BaseRequest{},
		user:        &User{},
		handler:     &Handler{},
	}
}

func (wxApi *WXApi) GetUser(userName string) (*User, error) {
	u, ok := wxApi.contacts[userName]
	if ok {
		return u, nil
	} else {
		return nil, errors.New("Error User Not Exist")
	}
}

func (wxApi *WXApi) GetUserName(userName string) string {
	u, err := wxApi.GetUser(userName)
	if err != nil {
		return "[myself]"
	}
	if u.RemarkName != "" {
		return u.RemarkName
	} else {
		return u.NickName
	}
}

func (wxApi *WXApi) getUuid() (string, error) {
	values := &url.Values{}
	values.Set("appid", "wx782c26e4c19acffb")
	values.Set("fun", "new")
	values.Set("lang", "zh_CN")
	values.Set("_", TimestampStr())
	uri := fmt.Sprintf("%s/jslogin", LOGIN_URL)
	b, err := wxApi.httpClient.Get(uri, values)
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

func (wxApi *WXApi) ShowQRcodeUrl(uuid string) error {
	//qr url
	qrStr := fmt.Sprintf("%s/qrcode/%s", LOGIN_URL, uuid)

	log.Println("Please open link in browser: "+qrStr, " , or you can scan this qr :")
	//qr img
	qrStrP := &qrStr
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

	wxApi.handler.ShowQRHandler(qrStrP)

	return nil
}

func (wxApi *WXApi) WaitingForLoginConfirm(uuid string) (string, error) {
	re := regexp.MustCompile(`window.code=([0-9]*);`)
	tip := "1"
	for {
		values := &url.Values{}
		values.Set("uuid", uuid)
		values.Set("tip", tip)
		values.Set("_", TimestampStr())
		b, err := wxApi.httpClient.Get("https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login", values)
		if err != nil {
			log.Printf("HTTP GET err: %s", err.Error())
			return "", err
		}
		s := string(b)
		codes := re.FindStringSubmatch(s)
		if len(codes) == 0 {
			log.Printf("find window.code failed, origin response: %s\n", s)
			return "", errors.New("Unknow Error")
		} else {
			code := codes[1]
			if code == "408" {
				log.Println("login timeout, reconnecting...")
			} else if code == "400" {
				log.Println("login timeout, need refresh qrcode")
				return "", errors.New("need refresh qr")
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
				return "", errors.New("Unknow Error")
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

func (wxApi *WXApi) NewLoginPage(newLoginUri string) error {
	b, err := wxApi.httpClient.Get(newLoginUri+"&fun=new", nil)
	if err != nil {
		log.Printf("HTTP GET err: %s", err.Error())
		return err
	}
	err = xml.Unmarshal(b, wxApi.secret)
	if err != nil {
		log.Printf("parse wxSecret from xml failed: %v", err)
		return err
	}
	if wxApi.secret.Code == "0" {
		u, _ := url.Parse(newLoginUri)
		wxApi.secret.BaseUri = newLoginUri[:strings.LastIndex(newLoginUri, "/")]
		wxApi.secret.Host = u.Host
		wxApi.secret.DeviceID = "e" + RandNumbers(15)
		return nil
	} else {
		return errors.New("Get wxSecret Error")
	}

}

func (wxApi *WXApi) Init() error {
	values := &url.Values{}
	values.Set("r", TimestampStr())
	values.Set("lang", "en_US")
	values.Set("pass_ticket", wxApi.secret.PassTicket)
	url := fmt.Sprintf("%s/webwxinit?%s", wxApi.secret.BaseUri, values.Encode())
	wxApi.baseRequest = &BaseRequest{
		Uin:      wxApi.secret.Uin,
		Sid:      wxApi.secret.Sid,
		Skey:     wxApi.secret.Skey,
		DeviceID: wxApi.secret.DeviceID,
	}
	b, err := wxApi.httpClient.PostJson(url, map[string]interface{}{
		"BaseRequest": wxApi.baseRequest,
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
		wxApi.user = r.User
		wxApi.updateSyncKey(r.SyncKey)
		return nil
	}
	return fmt.Errorf("Init error: %+v", r.BaseResponse)
}

func (wxApi *WXApi) updateSyncKey(s *SyncKey) {
	wxApi.secret.SyncKey = s
	syncKeys := make([]string, s.Count)
	for n, k := range s.List {
		syncKeys[n] = fmt.Sprintf("%d_%d", k.Key, k.Val)
	}
	wxApi.secret.SyncKeyStr = strings.Join(syncKeys, "|")
}

func (wxApi *WXApi) GetNewLoginUrl() (string, error) {
	uuid, err := wxApi.getUuid()
	if err != nil {
		return "", err
	}
	err = wxApi.ShowQRcodeUrl(uuid)
	if err != nil {
		return "", err
	}
	newLoginUri, err := wxApi.WaitingForLoginConfirm(uuid)
	if err != nil {
		return "", err
	}
	return newLoginUri, nil
}


func (wxApi *WXApi) StatusNotify() error {
	values := &url.Values{}
	values.Set("lang", "zh_CN")
	values.Set("pass_ticket", wxApi.secret.PassTicket)
	url := fmt.Sprintf("%s/webwxstatusnotify?%s", wxApi.secret.BaseUri, values.Encode())
	b, err := wxApi.httpClient.PostJson(url, map[string]interface{}{
		"BaseRequest":  wxApi.baseRequest,
		"code":         3,
		"FromUserName": wxApi.user.UserName,
		"ToUserName":   wxApi.user.UserName,
		"ClientMsgId":  TimestampMicroSecond(),
	})
	if err != nil {
		return err
	}
	return wxApi.CheckCode(b, "Status Notify error")
}

func (wxApi *WXApi) CheckCode(b []byte, errmsg string) error {
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

func (wxApi *WXApi) GetContacts() error {
	values := &url.Values{}
	values.Set("seq", "0")
	values.Set("pass_ticket", wxApi.secret.PassTicket)
	values.Set("skey", wxApi.secret.Skey)
	values.Set("r", TimestampStr())
	url := fmt.Sprintf("%s/webwxgetcontact?%s", wxApi.secret.BaseUri, values.Encode())
	b, err := wxApi.httpClient.PostJson(url, map[string]interface{}{})
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
	wxApi.contacts = make(map[string]*User, r.MemberCount)
	return wxApi.updateContacts(r.MemberList)
}

func (wxApi *WXApi) updateContacts(us []*User) error {
	for _, u := range us {
		wxApi.contacts[u.UserName] = u
		log.Printf("%s => %s", u.UserName, u.NickName)
	}
	return nil
}

func (wxApi *WXApi) TestSyncCheck() error {
	for _, h := range []string{"webpush.", "webpush2."} {
		wxApi.secret.PushHost = h + wxApi.secret.Host
		syncStatus, err := wxApi.SyncCheck()
		if err == nil {
			if syncStatus.Retcode == "0" {
				return nil
			}
		}
	}
	return errors.New("Test SyncCheck error")
}

func (wxApi *WXApi) SyncCheck() (*SyncStatus, error) {
	uri := fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin/synccheck", wxApi.secret.PushHost)
	values := &url.Values{}
	values.Set("r", TimestampStr())
	values.Set("sid", wxApi.secret.Sid)
	values.Set("uin", strconv.FormatInt(wxApi.secret.Uin, 10))
	values.Set("skey", wxApi.secret.Skey)
	values.Set("deviceid", wxApi.secret.DeviceID)
	values.Set("synckey", wxApi.secret.SyncKeyStr)
	values.Set("_", TimestampStr())

	b, err := wxApi.httpClient.Get(uri, values)
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
	syncStatus := &SyncStatus{Retcode: matchs[1], Selector: matchs[2]}
	return syncStatus, nil
}

func (wxApi *WXApi) Sync() ([]*Message, error) {
	values := &url.Values{}
	values.Set("sid", wxApi.secret.Sid)
	values.Set("skey", wxApi.secret.Skey)
	values.Set("lang", "en_US")
	values.Set("pass_ticket", wxApi.secret.PassTicket)
	url := fmt.Sprintf("%s/webwxsync?%s", wxApi.secret.BaseUri, values.Encode())
	b, err := wxApi.httpClient.PostJson(url, map[string]interface{}{
		"BaseRequest": wxApi.baseRequest,
		"SyncKey":     wxApi.secret.SyncKey,
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
	wxApi.updateSyncKey(r.SyncKey)
	return r.MsgList, nil
}

func (wxApi *WXApi) HandleMsgs(ms []*Message) {
	for _, m := range ms {
		wxApi.HandleMsg(m)
	}
}

func (wxApi *WXApi) SendMsgToMyself(msg string) error {
	return wxApi.SendMsg(wxApi.user.UserName, msg)
}

func (wxApi *WXApi) SendMsg(userName, msg string) error {
	values := &url.Values{}
	values.Set("pass_ticket", wxApi.secret.PassTicket)
	url := fmt.Sprintf("%s/webwxsendmsg?%s", wxApi.secret.BaseUri, values.Encode())
	msgId := fmt.Sprintf("%d%s", Timestamp()*1000, RandNumbers(4))
	b, err := wxApi.httpClient.PostJson(url, map[string]interface{}{
		"BaseRequest": wxApi.baseRequest,
		"Msg": map[string]interface{}{
			"Type":         1,
			"Content":      msg,
			"FromUserName": wxApi.user.UserName,
			"ToUserName":   userName,
			"LocalID":      msgId,
			"ClientMsgId":  msgId,
		},
		"Scene": 0,
	})
	if err != nil {
		return err
	}
	return wxApi.CheckCode(b, "发送消息失败")
}

var (
	MSG_TEXT       int64 = 1
	MSG_IMG        int64 = 3
	MSG_VOICE      int64 = 34
	MSG_FACE_0     int64 = 43
	MSG_FACE_1     int64 = 47
	MSG_LINK       int64 = 49
	MSG_ENTER_CHAT int64 = 51

	MSG_TYPE_MAP = map[int64]string{
		MSG_TEXT:       "MSG_TEXT",
		MSG_IMG:        "MSG_IMG",
		MSG_VOICE:      "MSG_VOICE",
		MSG_FACE_0:     "MSG_FACE_0",
		MSG_FACE_1:     "MSG_FACE_1",
		MSG_LINK:       "MSG_LINK",
		MSG_ENTER_CHAT: "MSG_ENTER_CHAT",
	}
)

func (wxApi *WXApi) HandleMsg(m *Message) {
	log.Printf("[%s] from %s to %s : %s", MSG_TYPE_MAP[m.MsgType], wxApi.GetUserName(m.FromUserName), wxApi.GetUserName(m.ToUserName), m.Content)
	switch m.MsgType {
	case MSG_TEXT: // 文本消息
		if wxApi.handler.TextHandler != nil {
			wxApi.handler.TextHandler(m)
		}
	case MSG_IMG: // 图片消息
		if wxApi.handler.ImgHandler != nil {
			wxApi.handler.ImgHandler(m)
		}
	case MSG_VOICE: // 语音消息
		if wxApi.handler.VoiceHandler != nil {
			wxApi.handler.VoiceHandler(m)
		}
	case MSG_FACE_0: // 表情消息
		if wxApi.handler.FaceHandler != nil {
			wxApi.handler.FaceHandler(m)
		}
	case MSG_FACE_1: // 表情消息
		if wxApi.handler.FaceHandler != nil {
			wxApi.handler.FaceHandler(m)
		}
	case MSG_LINK: // 链接消息
		if wxApi.handler.LinkHandler != nil {
			wxApi.handler.LinkHandler(m)
		}
	case MSG_ENTER_CHAT: // 用户在手机进入某个联系人聊天界面时收到的消息
		if wxApi.handler.EnterChatHandler != nil {
			wxApi.handler.EnterChatHandler(m)
		}
	default:

		if wxApi.handler.UnKnowHandler != nil {
			wxApi.handler.UnKnowHandler(m)
		}
	}

}

const (
	SYSNC_STATUS_RETCODE_LOGOUT_FROM_WX_CLIENT = "1100"
	SYSNC_STATUS_RETCODE_LOGIN_WEB             = "1101"
	SYSNC_STATUS_RETCODE_NORMAL                = "0"
	SYSNC_STATUS_RETCODE_ERROR                 = "1102"
	SYSNC_STATUS_SELECTOR_NO_UPDATE            = "0"
	SYSNC_STATUS_SELECTOR_HAVE_UPDATE          = "2"
)

func (wxApi *WXApi) Listening() error {
	err := wxApi.TestSyncCheck()
	if err != nil {
		return err
	}
	for {
		syncStatus, err := wxApi.SyncCheck()
		if err != nil {
			log.Printf("sync check error: %s", err.Error())
			time.Sleep(3 * time.Second)
			continue
		}
		switch syncStatus.Retcode {
		case SYSNC_STATUS_RETCODE_LOGOUT_FROM_WX_CLIENT:
			return errors.New("从微信客户端上登出")
		case SYSNC_STATUS_RETCODE_LOGIN_WEB:
			return errors.New("从其它设备上登了网页微信")
		case SYSNC_STATUS_RETCODE_NORMAL:
			wxApi.handleSysncRetCodeNormal(syncStatus)
		case SYSNC_STATUS_RETCODE_ERROR:
			return fmt.Errorf("Sync Error %+v", syncStatus)
		default:
			log.Printf("sync check Unknow Code: %+v", syncStatus)

		}

	}
}
func (wxApi *WXApi)handleSysncRetCodeNormal(syncStatus *SyncStatus){
	switch syncStatus.Selector {
	case SYSNC_STATUS_SELECTOR_NO_UPDATE:
		break
	case SYSNC_STATUS_SELECTOR_HAVE_UPDATE:
		ms, err := wxApi.Sync()
		if err != nil {
			log.Printf("sync err: %s", err.Error())
		}
		wxApi.HandleMsgs(ms)
	default:
		log.Printf("New Message, Unknow type: %+v", syncStatus)
		_, err := wxApi.Sync()
		if err != nil {

		}
	}
}

func (wxApi *WXApi) Start() error {
	newLoginUri, err := wxApi.GetNewLoginUrl()
	if err != nil {
		return err
	}

	err = wxApi.NewLoginPage(newLoginUri)
	if err != nil {
		return err
	}

	err = wxApi.Init()
	if err != nil {
		return err
	}

	err = wxApi.StatusNotify()
	if err != nil {
		return err
	}

	err = wxApi.GetContacts()
	if err != nil {
		return err
	}
	return wxApi.Listening()
}
