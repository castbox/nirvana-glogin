package cgi

import (
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"glogin/constant"
	"glogin/internal"
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
	Ctx *gin.Context
}

// http://127.0.0.1:8080/login/sms verify 获取验证码登陆 login 短信验证码登陆
func (l Login) SMSEx(request *glogin.SmsLoginReq, ctx *gin.Context) (response *glogin.SmsLoginRsp, err error) {
	l.Ctx = ctx
	return l.SMS(request)
}

func (l Login) SMS(request *glogin.SmsLoginReq) (response *glogin.SmsLoginRsp, err error) {
	response = &glogin.SmsLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
		ExtendData: &glogin.ExtendData{
			Authentication: &glogin.StateQueryResponse{},
		},
	}
	ip := l.Ctx.ClientIP()
	log.Infow("SMS login", "request", request, "ip", ip)
	// sms step verify 获得Phone 验证码
	if request.Step == constant.Verify {
		log.Infow("SMS verify", "request", request)
		code, err := sms.SmsVerify(request)
		if err != nil {
			response.Code = code
			response.Errmsg = fmt.Sprintf("get verify error %s ", request.Phone)
			return response, nil
		}
		return response, nil
		// sms step login  手机验证码登陆
	} else if request.Step == constant.Login {
		log.Infow("SMS login", "request", request)
		_, errCk := sms.CheckPhone(request.Phone, request.Verifycode)
		if errCk != nil {
			response.Code = constant.ErrCodeSMSCheckVerifyFail
			response.Errmsg = fmt.Sprintf("check verify faild phone %s ,verify %s", request.Phone, request.Verifycode)
			return response, nil
		}
		// 数美ID解析
		smID := smfpcrypto.ParseSMID(request.Client.Dhid)
		request.GetClient().Dhid = smID
		if account.CheckNotExist(bson.M{"phone": request.Phone}) {
			// 账号不存在,创建
			accountId, errCreate := account.CreatePhone(request, ip)
			if errCreate != nil {
				response.Code = constant.ErrCodeCreateInternal
				response.Errmsg = fmt.Sprintf("sms login %s create error: %s", request.Phone, errCreate)
				return response, errCreate
			}
			response.Code = constant.ErrCodeOk
			response.DhAccount = accountId
			response.SmId = smID
			response.Errmsg = "success"
			response.ExtendData.Nick = request.Phone
			response.DhToken = util.GenDHToken(accountId)
			log.Infow("sms login success", " request.phone", request.Phone, "rsp", response)
			return response, nil
		} else {
			// 账号存在, 直接登录
			loginRsp, errLogin := account.LoginPhone(request, ip)
			if errLogin != nil {
				response.Code = constant.ErrCodeLoginInternal
				response.Errmsg = fmt.Sprintf("sms login %s  error: %s", request.Phone, errLogin)
				return response, nil
			}
			// 返回InternalRsp
			rsp, ok := loginRsp.(internal.Rsp)
			if !ok {
				response.Code = constant.ErrCodeLoginInternal
				response.Errmsg = fmt.Sprintf("sms login %s  error: %s", request.Phone, errLogin)
				return response, nil
			}
			response.Code = constant.ErrCodeOk
			response.DhAccount = rsp.AccountData.ID
			response.SmId = smID
			response.DhToken = util.GenDHToken(rsp.AccountData.ID)
			response.ExtendData.Nick = rsp.AccountData.Phone
			response.Errmsg = "success"
			log.Infow("sms login success", " request.phone", request.Phone, "rsp", response)
			return response, nil
		}
	}
	return response, nil
}

