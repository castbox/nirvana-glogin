package plat

import (
	"glogin/pbs/glogin"
)

var WeChat wechat

type wechat struct{}

// Auth 登录返回第三方账号id 和 错误信息
func (w wechat) Auth(request *glogin.ThirdLoginReq) (string, string, error) {
	return "", "", nil
}

func (w wechat) String() string {
	return "wechat"
}

func (w wechat) DbFieldName() string {
	return "we_chat"
}
