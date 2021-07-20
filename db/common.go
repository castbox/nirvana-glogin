package db

import (
	"git.dhgames.cn/svr_comm/gcore/gmongo"
	"glogin/config"
)

func InitMongo() {
	gmongo.Init(config.Field("mongo_url").String())
	InitAccount()
	InitVerifyCode()
}
