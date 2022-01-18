package db

import (
	"fmt"
	log "gitlab.degames.cn/svr_comm/gcore/glog"
	"gitlab.degames.cn/svr_comm/gcore/gmongo"
	"glogin/config"
	"glogin/db/db_core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

const (
	VerifyCodeTable = "glogin_verifycode"
	SMSLlt          = 600
	SMSCountLimit   = 3
)

// 支持配置优先
func VerifyCodeTableName() string {
	verifyTableName := config.Field("code_table_name").String()
	if verifyTableName == "" {
		verifyTableName = VerifyCodeTable
	}
	return verifyTableName
}

//indexModel := mongo.IndexModel{
//Keys: bsonx.Doc{{"expire_date", bsonx.Int32(1)}}, // 设置TTL索引列"expire_date"
//Options:options.Index().SetExpireAfterSeconds(1*24*3600), // 设置过期时间1天，即，条目过期一天过自动删除
// 创建索引
func InitVerifyCode() {
	log.Infow("verifyCode mongodb init", "table", VerifyCodeTableName())
	indexFiles := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{"phone", bsonx.Int32(1)}},
		},
		{
			Keys: bsonx.Doc{{"phone", bsonx.Int32(1)}, {"send_time", bsonx.Int32(1)}},
		},
		{
			Keys:    bsonx.Doc{{"expire", bsonx.Int32(1)}},         // 设置TTL索引列"expire"
			Options: options.Index().SetExpireAfterSeconds(SMSLlt), // 设置过期时间1天，即，
		},
	}
	gmongo.CreateIndexes(config.MongoUrl(), config.MongoDb(), VerifyCodeTableName(), indexFiles)
}

func AddSmsVerify(phone string, verifyCode string) (interface{}, error) {
	document := bson.M{}
	document["phone"] = phone
	document["verify_code"] = verifyCode
	document["send_time"] = time.Now().Unix()
	//timeStr := time.Now().Format("2006-01-02 15:04:05.000")
	time := time.Now()
	document["expire"] = time
	_, errInsert := gmongo.InsertOne(config.MongoUrl(), config.MongoDb(), VerifyCodeTableName(), document)
	if errInsert != nil {
		return nil, errInsert
	}
	log.Infow("AddSmsVerify ok", "phone", phone, "verify_code", verifyCode, "time", time)
	return nil, nil
}

func CheckSmsVerifyCode(phone string, verifyCode string) (bool, error) {
	filter := bson.M{}
	filter["phone"] = phone
	filter["verify_code"] = verifyCode
	doc, errFind := gmongo.FindOne(config.MongoUrl(), config.MongoDb(), VerifyCodeTableName(), filter)
	if errFind != nil {
		log.Warnw("CheckSmsVerifyCode Load", "err", errFind)
		return false, errFind
	}

	result := &db_core.VerifyCodeData{}
	errDecode := doc.Decode(result)
	if errDecode != nil {
		return false, errDecode
	}
	if result.Phone == phone {
		return true, nil
	}
	errRsp := fmt.Errorf("CheckSmsVerifyCode fail phone:%v ,dbFind:%v", phone, result)
	return false, errRsp
}

// 60s 只允许发送1次
func CheckSmsInterval(phone string) (bool, error) {
	filter := bson.M{}
	filter["phone"] = phone
	filter["send_time"] = bson.M{
		"$gt": time.Now().Unix() - 60,
	}
	doc, errFind := gmongo.FindOne(config.MongoUrl(), config.MongoDb(), VerifyCodeTableName(), filter)
	// 数据库操作失败
	if errFind != nil {
		log.Warnw("CheckSmsVerifyCode Load", "err", errFind)
		return false, errFind
	}
	result := &db_core.VerifyCodeData{}
	errDecode := doc.Decode(result)
	log.Infow("CheckSmsInterval find AccountData", "result", result)
	if errDecode != nil {
		if errDecode == mongo.ErrNoDocuments {
			return true, nil
		} else {
			return false, errDecode
		}
	}
	if result.Phone == "" {
		return true, nil
	}
	errRsp := fmt.Errorf("CheckSmsInterval fail:%v ,sendtime:%v", result.Phone, result)
	return false, errRsp
}

// 10min 只允许发送3次
func CheckSmsVerifyCount(phone string) (bool, error) {
	filter := bson.M{}
	filter["phone"] = phone
	count, errCount := gmongo.CountDocuments(config.MongoUrl(), config.MongoDb(), VerifyCodeTableName(), filter)
	if errCount != nil {
		log.Warnw("CheckSmsVerifyCount", "err", errCount)
		return false, errCount
	}
	if count < SMSCountLimit {
		return true, nil
	}
	resultErr := fmt.Errorf("CheckSmsVerifyCount Limit CurCount: %v", count)
	return false, resultErr
}
