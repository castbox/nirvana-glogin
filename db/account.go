package db

import (
	"context"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"git.dhgames.cn/svr_comm/gcore/gmongo"
	"glogin/config"
	"glogin/db/db_core"
	"glogin/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	AccountTable = "glogin_account"
	MinAccount   = 100000000
	MaxAccount   = 999999999
	MaxTryTime   = 10
)

var TableName = ""

// 支持配置优先
func AccountTableName() string {
	if TableName != "" {
		return TableName
	}
	TableName = config.Field("account_table_name").String()
	if TableName == "" {
		TableName = AccountTable
	}
	return TableName
}

// 创建索引
func InitAccount() {
	log.Infow("account mongodb init", "table", AccountTableName())
	indexFiles := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{"google", bsonx.Int32(1)}, {"bundle_id", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"facebook", bsonx.Int32(1)}, {"bundle_id", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"ios", bsonx.Int32(1)}, {"bundle_id", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"google", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"facebook", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"ios", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"visitor", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"phone", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"we_chat", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"qq", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"huawei", bsonx.Int32(1)}},
		},
	}
	gmongo.CreateIndexes(config.MongoUrl(), config.MongoDb(), AccountTableName(), indexFiles)
}

func Load(filter interface{}) (result db_core.AccountData, err error) {
	doc, errFind := gmongo.FindOne(config.MongoUrl(), config.MongoDb(), AccountTableName(), filter)
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
		//accountId = rand.Int31n(MaxAccount-MinAccount) + MinAccount
		accountId = util.Rand32Num(MinAccount, MaxAccount)
		log.Infow("rand a 32 account ", "account num", accountId)
		document["_id"] = accountId
		_, errInsert := gmongo.InsertOne(config.MongoUrl(), config.MongoDb(), AccountTableName(), document)
		if errInsert == nil {
			log.Infow("new account ok", "account", accountId, "times", i)
			return
		}
	}
	return
}

// gmp用
func Lookup(filter bson.M, option bson.M, ptrToSlice interface{}) (count int32, err error) {
	var pageSize = option["page_size"].(int32)
	var pageNum = option["page_num"].(int32)
	if pageSize == 0 || pageNum == 0 {
		pageSize = 50
		pageNum = 1
	}
	findOption := options.Find().SetLimit(int64(pageSize)).SetSkip(int64((pageNum - 1) * pageSize))
	//filter2 := bson.M{}
	doc, errFind := gmongo.Find(config.MongoUrl(), config.MongoDb(), AccountTableName(), filter, findOption)
	if errFind != nil {
		log.Warnw("AccountTable Lookup gmongo.Find", "errFind", errFind)
		err = errFind
		return
	}
	if err = doc.All(context.TODO(), ptrToSlice); err != nil {
		panic(err)
	}
	count2, errCount := gmongo.CountDocuments(config.MongoUrl(), config.MongoDb(), AccountTableName(), filter)
	if errCount != nil {
		log.Warnw("AccountTable Lookup  gmongo.CountDocuments", "errCount", errCount)
		err = errCount
		return
	}
	count = int32(count2)
	return
}
