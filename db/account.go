package db

import (
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"git.dhgames.cn/svr_comm/gcore/gmongo"
	"glogin/config"
	"glogin/db/db_core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
)

const (
	AccountTableName = "glogin_account"
	MinAccount       = 100000000
	MaxAccount       = 999999999
	MaxTryTime       = 10
)

// 创建索引
func InitAccount() {
	log.Infow("account mongodb init", "table", AccountTableName)
	indexFiles := []mongo.IndexModel{
		{
			Keys: bson.D{{"google", int32(1)}, {"bundle_id", int32(1)}},
		},
		{
			Keys: bson.D{{"facebook", int32(1)}, {"bundle_id", int32(1)}},
		},
		{
			Keys: bson.D{{"ios", int32(1)}, {"bundle_id", int32(1)}},
		},
		{
			Keys: bson.D{{"google", int32(1)}},
		},
		{
			Keys: bson.D{{"facebook", int32(1)}},
		},
		{
			Keys: bson.D{{"ios", int32(1)}},
		},
		{
			Keys: bson.D{{"visitor", int32(1)}},
		},
		{
			Keys: bson.D{{"phone", int32(1)}},
		},
		{
			Keys: bson.D{{"we_chat", int32(1)}},
		},
		{
			Keys: bson.D{{"qq", int32(1)}},
		},
		{
			Keys: bson.D{{"huawei", int32(1)}},
		},
	}
	gmongo.CreateIndexes(config.MongoUrl(), config.MongoDb(), AccountTableName, indexFiles)
}

func CheckNotExist(filter interface{}) bool {
	count, err := gmongo.CountDocuments(config.MongoUrl(), config.MongoDb(), AccountTableName, filter)
	if err != nil {
		log.Warnw("CheckNotExist", "err", err)
		return false
	}
	if count == 0 {
		return true
	}
	return false
}

func Load(filter interface{}) (result db_core.AccountData, err error) {
	doc, errFind := gmongo.FindOne(config.MongoUrl(), config.MongoDb(), AccountTableName, filter)
	if errFind != nil {
		log.Warnw("AccountTable Load", "err", err)
		err = errFind
		return
	}
	err = doc.Decode(&result)
	return
}

func CreateDhId(document bson.M) (accountId int32, err error) {
	i := 0
	for ; i < MaxTryTime; i++ {
		accountId = rand.Int31n(MaxAccount-MinAccount) + MinAccount
		document["_id"] = accountId
		_, errInsert := gmongo.InsertOne(config.MongoUrl(), config.MongoDb(), AccountTableName, document)
		if errInsert == nil {
			log.Infow("new account ok", "account", accountId, "times", i)
			return
		}
	}
	return
}
