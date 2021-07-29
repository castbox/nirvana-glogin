package db

import (
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"git.dhgames.cn/svr_comm/gcore/gmongo"
	"glogin/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	TokenForBusinessTable = "fb_token_for_business"
)

// 创建索引
func InitForBusiness() {
	log.Infow("ForBusiness mongodb init", "table", TokenForBusinessTable)
	indexFiles := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{"facebook_token", bsonx.Int32(1)}, {"token_for_business", bsonx.Int32(1)}},
		},
	}
	gmongo.CreateIndexes(config.MongoUrl(), config.MongoDb(), TokenForBusinessTable, indexFiles)
}

func AddFbTokenForBusiness(facebookId string, forBusiness string, bundleId string) (interface{}, error) {
	document := bson.M{}
	document["facebook_token"] = facebookId
	document["token_for_business"] = forBusiness
	document["bundle_id"] = bundleId
	_, errInsert := gmongo.InsertOne(config.MongoUrl(), config.MongoDb(), TokenForBusinessTable, document)
	if errInsert != nil {
		return nil, errInsert
	}
	log.Infow("AddFbTokenForBusiness ok", "document", document)
	return nil, nil
}
