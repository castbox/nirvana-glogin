package cgi

import (
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"glogin/constant"
	"glogin/internal/account"
	"glogin/internal/plat"
	"glogin/internal/smfpcrypto"
	"glogin/internal/sms"
	glogin "glogin/pbs/glogin"
	"glogin/util"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

var (
	PlatIsWrong = fmt.Errorf("third plat is wrong")
)

type Login struct {
}

// http://127.0.0.1:8080/login/sms verify 获取验证码登陆 login 短信验证码登陆
func (Login) SMS(request *glogin.SmsLoginReq) (response *glogin.SmsLoginRsp, err error) {
	response = &glogin.SmsLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
	}
	// sms step verify 并获得Phone 验证码
	// sms step login  手机验证码登陆
	if request.Step == "verify" {
		log.Infow("SMS verify", "request", request)
		_, err := sms.SmsVerify(request)
		if err != nil {
			response.Code = constant.ErrCodeSMSGetVerifyFaild
			response.Errmsg = fmt.Sprintf("get verify error %s request", request.Phone)
			return response, nil
		}
		return response, nil
	} else if request.Step == "login" {
		log.Infow("SMS login", "request", request)
		bCheck, _ := sms.CheckPhone(request.Phone, request.Verifycode)
		if bCheck {

		}
	}
	return response, nil
}

func (Login) Third(request *glogin.ThirdLoginReq) (response *glogin.ThridLoginRsp, err error) {
	response = &glogin.ThridLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
	}
	uid, openId, errAuth := ThirdAuth(request)
	log.Infow("Third", "uid", uid, "openid", openId, "err", err)
	// 平台错误
	if errAuth == PlatIsWrong {
		response.Code = constant.ErrCodePlatIsWrong
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
	smID := smfpcrypto.ParseSMID(request.Client.Dhid)
	request.GetClient().Dhid = smID
	if account.CheckNotExist(bson.M{request.ThirdPlat: openId}) {
		// 账号不存在,创建
		accountId, errCreate := account.CreateThird(request, uid, openId)
		if errCreate != nil {
			response.Code = 500
			response.Errmsg = fmt.Sprintf("thrid plat %s auth error: %s", request.ThirdPlat, errAuth)
			return response, errCreate
		}
		response.Code = 0
		response.DhToken = util.GenDHToken(accountId)
	} else {
		// 账号存在, 直接登录
		loginRsp, errLogin := account.LoginThird(request, uid, openId)
		if errLogin != nil {
			response.Code = 500
			response.Errmsg = fmt.Sprintf("thrid plat %s auth error: %s", request.ThirdPlat, errAuth)
			return response, nil
		}
		// 返回
		rsp, ok := loginRsp.(account.InternalRsp)
		if !ok {
			response.Code = 500
			response.Errmsg = fmt.Sprintf("thrid plat %s auth error: %s", request.ThirdPlat, errLogin)
			return response, nil
		}
		response.Code = 0
		response.DhAccount = rsp.AccountData.ID
		response.SmId = smID
		response.DhToken = util.GenDHToken(rsp.AccountData.ID)
		response.Errmsg = fmt.Sprintf("thrid plat %s auth error: %s", request.ThirdPlat, errAuth)
	}
	return response, nil
}

