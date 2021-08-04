package plat

import (
	"fmt"
	"github.com/tidwall/gjson"
	"glogin/pbs/glogin"
	"io/ioutil"
	"net/http"
)

const (
	QQAuthUrl   = "qq_oauth_url"
	UnionIDLast = "&unionid=1"
)

var QQ qq

type qq struct{}

// Auth 登录返回第三方账号tokenId openId 错误信息
// https://graph.qq.com/oauth2.0/me?access_token=ACCESSTOKEN&unionid=1
func (q qq) Auth(request *glogin.ThirdLoginReq) (*AuthRsp, error) {
	baseUrl := authURL(request.Game.BundleId, QQAuthUrl)
	url := baseUrl + request.ThirdToken + UnionIDLast
	resp, err := http.Get(url)
	if err != nil {
		resErr := fmt.Errorf("failed communicating with server: %v", err)
		elkAlarm("error", url, resErr)
		return nil, resErr
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		elkAlarm(resp.Status, url, "")
		return nil, fmt.Errorf(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resErr := fmt.Errorf("failed reading from metadata server: %v", err)
		elkAlarm("error", url, resErr)
		return nil, resErr
	}

	openID := gjson.GetBytes(body, "openid").String()
	unionID := gjson.GetBytes(body, "unionid").String()
	return &AuthRsp{
		Uid:     openID,
		UnionId: unionID,
	}, nil

}

func (q qq) String() string {
	return "qq"
}

func (q qq) DbFieldName() string {
	return "qq"
}
