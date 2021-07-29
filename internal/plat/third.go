package plat

import (
	"github.com/sirupsen/logrus"
	"glogin/pbs/glogin"
)

var (
	logger    = logrus.WithField("component", "plat")
	ThirdList = map[string]third{
		"google":   Google,
		"facebook": Facebook,
		"ios":      IOS,
		"yedun":    YeDun,
		"qq":       QQ,
		"wechat":   WeChat,
	}
)

type third interface {
	// Auth 登录返回第三方账号tokenId openId 错误信息
	Auth(request *glogin.ThirdLoginReq) (string, string, error)
	String() string
	DbFieldName() string
}

// elkAlarm 运维日志
func elkAlarm(status string, url string, msg interface{}) {
	logger.Errorf("#elkAlarm#http#%s#%s#%v", status, url, msg)
}
