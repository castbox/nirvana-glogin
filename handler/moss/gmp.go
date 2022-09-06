package moss

import (
	"fmt"
	log "gitlab.degames.cn/svr_comm/gcore/glog"
	"github.com/pborman/uuid"
	"glogin/constant"
	"glogin/db"
	"glogin/db/db_core"
	"glogin/internal/account"
	"glogin/internal/plat"
	"glogin/pbs/glogin"
	"glogin/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"strconv"
	"time"
)

const (
	NotSpecified = "notspecified"
)

type Gmp struct {
}

func (Gmp) ChangeBind(request *glogin.ChangeBindReq) (response *glogin.ChangeBindRsp, err error) {
	reqId := uuid.New()
	log.Infow("sofa rpc ChangeBind req", "reqId", reqId, "request", request)
	response = &glogin.ChangeBindRsp{
		Code: constant.ErrCodeOk,
		Data: &glogin.AccountData{},
	}
	before := time.Now().UnixNano()
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("got panic", "err", err)
		}
		log.Infow("ofa rpc ChangeBind rsp", "reqId", reqId, "response", response, "time_cost", (time.Now().UnixNano()-before)/1000000)
	}()
	if request.Account == "" || request.Phone == "" {
		response.Code = constant.ErrCodeThirdBindFail
		return response, nil
	}
	// 1.判断手机号是否已经使用
	if !account.CheckNotExist(bson.M{"phone": request.Phone}) {
		response.Code = constant.ErrCodePhoneAlreadyBind
		return response, nil
	}
	// 2.判断账号是否存在
	// 指定账号 直接查出 #{<<"_id">> => AccountID}
	dhAccount, err := strconv.Atoi(request.Account)
	if err != nil {
		return response, err
	}
	doc := db_core.AccountData{}
	errLoad := db.LoadOne(bson.M{"_id": dhAccount}, &doc, db.AccountTableName())
	if errLoad != nil {
		response.Code = constant.ErrCodeNoAccount
		return response, errLoad
	}
	// 3.执行换绑操作；原第三方token屏蔽加上前缀
	thirdPlat := request.Plat
	dbField := ""
	if thirdPlat == "visitor" {
		dbField = thirdPlat
	}
	if third, ok := plat.ThirdList[thirdPlat]; ok {
		dbField = third.DbFieldName()
	}
	// 原先的值加前缀
	//"{gmp}#oOYr9t5gv7RsJNoTlc5icYy_b-NQ"
	prefix := "{gmp}#"
	thirdToken, getErr := GetToken(doc, dbField)
	if getErr != nil || thirdToken == "" {
		response.Code = constant.ErrCodeNoThirdErr
		return response, errLoad
	}
	upData := bson.M{
		dbField: prefix + thirdToken,
		"phone": request.Phone,
	}
	setData := bson.M{"$set": upData}
	upDataErr := db.UpdateOne(bson.M{"_id": dhAccount}, setData, db.AccountTableName())
	if upDataErr != nil {
		response.Code = constant.ErrCodeUpdateFail
		return response, upDataErr
	}
	// 4.返回最新
	doc = db_core.AccountData{}
	_ = db.LoadOne(bson.M{"_id": dhAccount}, &doc, db.AccountTableName())
	response.Msg = "success"
	pbAcc, errC := dbConvertToPb(doc)
	if errC == nil {
		response.Data = &pbAcc
	}
	return response, nil
}

func GetToken(in interface{}, dbField string) (platToken string, err error) {
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
	if _, ok := data[dbField]; ok { // val 存储的是第三方uid，只有找到返回key 就可以了
		platToken = data[dbField].(string)
		return
	}
	return
}

func (Gmp) QueryAccount(request *glogin.QueryReq) (response *glogin.QueryRsp, err error) {
	log.Infow(" sofa rpc QueryAccount req", "request", request)
	response = &glogin.QueryRsp{
		Code: constant.ErrCodeOk,
		Data: []*glogin.AccountData{},
	}
	// 指定账号
	if len(request.Accounts) == 1 {
		// 指定账号 直接查出 #{<<"_id">> => AccountID}
		dhAccount, err := strconv.Atoi(request.Accounts[0])
		if err != nil {
			return response, err
		}
		filter := bson.M{
			"_id": dhAccount,
		}
		doc := db_core.AccountData{}
		errLoad := db.LoadOne(filter, &doc, db.AccountTableName())
		if errLoad != nil {
			if errLoad == mongo.ErrNoDocuments {
				response.Code = 200
				response.Count = 0
				response.Msg = "success"
				log.Infow(" sofa rpc LoadAccountInfo 1 response", "response", response)
				return response, nil
			} else {
				return response, err
			}
		}
		response.Code = 200
		response.Count = 1
		response.Msg = "success"
		pbAcc, errC := dbConvertToPb(doc)
		if errC == nil {
			response.Data = append(response.Data, &pbAcc)
		}
		log.Infow(" sofa rpc LoadAccountInfo 2 response", "response", response)
		return response, nil
	} else {
		queryCond := getQueryConditions2(request)
		var accounts []db_core.AccountData
		//option
		option := bson.M{}
		option["page_size"] = request.PageSize
		option["page_num"] = request.PageNum
		count, errLook := db.Lookup(queryCond, option, &accounts)
		if errLook != err {
			return response, nil
		}
		for _, accountInfo := range accounts {
			pbAcc, errC := dbConvertToPb(accountInfo)
			if errC == nil {
				response.Data = append(response.Data, &pbAcc)
			}
		}
		response.Code = 200
		response.Count = count
		response.Msg = "success"
		log.Infow(" sofa rpc LoadAccountInfo response", "response", response)
		return response, nil
	}
	return response, nil
}