func (l Login) ThirdEx(request *glogin.ThirdLoginReq, ctx *gin.Context) (response *glogin.ThridLoginRsp, err error) {
	l.Ctx = ctx
	return l.Third(request)
}
func (l Login) Third(request *glogin.ThirdLoginReq) (response *glogin.ThridLoginRsp, err error) {
	ip := l.Ctx.ClientIP()
	log.Infow("Third login", "request", request, "ip", ip)
	response = &glogin.ThridLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
		ExtendData: &glogin.ExtendData{
			Authentication: &glogin.StateQueryResponse{},
		},
	}
	uid, openId, dbField, errAuth := ThirdAuth(request)
	log.Infow("ThirdAuth Rsp", "uid", uid, "openid", openId)
	// 平台错误
	if errAuth == PlatIsWrong {
		response.Code = constant.ErrCodePlatWrong
		response.Errmsg = fmt.Sprintf("third %s plat is wrong", request.ThirdPlat)
		return response, nil
	}
	// 第三方账号验证失败
	if errAuth != nil {
		response.Code = constant.ErrCodeThirdAuthFail
		response.Errmsg = fmt.Sprintf("thrid plat %s auth error: %s", request.ThirdPlat, errAuth)
		return response, nil
	}
	// 数美ID解析
	smID := smfpcrypto.ParseSMID(request.Client.Dhid)
	request.GetClient().Dhid = smID
	if account.CheckNotExist(bson.M{dbField: openId}) {
		// 账号不存在,创建
		accountId, errCreate := account.CreateThird(request, dbField, openId, ip)
		if errCreate != nil {
			response.Code = constant.ErrCodeCreateInternal
			response.Errmsg = fmt.Sprintf("thrid plat %s login error: %s", request.ThirdPlat, errAuth)
			return response, errCreate
		}
		response.Code = 0
		response.DhToken = util.GenDHToken(accountId)
	} else {
		// 账号存在, 直接登录
		loginRsp, errLogin := account.LoginThird(request, dbField, uid, openId, ip)
		if errLogin != nil {
			response.Code = constant.ErrCodeLoginInternal
			response.Errmsg = fmt.Sprintf("thrid plat %s login error: %s", request.ThirdPlat, errAuth)
			return response, nil
		}
		// 返回
		rsp, ok := loginRsp.(internal.Rsp)
		if !ok {
			response.Code = constant.ErrCodeParsePbInternal
			response.Errmsg = fmt.Sprintf("thrid plat %s login error: %s", request.ThirdPlat, errLogin)
			return response, nil
		}
		log.Infow("thrid login success", " request.ThirdPlat", request.ThirdPlat, "openid", openId)
		response.Code = constant.ErrCodeOk
		response.DhAccount = rsp.AccountData.ID
		response.SmId = smID
		response.DhToken = util.GenDHToken(rsp.AccountData.ID)

		// ExtendData
		if request.ThirdPlat == "yedun" {
			response.ExtendData.Nick = openId
		}
		response.Errmsg = "success"
		return response, nil
	}
	return response, nil
}

func (l Login) VisitorEx(request *glogin.VisitorLoginReq, ctx *gin.Context) (response *glogin.VisitorLoginRsp, err error) {
	l.Ctx = ctx
	return l.Visitor(request)
}

