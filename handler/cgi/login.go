package cgi

import (
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	_ "github.com/gin-gonic/gin"
	"glogin/constant"
	"glogin/internal/account"
	"glogin/internal/plat"
	glogin "glogin/pbs/glogin"
	"glogin/utils"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	PlatIsWrong = fmt.Errorf("third plat is wrong")
)

type Login struct {
}

func (Login) SMS(request *glogin.SmsLoginReq) (response *glogin.SmsLoginRsp, err error) {
	response = &glogin.SmsLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
	}
	return response, nil
}

func (Login) Third(request *glogin.ThirdLoginReq) (response *glogin.ThridLoginRsp, err error) {
	response = &glogin.ThridLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
	}
	uid, openId, errAuth := ThirdAuth(request)
	log.Fatalw("Third", "uid", "openid", "err", uid, openId, err)
	// 平台参数错误
	if errAuth == PlatIsWrong {
		response.Code = 500
		response.Errmsg = fmt.Sprintf("third %s plat is wrong", request.ThirdPlat)
		return response, nil
	}
	// 第三方账号验证失败
	if errAuth != nil {
		response.Code = 500
		response.Errmsg = fmt.Sprintf("thrid plat %s auth error: %s", request.ThirdPlat, errAuth)
		return response, nil
	}
	// 数美ID解析
	if account.CheckNotExist(bson.M{request.ThirdPlat: openId}) {
		// 账号不存在,创建
		accountId, errCreate := account.CreateThird(request, uid, openId)
		if errCreate != nil {
			response.Code = 500
			response.Errmsg = fmt.Sprintf("thrid plat %s auth error: %s", request.ThirdPlat, errAuth)
			return response, errCreate
		}
		response.Code = 0
		response.DhToken = utils.GenDHToken(accountId)
	} else {
		// 账号存在, 直接登录
		_, errLogin := account.LoginThird(request, uid, openId)
		if errLogin != nil {
			response.Code = 500
			response.Errmsg = fmt.Sprintf("thrid plat %s auth error: %s", request.ThirdPlat, errAuth)
			return response, nil
		}
		//loginRsp.xxx
		response.Code = 0
		//response.DhToken = utils.GenDHToken(loginRsp)
		response.Errmsg = fmt.Sprintf("thrid plat %s auth error: %s", request.ThirdPlat, errAuth)
	}
	return response, nil
}

func (Login) Visitor(request *glogin.VisitorLoginReq) (response *glogin.VisitorLoginRsp, err error) {
	response = &glogin.VisitorLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
	}
	return response, nil
}

func ThirdAuth(request *glogin.ThirdLoginReq) (uid string, openid string, err error) {
	thirdPlat := request.ThirdPlat
	if third, ok := plat.ThirdList[thirdPlat]; ok {
		uid, openid, err = third.Auth(request)
	} else {
		err = PlatIsWrong
	}
	return
}
