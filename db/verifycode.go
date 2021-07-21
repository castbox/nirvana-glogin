package db

import (
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"git.dhgames.cn/svr_comm/gcore/gmongo"
	"glogin/config"
	"glogin/db/db_core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

const (
	VerifyCodeTableName = "glogin_verifycode"
	SMSLlt              = 600
	SMSCountLimit       = 3
)

//
//indexModel := mongo.IndexModel{
//Keys: bsonx.Doc{{"expire_date", bsonx.Int32(1)}}, // 设置TTL索引列"expire_date"
//Options:options.Index().SetExpireAfterSeconds(1*24*3600), // 设置过期时间1天，即，条目过期一天过自动删除
// 创建索引
func InitVerifyCode() {
	log.Infow("verifyCode mongodb init", "table", VerifyCodeTableName)
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
			//Keys:    bson.D{{"expire", int32(1)}},
			//Options: options.Index().SetExpireAfterSeconds(SMSLlt),
		},
	}
	gmongo.CreateIndexes(config.MongoUrl(), config.MongoDb(), VerifyCodeTableName, indexFiles)
}

//-spec add_sms_verify(Phone :: binary(), VerifyCode :: binary()) -> ok | {error, any()}.
//add_sms_verify(Phone, VerifyCode) ->
//{M, S, Ms} = os:timestamp(),
//Doc = #{phone => Phone, verify_code => VerifyCode, send_time => now(), expire => {M, S, Ms}},
//emgo:insert(?URL, ?NS, Doc).
func AddSmsVerify(phone string, verifyCode string) (interface{}, error) {
	document := bson.M{}
	document["phone"] = phone
	document["verify_code"] = verifyCode
	document["send_time"] = time.Now().Unix()
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	document["expire"] = timeStr
	//localTime:=document["expire"] = time.Now().Local()
	_, errInsert := gmongo.InsertOne(config.MongoUrl(), config.MongoDb(), VerifyCodeTableName, document)
	if errInsert != nil {
		return nil, errInsert
	}
	log.Infow("AddSmsVerify ok", "phone", phone, "verify_code", verifyCode, "timeStr", timeStr)
	return nil, nil
}

//-spec check_sms_verify_code(Phone :: binary(), VerifyCode :: binary()) -> ok | {error, integer()}.
//check_sms_verify_code(Phone, VerifyCode) ->
//case emgo:findOne(?URL, ?NS, #{phone => Phone, verify_code => VerifyCode}) of
//{ok, #{<<"phone">> := Phone}} -> ok;
//{ok, #{}} -> {error, ?ERR_ULOGIN_Login_VerifyCode};
//{error, _} -> {error, ?ERR_DB}
//end.
func CheckSmsVerifyCode(phone string, verifyCode string) (bool, error) {
	filter := bson.M{}
	filter["phone"] = phone
	filter["verify_code"] = verifyCode
	doc, errFind := gmongo.FindOne(config.MongoUrl(), config.MongoDb(), VerifyCodeTableName, filter)
	if errFind != nil {
		log.Warnw("CheckSmsVerifyCode Load", "err", errFind)
		return false, errFind
	}
	result := db_core.AccountData{}
	errDecode := doc.Decode(&result)
	if errDecode == nil {
		if result.Phone == phone {
			return true, nil
		}
	}
	return false, errFind
}

//%% 60s 只允许发送1次
//-spec check_sms_interval(Phone :: binary()) -> ok | {error, any()}.
//check_sms_interval(Phone) ->
//NowTime = util:now(),
//case emgo:find(?URL, ?NS, #{<<"phone">> => Phone, <<"send_time">> => #{<<"$gt">> => NowTime - 60}}) of
//{ok, Cursor} ->
//Docs = em_cursor:rest(Cursor),
//case Docs of
//[] -> ok;
//_ -> {error, ?ERR_ULOGIN_SMS_INTERVAL}
//end;
//Err ->
//lager:error("check_sms_interval err: ~p", [Err]),
//{error, ?ERR_DB}
//end.

func CheckSmsInterval(phone string) (bool, error) {
	filter := bson.M{}
	filter["phone"] = phone
	nowTime := time.Now().Unix()
	filter["send_time"] = bson.M{
		"$gt": nowTime - 60,
	}
	doc, errFind := gmongo.FindOne(config.MongoUrl(), config.MongoDb(), VerifyCodeTableName, filter)
	if errFind != nil {
		log.Warnw("CheckSmsVerifyCode Load", "err", errFind)
		return false, errFind
	}
	result := db_core.AccountData{}
	errDecode := doc.Decode(&result)
	if errDecode == nil {
		if result.Phone == "" {
			return true, nil
		}
	}
	return false, errFind
}

//%% 10min 只允许发送3次
//-spec check_sms_verify_count(Phone :: binary()) -> ok | {error, any()}.
//check_sms_verify_count(Phone) ->
//case emgo:count(?URL, ?NS, #{<<"phone">> => Phone}) of
//{ok, Cnt} ->
//if
//Cnt < ?SMS_COUNT_LIMIT -> ok;
//true -> {error, ?ERR_ULOGIN_SMS_COUNT}
//end;
//Err ->
//lager:error("check_sms_verify_count err: ~p", [Err]),
//{error, ?ERR_DB}
//end.
func CheckSmsVerifyCount(phone string) (bool, error) {
	filter := bson.M{}
	filter["phone"] = phone
	count, errCount := gmongo.CountDocuments(config.MongoUrl(), config.MongoDb(), VerifyCodeTableName, filter)
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
