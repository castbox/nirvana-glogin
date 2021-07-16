package plat

import (
	_ "github.com/gogf/gf/encoding/gjson"
	"glogin/pbs/glogin"
)

const (
	facebookAuthKey2 = "google_oauth_url"
)

var YeDun yedun

type yedun struct{}

// Auth 登录返回第三方账号id 和 错误信息
func (y yedun) Auth(request *glogin.ThirdLoginReq) (string, string, error) {
	return request.ThirdToken, request.ThirdToken, nil
}

func (y yedun) String() string {
	return "yedun"
}
