package moss

import (
	"encoding/json"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"glogin/constant"
	"glogin/db"
	"glogin/db/db_core"
	"glogin/pbs/glogin"
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
		count, errLook := db.Lookup(queryCond, &accounts)
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
		errLoad := db.LoadOne(filter, &doc, db.AccountTableName)
		log.Infow(" sofa rpc LoadAccountInfo 1 response", "response", response)
		if errLoad != nil {
			if errLoad == mongo.ErrNoDocuments {
				response.Code = 200
				response.Count = 0
				response.Msg = "success"
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
		return response, nil
		log.Infow(" sofa rpc LoadAccountInfo 2 response", "response", response)
	}
	return response, nil
}

func dbConvertToPb(dbAccount db_core.AccountData) (pbAccount glogin.AccountData, err error) {
	jsonBlob, errB := json.Marshal(dbAccount)
	if err != nil {
		err = errB
		return
	}
	errUb := json.Unmarshal(jsonBlob, &pbAccount)
	if errUb != nil {
		err = errUb
	}
	//pbAccount.Id = dbAccount.ID
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
		filter[request.LoginType] = bson.M{
			"$exists": 1,
		}
	}
	//options
	filter["page_size"] = 50
	filter["page_num"] = 1
	return filter
}