func (l Login) Visitor(req *glogin.VisitorLoginReq) (rsp *glogin.VisitorLoginRsp, err error) {
	rsp = &glogin.VisitorLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
		ExtendData: &glogin.ExtendData{
			Authentication: &glogin.StateQueryResponse{},
		},
	}
	ip := l.Ctx.ClientIP()
	log.Infow("Visitor login", "request", req, "ip", ip)
	// 参数验证
	if req.Dhid == "" {
		rsp.Code = constant.ErrCodeParamError
		rsp.Errmsg = fmt.Sprintf("visitor is agrs error dhid: %s", req.Dhid)
		return rsp, nil
	}
	// 数美ID解析
	smId := smfpcrypto.ParseSMID(req.Client.Dhid)
	// 创建唯一设备账号，如果有SMID就用SMID创建，如果没有就随机
	visitorId := CreateVisitorID(req.Dhid, smId)
	if visitorId == "" {
		rsp.Code = constant.ErrCodeCreateVisitorIdFail
		rsp.Errmsg = fmt.Sprintf("visitor is agrs smid error: %s", req.Dhid)
		return rsp, nil
	}

	if account.CheckNotExist(bson.M{"visitor": visitorId}) {
		// 账号不存在,创建
		accountId, errCreate := account.CreateVisitor(req, visitorId, ip)
		if errCreate != nil {
			log.Infow("visitor fast login create account error ", "visitor", visitorId)
			rsp.Code = constant.ErrCodeCreateInternal
			rsp.Errmsg = fmt.Sprintf("visitor  %s fast login create account error: %s", visitorId, errCreate)
			return rsp, errCreate
		}
		rsp.Code = constant.ErrCodeOk
		rsp.DhToken = util.GenDHToken(accountId)
		rsp.SmId = smId
		rsp.DhAccount = accountId
		rsp.Visitor = visitorId
		rsp.ExtendData.Nick = ""
		rsp.Errmsg = "success"
		log.Infow("visitor fast login success ", "rsp", rsp)
		return rsp, nil
	} else {
		// 账号存在, 直接登录
		loginRsp, errLogin := account.LoginVisitor(req, visitorId, ip)
		if errLogin != nil {
			log.Infow("visitor fast login 1 error", "visitor", visitorId)
			rsp.Code = constant.ErrCodeLoginInternal
			rsp.Errmsg = fmt.Sprintf("visitor %s fast login error: %s", req.Dhid, errLogin)
			return rsp, nil
		}
		// 返回
		value, ok := loginRsp.(internal.Rsp)
		if !ok {
			log.Infow("visitor fast login  2 error", "visitor", visitorId)
			rsp.Code = constant.ErrCodeParsePbInternal
			rsp.Errmsg = fmt.Sprintf("visitor %s fast login error: %s", req.Dhid, errLogin)
			return rsp, nil
		}
		log.Infow("visitor login success", "visitor", visitorId)
		rsp.Code = constant.ErrCodeOk
		rsp.DhAccount = value.AccountData.ID
		rsp.SmId = smId
		rsp.DhToken = util.GenDHToken(value.AccountData.ID)
		rsp.Visitor = visitorId
		rsp.Errmsg = "success"
		return rsp, nil
	}
	return rsp, nil
}

func (l Login) FastEx(request *glogin.FastLoginReq, ctx *gin.Context) (response *glogin.FastLoginRsp, err error) {
	l.Ctx = ctx
	return l.Fast(request)
}

func (l Login) Fast(request *glogin.FastLoginReq) (response *glogin.FastLoginRsp, err error) {
	response = &glogin.FastLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
		ExtendData: &glogin.ExtendData{
			Authentication: &glogin.StateQueryResponse{},
		},
	}
	if request.DhToken == "" {
		response.Code = constant.ErrCodeFastTokenError
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
	ip := l.Ctx.ClientIP()
	// 数美ID解析
	smId := smfpcrypto.ParseSMID(request.Client.Dhid)
	loginRsp, errLogin := account.LoginFast(request, dhAccount, ip)
	if errLogin != nil {
		response.Code = constant.ErrCodeLoginInternal
		response.Errmsg = fmt.Sprintf("fast login %s auth error: %s", request.DhToken, errLogin)
		return response, nil
	}

	// 返回
	value, ok := loginRsp.(internal.Rsp)
	if !ok {
		response.Code = constant.ErrCodeParsePbInternal
		response.Errmsg = fmt.Sprintf("fast login %s auth error: %s", request.DhToken, errLogin)
		return response, nil
	}
	response.Code = constant.ErrCodeOk
	response.SmId = smId
	response.DhAccount = value.AccountData.ID
	response.DhToken = util.GenDHToken(value.AccountData.ID)
	response.Errmsg = "success"
	log.Infow("fast login success", "response", response)
	return response, nil
}

func ThirdAuth(request *glogin.ThirdLoginReq) (uid string, openid string, dbField string, err error) {
	thirdPlat := request.ThirdPlat
	if third, ok := plat.ThirdList[thirdPlat]; ok {
		uid, openid, err = third.Auth(request)
		dbField = third.DbFieldName()
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
