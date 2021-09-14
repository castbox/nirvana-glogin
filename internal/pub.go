package internal

import (
	"glogin/db/db_core"
	"glogin/pbs/glogin"
)

type Rsp struct {
	AccountData db_core.AccountData `json:"acc_data"`
	GameRsp     db_core.GameRsp     `json:"game_rsp"`
	HawkRsp     interface{}         `json:"hawk_rsp"`
	AntiRsp     interface{}         `json:"anti_rsp"`
}

type Req struct {
	GameCd  string              `json:"game_cd"`
	Account string              `json:"account"`
	IP      string              `json:"ip"`
	Client  *glogin.LoginClient `json:"client"`
	Game    *glogin.LoginGame   `json:"game"`
}
