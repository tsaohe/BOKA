package BOKA

import (
	"time"
)

type Token struct {
	Token string
	ID    string
	Expir int64
}

var TOKEN = Token{}

func (T *Token) get() (string, string) {
	timeNow := time.Now().Unix()

	if T.Token == "" || T.ID == "" || timeNow-T.Expir > CONFIG.BoKa.Sec {
		T.RequestToken()
	}
	return T.Token, T.ID
}

func (T *Token) RequestToken() {
	client := HTTP{
		Url:    "https://api.bokao2o.com/auth/merchant/v2/user/login",
		Params: nil,
		Headers: map[string]string{
			"referer": "https://s3.boka.vc/",
		},
		Body: map[string]interface{}{
			"custId":   CONFIG.BoKa.CustId,
			"compId":   CONFIG.BoKa.CompId,
			"userName": CONFIG.BoKa.UserName,
			"passWord": CONFIG.BoKa.PassWord,
			"source":   CONFIG.BoKa.Source,
		},
	}
	err := client.post()
	if err != nil {
		return
	}
	token := client.JSON.Get("result.token").String()
	shopid := client.JSON.Get("result.shopId").String()
	TOKEN.Token, TOKEN.ID = token, shopid
	TOKEN.Expir = time.Now().Unix()
}
