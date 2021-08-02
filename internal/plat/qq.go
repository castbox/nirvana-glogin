package plat

import (
	"glogin/pbs/glogin"
)

var QQ qq

type qq struct{}

// Auth 登录返回第三方账号tokenId openId 错误信息
func (q qq) Auth(request *glogin.ThirdLoginReq) (string, string, error) {
	return "", "", nil
}

func (q qq) String() string {
	return "qq"
}

func (q qq) DbFieldName() string {
	return "qq"
}
