package account

import (
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"glogin/db"
	"glogin/db/db_core"
	"glogin/internal"
	"glogin/internal/anti"
	"glogin/internal/appsflyer"
	"glogin/internal/hawkeye"
	"glogin/internal/plat"
	"glogin/pbs/glogin"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"time"
)

func CheckNotExist(filter interface{}) bool {
	return db.CheckNotExist(filter, db.AccountTableName())
}

func CreateVisitor(request *glogin.VisitorLoginReq, visitor string, ip string) (internal.Rsp, error) {
	document := bson.M{"visitor": visitor, "create": bson.M{"time": time.Now().Unix(), "ip": ip, "bundle_id": request.Game.BundleId}}
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	return create(document, req)
}

func CreateThird(request *glogin.ThirdLoginReq, dbField string, unionId string, ip string) (internal.Rsp, error) {
	document := bson.M{dbField: unionId, "create": bson.M{"time": time.Now().Unix(), "ip": ip, "bundle_id": request.Game.BundleId}}
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	return create(document, req)
}

func CreatePhone(request *glogin.SmsLoginReq, ip string) (internal.Rsp, error) {
	document := bson.M{"phone": request.Phone, "create": bson.M{"time": time.Now().Unix(), "ip": ip, "bundle_id": request.Game.BundleId}}
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	return create(document, req)
}

func LoginPhone(request *glogin.SmsLoginReq, ip string) (internal.Rsp, error) {
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	loginRsp, err := login(bson.M{"phone": request.Phone}, req)
	if err != nil {
		return internal.Rsp{}, err
	}
	return loginRsp, nil
}

func LoginBundleThird(request *glogin.ThirdLoginReq, dbField string, uid interface{}, ip string) (internal.Rsp, error) {
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	loginRsp, err := login(bson.M{dbField: uid, "bundle_id": request.Game.BundleId}, req)
	if err != nil {
		return internal.Rsp{}, err
	}
	return loginRsp, nil
}

func LoginThird(request *glogin.ThirdLoginReq, dbField string, unionId string, ip string) (internal.Rsp, error) {
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	loginRsp, err := login(bson.M{dbField: unionId}, req)
	if err != nil {
		return internal.Rsp{}, err
	}
	return loginRsp, nil
}

func LoginFast(request *glogin.FastLoginReq, dhAccountId int32, ip string) (internal.Rsp, error) {
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	loginRsp, err := login(bson.M{"_id": dhAccountId}, req)
	if err != nil {
		return internal.Rsp{}, err
	}
	return loginRsp, nil
}

func LoginVisitor(request *glogin.VisitorLoginReq, visitor string, ip string) (internal.Rsp, error) {
	req := internal.Req{IP: ip, Client: request.Client, Game: request.Game}
	loginRsp, err := login(bson.M{"visitor": visitor}, req)
	if err != nil {
		return internal.Rsp{}, err
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
	// 自动激活gameCD  k-v  gameCd - time
	gameMap := make(map[string]db_core.GameInfo)
	gameMap[req.Game.GameCd] = db_core.GameInfo{
		Time: time.Now().Unix(),
	}
	accountInfo["games"] = gameMap
	dhid, errInsert := db.CreateDhId(accountInfo)
	if errInsert != nil {
		err = errInsert
		return
	} else {
		rsp.AccountData.ID = dhid
		req.Account = fmt.Sprintf("%d", dhid)
		req.GameCd = req.Game.GameCd
		//anti_addiction
		antiRsp, antiErr := anti.StateQuery(req)
		if antiErr != nil {
			err = antiErr
			return
		}
		rsp.AntiRsp = antiRsp
		rsp.GameRsp.FirstLogin = true
		// todo:: appsflyer
		appsflyer.SendAppsFlyer(req)
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
	req.Account = fmt.Sprintf("%d", accountData.ID)
	hawkRsp, hawkErr := hawkeye.CheckLogin(req)
	if hawkErr != nil {
		return internalRsp, hawkErr
	}
	// 防沉迷查询
	req.GameCd = req.Game.GameCd
	log.Infow("anti.StateQuery req", "req", req)
	antiRsp, antiErr := anti.StateQuery(req)
	if antiErr != nil {
		return internalRsp, antiErr
	}
	log.Infow("anti.StateQuery rsp ", "rsp", antiRsp, "req", req)
	internalRsp.AccountData = accountData
	internalRsp.HawkRsp = hawkRsp
	internalRsp.AntiRsp = antiRsp
	// 更新信息
	gameRsp, upErr := updateLoginInfo(req, accountData)
	internalRsp.GameRsp = gameRsp
	log.Infow("login internalRsp ", "rsp", internalRsp, "upErr", upErr)
	return internalRsp, nil
}

func updateLoginInfo(req internal.Req, accountData db_core.AccountData) (db_core.GameRsp, error) {
	upData := bson.M{"last_login": time.Now().Unix()}
	gameRsp := db_core.GameRsp{
		FirstLogin: false,
	}
	if accountData.Games == nil {
		accountData.Games = make(map[string]db_core.GameInfo)
		accountData.Games[req.GameCd] = db_core.GameInfo{
			Time: time.Now().Unix(),
		}
		upData["games"] = accountData.Games
		gameRsp.FirstLogin = true
	} else {
		_, ok := accountData.Games[req.GameCd]
		if !ok {
			accountData.Games[req.GameCd] = db_core.GameInfo{
				Time: time.Now().Unix(),
			}
			upData["games"] = accountData.Games
			gameRsp.FirstLogin = true
		}
	}

	setData := bson.M{"$set": upData}
	return gameRsp, db.UpdateOne(bson.M{"_id": accountData.ID}, setData, db.AccountTableName())
}

// BindThird 游客账号绑定第三方
func Load(filter interface{}) (db_core.AccountData, error) {
	return db.Load(filter)
}

func BindThird(accountId int32, thirdPlat, thirdUid string) error {
	update := bson.M{"$set": bson.M{thirdPlat: thirdUid}}
	return db.UpdateOne(bson.M{"_id": accountId}, update, db.AccountTableName())
}

func GetPlat(in interface{}) (platString string, err error) {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return "", fmt.Errorf("only accepts struct or struct pointer; got %T", v)
	}
	t := reflect.TypeOf(in)
	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			data[t.Field(i).Tag.Get("bson")] = v.Field(i).Interface()
		}
	}
	for key := range plat.ThirdList {
		if _, ok := data[key]; ok { // val 存储的是第三方uid，只有找到返回key 就可以了
			platString = key
			return
		}
	}
	return
}
