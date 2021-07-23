package internal

import (
	"glogin/db/db_core"
	"glogin/pbs/glogin"
)

type Rsp struct {
	AccountData db_core.AccountData `json:"acc_data"`
	HawkRsp     interface{}         `json:"hawk_rsp"`
	AntiRsp     interface{}         `json:"anti_rsp"`
}

type Req struct {
	Account string              `json:"account"`
	IP      string              `json:"ip"`
	Client  *glogin.LoginClient `json:"client"`
	Game    *glogin.LoginGame   `json:"game"`
}
