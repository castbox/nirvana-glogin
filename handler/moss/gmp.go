package moss

import (
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"glogin/constant"
	"glogin/db"
	"glogin/db/db_core"
	"glogin/pbs/glogin"
	"glogin/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
)

const (
	NotSpecified = "notspecified"
)

type Gmp struct {
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