func (Login) Visitor(req *glogin.VisitorLoginReq) (rsp *glogin.VisitorLoginRsp, err error) {
	rsp = &glogin.VisitorLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
	}
	// 参数验证
	if req.Dhid == "" {
		rsp.Code = constant.ErrCodeSMIDError
		rsp.Errmsg = fmt.Sprintf("visitor is agrs  error: %s", req.Dhid)
		return rsp, nil
	}
	// 数美ID解析
	smId := smfpcrypto.ParseSMID(req.Client.Dhid)
	// 创建唯一设备账号，如果有SMID就用SMID创建，如果没有就随机
	visitorId := CreateVisitorID(req.Dhid, smId)
	if visitorId == "" {
		rsp.Code = constant.ErrCodeSMIDError
		rsp.Errmsg = fmt.Sprintf("visitor is agrs smid error: %s", req.Dhid)
		return rsp, nil
	}

	if account.CheckNotExist(bson.M{"visitor": visitorId}) {
		// 账号不存在,创建
		accountId, errCreate := account.CreateVisitor(req, visitorId)
		if errCreate != nil {
			log.Infow("visitor fast login create account error ", "visitor", visitorId)
			rsp.Code = 500
			rsp.Errmsg = fmt.Sprintf("visitor  %s fast login create account error: %s", visitorId, errCreate)
			return rsp, errCreate
		}
		rsp.Code = 0
		rsp.DhToken = util.GenDHToken(accountId)
		rsp.SmId = smId
		rsp.DhAccount = accountId
		rsp.Errmsg = "success"
		return rsp, nil
	} else {
		// 账号存在, 直接登录
		loginRsp, errLogin := account.LoginVisitor(req, visitorId)
		if errLogin != nil {
			log.Infow("visitor fast login 1 error", "visitor", visitorId)
			rsp.Code = 500
			rsp.Errmsg = fmt.Sprintf("visitor %s fast login error: %s", req.Dhid, errLogin)
			return rsp, nil
		}
		// 返回
		value, ok := loginRsp.(account.InternalRsp)
		if !ok {
			log.Infow("visitor fast login  2 error", "visitor", visitorId)
			rsp.Code = 500
			rsp.Errmsg = fmt.Sprintf("visitor %s fast login error: %s", req.Dhid, errLogin)
			return rsp, nil
		}
		log.Infow("visitor login success", "visitor", visitorId)
		rsp.Code = 0
		rsp.DhAccount = value.AccountData.ID
		rsp.SmId = smId
		rsp.DhToken = util.GenDHToken(value.AccountData.ID)
		rsp.Visitor = visitorId
		rsp.Errmsg = "success"
	}
	return rsp, nil
}

func (l Login) Fast2(request *glogin.FastLoginReq, ctx *gin.Context) (response *glogin.FastLoginRsp, err error) {
	return l.Fast(request)
}

func (Login) Fast(request *glogin.FastLoginReq) (response *glogin.FastLoginRsp, err error) {
	response = &glogin.FastLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
	}
	if request.DhToken == "" {
		response.Code = constant.ErrCodeTokenError
		response.Errmsg = fmt.Sprintf("fast is token null game: %s client: %s", request.Game, request.Client)
		return response, nil
	}

	// 校验token
	dhAccount, errToken := util.ValidDHToken(request.DhToken)
	if errToken != nil {
		if errToken == util.ExpiredToken {
			// 1:快速登录失败，token 过期 30天
			response.Code = constant.ErrCodeFastTokenExpired
			response.Errmsg = fmt.Sprintf("fast login: %s", errToken)
			return response, nil
		} else {
			response.Code = constant.ErrCodeFastTokenVaild
			response.Errmsg = fmt.Sprintf("fast login: %s", errToken)
			return response, nil
		}
	}
	// 数美ID解析
	smId := smfpcrypto.ParseSMID(request.Client.Dhid)
	loginRsp, errLogin := account.LoginFast(request, dhAccount)
	if errLogin != nil {
		response.Code = 500
		response.Errmsg = fmt.Sprintf("fast login %s auth error: %s", request.DhToken, errLogin)
		return response, nil
	}

	// 返回
	value, ok := loginRsp.(account.InternalRsp)
	if !ok {
		response.Code = 500
		response.Errmsg = fmt.Sprintf("fast login %s auth error: %s", request.DhToken, errLogin)
		return response, nil
	}
	response.Code = 0
	response.SmId = smId
	response.DhAccount = value.AccountData.ID
	response.DhToken = util.GenDHToken(value.AccountData.ID)
	response.Errmsg = "success"

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

func CreateVisitorID(srcVisitor string, dhId string) string {
	if srcVisitor == "" {
		return util.GetRandomString(32)
	}
	// 去掉时间戳的DHID
	lastPos := strings.LastIndex(srcVisitor, "-")
	if lastPos == -1 {
		return dhId
	}
	// 拼接时间戳
	timeParam := srcVisitor[lastPos:]
	return dhId + timeParam
}
