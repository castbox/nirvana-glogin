package moss

import (
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"glogin/constant"
	"glogin/db"
	"glogin/pbs/glogin"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	NotSpecified = "notspecified"
)

type Gmp struct {
}

func (Gmp) LoadAccountInfo(request *glogin.QueryRequest) (response *glogin.QueryResponse, err error) {
	response = &glogin.QueryResponse{
		Code: constant.ErrCodeOk,
		Data: []*glogin.AccountData{},
	}
	// 未指定账号 需要根据条件进行组合查询
	if request.Account == NotSpecified {
		queryCond := getQueryConditions(request)
		//var accounts []glogin.AccountData
		accounts, count, errLook := db.Lookup(queryCond)
		if errLook != err {
			return response, nil
		}
		r2, ok := accounts.([]glogin.AccountData)
		if ok {
			for _, accountInfo := range r2 {
				response.Data = append(response.Data, &accountInfo)
			}
		}
		response.Code = 200
		response.Count = count
		response.Msg = "success"
		return response, nil
	} else {

	}
	return response, nil
}

// 获得账号信息 这里需要指定查询条件
func getQueryConditions(request *glogin.QueryRequest) bson.M {
	filter := bson.M{}
	log.Infow(" sofa rpc getQueryConditions req", "request", request)
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
	return filter
}
