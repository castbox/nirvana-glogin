package account

import (
	"github.com/gin-gonic/gin"
	"glogin/db"
	"glogin/db/db_core"
	"glogin/internal/anti"
	"glogin/internal/hawkeye"
	"glogin/internal/plat"
	"glogin/pbs/glogin"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func CheckNotExist(filter interface{}) bool {
	return db.CheckNotExist(filter)
}

func CreateVisitor(request *glogin.VisitorLoginReq, visitor string) (DhAccount int32, err error) {
	document := bson.M{"visitor": visitor, "create": bson.M{"time": time.Now().Unix(), "bundle_id": request.Game.BundleId}}
	return create(document, request)
}

func CreateThird(request *glogin.ThirdLoginReq, uid string, openId string) (DhAccount int32, err error) {
	document := bson.M{request.ThirdPlat: openId, "create": bson.M{"time": time.Now().Unix(), "third_uid": uid, "bundle_id": request.Game.BundleId}}
	return create(document, request)
}

func LoginThird(request *glogin.ThirdLoginReq, uid string, openId string) (interface{}, error) {
	loginRsp, err := login(bson.M{request.ThirdPlat: openId}, request)
	if err != nil {
		return nil, err
	}
	return loginRsp, nil
}

func LoginFast(request *glogin.FastLoginReq, dhAccountId int32) (interface{}, error) {
	loginRsp, err := login(bson.M{"_id": dhAccountId}, request)
	if err != nil {
		return nil, err
	}
	return loginRsp, nil
}

func LoginVisitor(request *glogin.VisitorLoginReq, visitor string) (interface{}, error) {
	loginRsp, err := login(bson.M{"visitor": visitor}, request)
	if err != nil {
		return nil, err
	}
	return loginRsp, nil
}

func create(accountInfo bson.M, request interface{}) (DhAccount int32, err error) {
	// todo 鹰眼check注册
	id, errInsert := db.CreateDhId(accountInfo)
	if errInsert != nil {
		err = errInsert
		return
	} else {
		DhAccount = id
		// todo appsflyer
		// todo anti_addiction
		// todo 防沉迷数据查询
		return
	}
}

type InternalRsp struct {
	AccountData db_core.AccountData `json:"acc_data"`
	HawkRsp     interface{}         `json:"hawk_rsp"`
	AntiRsp     interface{}         `json:"anti_rsp"`
}

func login(filter interface{}, request interface{}) (InternalRsp, error) {
	internalRsp := InternalRsp{}
	// 加载数据
	accountData, err := db.Load(filter)
	if err != nil {
		return internalRsp, err
	}
	// 鹰眼检查
	hawkRsp, hawkErr := hawkeye.CheckLogin()
	if hawkErr != nil {
		return internalRsp, hawkErr
	}
	// 防沉迷检查
	antiRsp, antiErr := anti.Check(gin.H{})
	if antiErr != nil {
		return internalRsp, antiErr
	}
	internalRsp.AccountData = accountData
	internalRsp.HawkRsp = hawkRsp
	internalRsp.AntiRsp = antiRsp
	// 返回数据
	return internalRsp, nil
}

func GetPlat(result bson.M) (platString string) {
	for key := range plat.ThirdList {
		if _, ok := result[key]; ok { // val 存储的是第三方uid，只有找到返回key 就可以了
			platString = key
			break
		}
	}
	return
}
