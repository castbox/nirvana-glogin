package cgi

import (
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"glogin/constant"
	"glogin/internal/account"
	"glogin/internal/smfpcrypto"
	glogin2 "glogin/pbs/glogin"
	"glogin/util"
	"go.mongodb.org/mongo-driver/bson"
)

type Bind struct {
}

//游客绑定第三方
func (Bind) BindThird(request *glogin2.VistorBindThridReq) (response *glogin2.VistorBindThridRsp, err error) {
	log.Infow("BindThird", "request", request)
	response = &glogin2.VistorBindThridRsp{
		Code: constant.ErrCodeOk,
		ExtendData: &glogin2.ExtendData{
			Authentication: &glogin2.StateQueryResponse{},
		},
	}
	// 请求参数错误
	if request.Visitor == "" || request.ThirdPlat == "" {
		response.Code = constant.ErrCodeBindType
		return response, nil
	}
	if account.CheckNotExist(bson.M{"visitor": request.Visitor}) {
		// 游客账号不存在
		log.Infow("BindThird visitor account not exist error ", "visitor", request.Visitor)
		response.Code = constant.ErrCodeBindVisitorNotExist
		return response, nil
	}
	// 游客账号不正确
	dhAccount, err := account.Load(bson.M{"visitor": request.Visitor})
	if err != nil || dhAccount.ID == 0 {
		response.Code = constant.ErrCodeVisitorLoadErr
		response.Errmsg = fmt.Sprintf("BindThird plat %s auth error: %s", request.ThirdPlat, err)
		return response, nil
	}

	data, err := bson.Marshal(dhAccount)
	if err != nil {
		log.Infow("BindThird thirdUid bson.Marshal error  ", "err", err)
		return response, nil
	}
	result := bson.M{}
	err3 := bson.Unmarshal(data, result)
	if err3 != nil {
		log.Infow("BindThird thirdUid bson.Unmarshal error  ", "err", err)
		return response, nil
	}
	platString := account.GetPlat(result)
	// 账号已经绑定
	if platString != "" {
		log.Infow("BindThird Visitor had bind error  ", "platString", platString, "Visitor", request.Visitor)
		return response, nil
	}
	// 第三方账号验证检查
	authReq := &glogin2.ThirdLoginReq{
		ThirdPlat:  request.ThirdPlat,
		ThirdToken: request.ThirdToken,
	}
	authRsp, dbField, errAuth := ThirdAuth(authReq)
	if errAuth == PlatIsWrong {
		response.Code = constant.ErrCodePlatWrong
		response.Errmsg = fmt.Sprintf("BindThird %s plat is wrong", request.ThirdPlat)
		return response, nil
	}
	// 第三方账号验证失败
	if errAuth != nil || authRsp == nil {
		response.Code = constant.ErrCodeThirdAuthFail
		response.Errmsg = fmt.Sprintf("BindThird plat %s auth error: %s", request.ThirdPlat, errAuth)
		return response, nil
	}
	uid := authRsp.Uid
	unionId := authRsp.UnionId
	log.Infow("BindThird ThirdAuth Rsp", "uid", uid, "unionId", unionId)
	// 第三方账号已经绑定
	if !account.CheckNotExist(bson.M{dbField: unionId}) {
		log.Infow("BindThird thirdUid third_already_bind ", "thirdUid", unionId, "thirdplat", request.ThirdPlat)
		response.Code = constant.ErrCodeThirdAlreadyBind
		return response, nil
	}
	// 绑定第三方账号
	errBind := account.BindThird(dhAccount.ID, request.ThirdPlat, unionId)
	if errBind != nil {
		response.Code = constant.ErrCodeThirdBindFail
		return response, nil
	}
	// 数美ID解析
	smID := smfpcrypto.ParseSMID(request.Dhid)
	response.SmId = smID
	response.ThirdPlat = request.ThirdPlat
	response.DhAccount = dhAccount.ID
	response.DhToken = util.GenDHToken(dhAccount.ID)
	response.Errmsg = "success"
	return response, nil
}
