package plat

import (
	"git.dhgames.cn/svr_comm/gmoss/v2"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"glogin/pbs/glogin"
)

var (
	logger    = logrus.WithField("component", "plat")
	ThirdList = map[string]third{
		"google":   Google,
		"facebook": Facebook,
		"ios":      IOS,
		"yedun":    YeDun,
	}
)

type third interface {
	// Auth 登录返回第三方账号tokenId openId 错误信息
	Auth(request *glogin.ThirdLoginReq) (string, string, error)
	String() string
}

// authURL 返回auth url地址
func authURL(bundleId string, platKey string) string {
	data := gmoss.DynamicCfg("glogin", bundleId, nil)
	if len(data) == 0 {
		return ""
	}

	return gjson.GetBytes(data, platKey).String()
}

// elkAlarm 运维日志
func elkAlarm(status string, url string, msg interface{}) {
	logger.Errorf("#elkAlarm#http#%s#%s#%v", status, url, msg)
}
