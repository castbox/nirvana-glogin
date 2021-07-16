package account

import (
	"github.com/gin-gonic/gin"
	"glogin/db"
	"glogin/internal/anti"
	"glogin/internal/hawkeye"
	"glogin/pbs/glogin"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func CheckNotExist(filter interface{}) bool {
	return db.CheckNotExist(filter)
}

func CreateVisitor(request *glogin.ThirdLoginReq, visitor string) (DhAccount int32, err error) {
	document := bson.M{visitor: visitor, "create": bson.M{"time": time.Now().Unix()}}
	return create(document, request)
}

func CreateThird(request *glogin.ThirdLoginReq, uid string, openId string) (DhAccount int32, err error) {
	document := bson.M{request.ThirdPlat: openId, "create": bson.M{"time": time.Now().Unix(), "third_uid": uid}}
	return create(document, request)
}

func LoginThird(request *glogin.ThirdLoginReq, uid string, openId string) (interface{}, error) {
	loginRsp, err := login(bson.M{request.ThirdPlat: openId}, request)
	if err != nil {
		return nil, err
	}
	return loginRsp, nil
}

func create(accountInfo bson.M, request *glogin.ThirdLoginReq) (DhAccount int32, err error) {
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

func login(filter interface{}, request *glogin.ThirdLoginReq) (interface{}, error) {
	// 加载数据
	accountData, err := db.Load(filter)
	if err != nil {
		return nil, err
	}
	// 鹰眼检查
	hawkRsp, hawkErr := hawkeye.CheckLogin()
	if hawkErr != nil {
		return nil, hawkErr
	}
	// 防沉迷检查
	antiRsp, antiErr := anti.Check(gin.H{})
	if antiErr != nil {
		return nil, antiErr
	}
	// 返回数据
	return gin.H{
		"acc_data": accountData,
		"hawk_rsp": hawkRsp,
		"anti_rsp": antiRsp,
	}, nil
}
