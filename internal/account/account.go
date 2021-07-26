package account

import (
	"glogin/db"
	"glogin/internal"
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

func CreateVisitor(request *glogin.VisitorLoginReq, visitor string, ip string) (interface{}, error) {
	document := bson.M{"visitor": visitor, "create": bson.M{"time": time.Now().Unix(), "ip": ip, "bundle_id": request.Game.BundleId}}
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	return create(document, req)
}

func CreateThird(request *glogin.ThirdLoginReq, dbField string, openId string, ip string) (interface{}, error) {
	document := bson.M{dbField: openId, "create": bson.M{"time": time.Now().Unix(), "ip": ip, "bundle_id": request.Game.BundleId}}
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	return create(document, req)
}

func CreatePhone(request *glogin.SmsLoginReq, ip string) (interface{}, error) {
	document := bson.M{"phone": request.Phone, "create": bson.M{"time": time.Now().Unix(), "ip": ip, "bundle_id": request.Game.BundleId}}
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	return create(document, req)
}

func LoginPhone(request *glogin.SmsLoginReq, ip string) (interface{}, error) {
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	loginRsp, err := login(bson.M{"phone": request.Phone}, req)
	if err != nil {
		return nil, err
	}
	return loginRsp, nil
}

func LoginThird(request *glogin.ThirdLoginReq, dbField string, uid string, openId string, ip string) (interface{}, error) {
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	loginRsp, err := login(bson.M{dbField: openId}, req)
	if err != nil {
		return nil, err
	}
	return loginRsp, nil
}

func LoginFast(request *glogin.FastLoginReq, dhAccountId int32, ip string) (interface{}, error) {
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	loginRsp, err := login(bson.M{"_id": dhAccountId}, req)
	if err != nil {
		return nil, err
	}
	return loginRsp, nil
}

func LoginVisitor(request *glogin.VisitorLoginReq, visitor string, ip string) (interface{}, error) {
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	loginRsp, err := login(bson.M{"visitor": visitor}, req)
	if err != nil {
		return nil, err
	}
	return loginRsp, nil
}

func create(accountInfo bson.M, req internal.Req) (rsp internal.Rsp, err error) {
	// 鹰眼check注册
	_, hErr := hawkeye.CheckRegister(req)
	if hErr != nil {
		err = hErr
		return
	}
	dhid, errInsert := db.CreateDhId(accountInfo)
	if errInsert != nil {
		err = errInsert
		return
	} else {
		rsp.AccountData.ID = dhid
		req.Account = string(dhid)
		req.GameCd = req.Game.GameCd
		// todo anti_addiction
		antiRsp, antiErr := anti.StateQuery(req)
		if antiErr != nil {
			err = antiErr
			return
		}
		rsp.AntiRsp = antiRsp
		// todo appsflyer
		return
	}
}

func login(filter interface{}, req internal.Req) (internal.Rsp, error) {
	internalRsp := internal.Rsp{}
	// 加载数据
	accountData, err := db.Load(filter)
	if err != nil {
		return internalRsp, err
	}
	// 鹰眼检查
	req.Account = string(accountData.ID)
	hawkRsp, hawkErr := hawkeye.CheckLogin(req)
	if hawkErr != nil {
		return internalRsp, hawkErr
	}
	// 防沉迷查询
	req.GameCd = req.Game.GameCd
	antiRsp, antiErr := anti.StateQuery(req)
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
