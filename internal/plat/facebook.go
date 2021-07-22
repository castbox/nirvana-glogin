package plat

import (
	"fmt"
	"glogin/pbs/glogin"
	"io/ioutil"
	"net/http"
)

const (
	facebookAuthKey = "google_oauth_url"
)

var Facebook facebook

type facebook struct{}

// Auth 登录返回第三方账号id 和 错误信息
func (f facebook) Auth(request *glogin.ThirdLoginReq) (string, string, error) {
	baseUrl := authURL(request.Game.BundleId, facebookAuthKey)
	url := baseUrl + request.ThirdToken
	resp, err := http.Get(url)
	if err != nil {
		resErr := fmt.Errorf("failed communicating with server: %v", err)
		elkAlarm("error", url, resErr)
		return "", "", resErr
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		elkAlarm(resp.Status, url, "")
		return "", "", fmt.Errorf(resp.Status)
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		resErr := fmt.Errorf("failed reading from metadata server: %v", err)
		elkAlarm("error", url, resErr)
		return "", "", resErr
	}
	return "", "", fmt.Errorf(resp.Status)
}

func (f facebook) String() string {
	return "facebook"
}

func (f facebook) DbFieldName() string {
	return "facebook"
}
