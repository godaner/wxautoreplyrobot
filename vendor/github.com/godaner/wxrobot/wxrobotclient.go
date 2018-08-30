package wxrobot

var wxApi *WXApi
func init(){
	wxApi = NewWXApi()
}
func SetClientHandler(messageHandler *Handler){
	wxApi.handler = messageHandler
}
func StartClient() error {

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

	err = wxApi.GetContacts()
	if err != nil {
		return err
	}
	err = wxApi.Listening()
	if err != nil {
		return err
	}
	return nil
}


func SendMsg(userName,msg string) error{
	return wxApi.SendMsg(userName,msg)
}