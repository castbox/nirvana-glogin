package plat

import (
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"glogin/pbs/glogin"
)

var (
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

// elkAlarm http
func elkAlarm(status string, url string, msg interface{}) {
	log.Warnw("elkAlarm http", "status", status, "url", url, "msg", msg)
}
