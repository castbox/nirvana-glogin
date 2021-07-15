package plat

import (
	"fmt"
	_ "github.com/gogf/gf/encoding/gjson"
	"io/ioutil"
	"net/http"
)

const (
	facebookAuthKey2 = "google_oauth_url"
)

var YeDun yedun

type yedun struct{}

// Auth 登录返回第三方账号id 和 错误信息
func (y yedun) Auth(bundleId string, token string) (string, error) {
	logger.Debugf("%s -> bundleId:%s, token:%s", y, bundleId, token)
	baseUrl := authURL(bundleId, facebookAuthKey2)
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

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		resErr := fmt.Errorf("failed reading from metadata server: %v", err)
		elkAlarm("error", url, resErr)
		return "", resErr
	}

	//data := gjson.GetBytes(body, "data").Map()
	//if len(data) == 0 {
	//	resErr := fmt.Errorf("failed reading data from metadata server: %v", err)
	//	elkAlarm("error", url, resErr)
	//	return "", resErr
	//}
	//
	//isValid := data["is_valid"]
	//uid := data["user_id"]
	//if isValid.Exists() == false || uid.Exists() == false {
	//	resErr := fmt.Errorf("failed reading isValid from metadata server: %v", err)
	//	elkAlarm("error", url, resErr)
	//	return "", resErr
	//}
	//
	//if isValid.Bool() == false {
	//	return "", fmt.Errorf("Invalid Auth Access Token.")
	//}

	return "", nil
}

func (y yedun) String() string {
	return "yedun"
}