func (Gmp) LoadAccountInfo(request *glogin.QueryRequest) (response *glogin.QueryResponse, err error) {
	log.Infow(" sofa rpc LoadAccountInfo req", "request", request)
	response = &glogin.QueryResponse{
		Code: constant.ErrCodeOk,
		Data: []*glogin.AccountData{},
	}
	// 未指定账号 需要根据条件进行组合查询
	if request.Account == NotSpecified {
		queryCond := getQueryConditions(request)
		accounts := []db_core.AccountData{}
		//option
		option := bson.M{}
		option["page_size"] = request.PageSize
		option["page_num"] = request.PageNum
		count, errLook := db.Lookup(queryCond, option, &accounts)
		if errLook != err {
			return response, nil
		}
		for _, accountInfo := range accounts {
			pbAcc, errC := dbConvertToPb(accountInfo)
			if errC == nil {
				response.Data = append(response.Data, &pbAcc)
			}
		}
		response.Code = 200
		response.Count = count
		response.Msg = "success"
		log.Infow(" sofa rpc LoadAccountInfo response", "response", response)
		return response, nil
	} else {
		// 指定账号 直接查出 #{<<"_id">> => AccountID}
		dhAccount, err := strconv.Atoi(request.Account)
		if err != nil {
			return response, err
		}
		filter := bson.M{
			"_id": dhAccount,
		}
		log.Infow(" LoadAccountInfo test 1 ", dhAccount)
		doc := db_core.AccountData{}
		errLoad := db.LoadOne(filter, &doc, db.AccountTableName())
		if errLoad != nil {
			if errLoad == mongo.ErrNoDocuments {
				response.Code = 200
				response.Count = 0
				response.Msg = "success"
				log.Infow(" sofa rpc LoadAccountInfo 1 response", "response", response)
				return response, nil
			} else {
				return response, err
			}
		}
		log.Infow(" LoadAccountInfo test 2 ", doc.Token)
		response.Code = 200
		response.Count = 1
		response.Msg = "success"
		pbAcc, errC := dbConvertToPb(doc)
		log.Infow(" LoadAccountInfo test 3 ", pbAcc.Token)
		if errC == nil {
			response.Data = append(response.Data, &pbAcc)
		}
		log.Infow(" sofa rpc LoadAccountInfo 2 response", "response", response)
		return response, nil
	}
	return response, nil
}

func dbConvertToPb(dbAccount db_core.AccountData) (pbAccount glogin.AccountData, err error) {
	pbAccount.XId = dbAccount.ID
	pbAccount.BundleId = dbAccount.BundleID
	pbAccount.Facebook = dbAccount.Facebook
	pbAccount.Ios = util.BsonAToStr(dbAccount.IOS)
	pbAccount.Google = dbAccount.Google
	pbAccount.Phone = dbAccount.Phone
	pbAccount.Visitor = dbAccount.Visitor
	pbAccount.LastLogin = dbAccount.LastLogin
	pbAccount.Token = dbAccount.Token
	strIp := util.BsonAToStr(dbAccount.Create.Ip)
	pbAccount.Create = &glogin.CreateData{
		Ip:       strIp,
		BundleId: dbAccount.Create.BundleId,
		SmId:     dbAccount.Create.SmId,
		Time:     dbAccount.Create.Time,
	}
	return
}

// 获得账号信息 这里需要指定查询条件
func getQueryConditions2(request *glogin.QueryReq) bson.M {
	filter := bson.M{}
	if len(request.Bundleids) > 0 {
		var bundleFilters []bson.M
		for _, bundleId := range request.Bundleids {
			bundleFilters = append(bundleFilters, bson.M{
				"bundle_id": bundleId,
			})
		}
		filter["$or"] = bundleFilters
	}
	if len(request.Accounts) > 0 {
		var accountFilters []bson.M
		for _, accountID := range request.Accounts {
			dhAccount, err := strconv.Atoi(accountID)
			if err != nil {
				continue
			}
			accountFilters = append(accountFilters, bson.M{
				"_id": dhAccount,
			})
		}
		filter["$or"] = accountFilters
	}
	if request.LoginType != "" {
		filter[request.LoginType] = bson.M{
			"$exists": 1,
		}
	}
	return filter
}

// 获得账号信息 这里需要指定查询条件
func getQueryConditions(request *glogin.QueryRequest) bson.M {
	filter := bson.M{}
	if len(request.Bundleids) > 0 {
		var bundleFilters []bson.M
		for bundleId := range request.Bundleids {
			bundleFilters = append(bundleFilters, bson.M{
				"bundle_id": bundleId,
			})
		}
		filter["$or"] = bundleFilters
	} else {
		if request.LoginType != "" {
			filter[request.LoginType] = bson.M{
				"$exists": 1,
			}
		}
	}
	return filter
}
