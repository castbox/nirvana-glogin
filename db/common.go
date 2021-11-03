package db

import (
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"git.dhgames.cn/svr_comm/gcore/gmongo"
	"glogin/config"
)

func InitMongo() {
	gmongo.Init(config.MongoUrl())
	InitAccount()
	InitVerifyCode()
	InitForBusiness()
}

func LoadOne(filter interface{}, result interface{}, tableName string) (err error) {
	doc, errFind := gmongo.FindOne(config.MongoUrl(), config.MongoDb(), tableName, filter)
	if errFind != nil {
		log.Warnw("AccountTable LoadOne", "err", err)
		err = errFind
		return
	}
	if errDecode := doc.Decode(result); errDecode != nil {
		err = errDecode
	}
	return
}

func CheckNotExist(filter interface{}, tableName string) bool {
	count, err := gmongo.CountDocuments(config.MongoUrl(), config.MongoDb(), tableName, filter)
	if err != nil {
		log.Warnw("CheckNotExist", "err", err)
		return false
	}
	if count == 0 {
		return true
	}
	return false
}

func UpdateOne(filter interface{}, update interface{}, tableName string) (err error) {
	_, errUpdate := gmongo.UpdateOne(config.MongoUrl(), config.MongoDb(), tableName, filter, update)
	if errUpdate != nil {
		log.Warnw("UpdateOne Table error", "tableName", tableName, "errUpdate", errUpdate)
		err = errUpdate
		return
	}
	return
}
