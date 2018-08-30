package wxrobot

var wxRobot *WXRobot

func Init(messageHandler *MessageHandler) error {
	wxRobot = NewWeixin(messageHandler)
	newLoginUri, err := wxRobot.GetNewLoginUrl()
	if err != nil {
		return err
	}

	err = wxRobot.NewLoginPage(newLoginUri)
	if err != nil {
		return err
	}

	err = wxRobot.Init()
	if err != nil {
		return err
	}

	err = wxRobot.GetContacts()
	if err != nil {
		return err
	}
	return nil
}

func Listening() error {
	return wxRobot.Listening()
}

func GetContacts() (map[string]*User, error) {
	return wxRobot.contacts, nil
}

func SendMsg(userId, msg string) error {
	return wxRobot.SendMsg(userId, msg)
}

