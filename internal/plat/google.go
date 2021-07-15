package plat

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

const (
	googleAuthKey = "google_oauth_url"
)

var Google google

type google struct{}

// Auth 登录返回第三方账号id 和 错误信息
func (g google) Auth(bundleId string, token string) (string, error) {
	logger.Debugf("%s -> bundleId:%s, token:%s", g, bundleId, token)
	baseUrl := authURL(bundleId, googleAuthKey)
	url := baseUrl + token
	resp, err := http.Get(url)
	if err != nil {
		resErr := fmt.Errorf("failed communicating with server: %v", err)
		elkAlarm("error", url, resErr)
		return "", resErr
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		elkAlarm(resp.Status, url, "")
		return "", fmt.Errorf(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resErr := fmt.Errorf("failed reading from metadata server: %v", err)
		elkAlarm("error", url, resErr)
		return "", resErr
	}

	uid := gjson.GetBytes(body, "sub").String()
	return uid, nil
}

func (g google) String() string {
	return "google"
}
