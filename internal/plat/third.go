package plat

import (
	log "gitlab.degames.cn/svr_comm/gcore/glog"
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

type AuthRsp struct {
	Uid     interface{} `json:"uid"`     // uid
	UnionId string      `json:"unionid"` // unionid
	Nick    string      `json:"nick"`    // 普通用户昵称
	Sex     int         `json:"sex"`     // 普通用户性别，1为男性，2为女性
	Country string      `json:"country"` // 国家，如中国为CN
}

type third interface {
	// Auth 登录返回第三方账号tokenId openId 错误信息
	Auth(request *glogin.ThirdLoginReq) (*AuthRsp, error)
	String() string
	DbFieldName() string
}

// elkAlarm http
func elkAlarm(status string, url string, msg interface{}) {
	log.Warnw("elkAlarm http", "status", status, "url", url, "msg", msg)
}
