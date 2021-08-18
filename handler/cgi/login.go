package cgi

import (
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"glogin/config"
	"glogin/constant"
	"glogin/internal"
	"glogin/internal/account"
	"glogin/internal/plat"
	"glogin/internal/smfpcrypto"
	"glogin/internal/sms"
	anti_authentication "glogin/pbs/authentication"
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
	// sms step verify 获得Phone验证码
	if request.Step == constant.Verify {
		log.Infow("SMS verify", "request", request)
		code, err := sms.GetVerify(request)
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
			createRsp, errCreate := account.CreatePhone(request, ip)
			if errCreate != nil {
				response.Code = constant.ErrCodeCreateInternal
				response.Errmsg = fmt.Sprintf("sms login %s create error: %s", request.Phone, errCreate)
				return response, errCreate
			}
			// 返回InternalRsp
			rsp, ok := createRsp.(internal.Rsp)
			if !ok {
				response.Code = constant.ErrCodeCreateInternal
				response.Errmsg = fmt.Sprintf("sms login %s  error: %v", request.Phone, rsp)
				return response, nil
			}
			response.Code = constant.ErrCodeOk
			response.DhAccount = rsp.AccountData.ID
			response.SmId = smID
			response.Errmsg = "success"
			response.ExtendData.Nick = util.HideStar(request.Phone)
			r2, ok := rsp.AntiRsp.(*anti_authentication.StateQueryResponse)
			if ok {
				response.ExtendData.Authentication = (*glogin.StateQueryResponse)(r2)
			}
			response.DhToken = util.GenDHToken(rsp.AccountData.ID)
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
			response.ExtendData.Nick = util.HideStar(rsp.AccountData.Phone)
			r2, ok := rsp.AntiRsp.(*anti_authentication.StateQueryResponse)
			if ok {
				response.ExtendData.Authentication = (*glogin.StateQueryResponse)(r2)
			}
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
	authRsp, dbField, errAuth := ThirdAuth(request)
	// 平台错误
	if errAuth == PlatIsWrong {
		response.Code = constant.ErrCodePlatWrong
		response.Errmsg = fmt.Sprintf("third %s plat is wrong", request.ThirdPlat)
		return response, nil
	}
	// 第三方账号验证失败
	if errAuth != nil || authRsp == nil {
		response.Code = constant.ErrCodeThirdAuthFail
		response.Errmsg = fmt.Sprintf("thrid plat %s auth error: %s", request.ThirdPlat, errAuth)
		return response, nil
	}
	uid := authRsp.Uid
	unionId := authRsp.UnionId
	log.Infow("ThirdAuth Rsp", "uid", uid, "unionId", unionId)
	// 数美ID解析
	smID := smfpcrypto.ParseSMID(request.Client.Dhid)
	request.GetClient().Dhid = smID
	// 为兼容海外版本老数据,1004项目若获得bundle账号，直接登录
	if config.Field("region_mark").Int() == constant.RegionOverseas &&
		request.Game.GameCd == constant.AODGameCd {
		if !account.CheckNotExist(bson.M{dbField: uid, "bundle_id": request.Game.BundleId}) {
			loginRsp, errLogin := account.LoginBundleThird(request, dbField, uid, ip)
			if errLogin != nil {
				response.Code = constant.ErrCodeLoginInternal
				response.Errmsg = fmt.Sprintf("thrid plat %s bundle account login error: %s", request.ThirdPlat, errAuth)
				return response, nil
			}
			rsp, ok := loginRsp.(internal.Rsp)
			if !ok {
				response.Code = constant.ErrCodeParsePbInternal
				response.Errmsg = fmt.Sprintf("thrid plat %s bundle account login error: %s", request.ThirdPlat, errLogin)
				return response, nil
			}
			response.Code = constant.ErrCodeOk
			response.DhAccount = rsp.AccountData.ID
			response.Errmsg = "success"
			response.SmId = smID
			response.DhToken = util.GenDHToken(rsp.AccountData.ID)
			r2, ok := rsp.AntiRsp.(*anti_authentication.StateQueryResponse)
			if ok {
				response.ExtendData.Authentication = (*glogin.StateQueryResponse)(r2)
			}
			response.ExtendData.Nick = authRsp.Nick
			log.Infow("third bundle account login success", "response", response, "uid", uid)
			return response, nil
		}
	}

	// 卓杭通行证
	if account.CheckNotExist(bson.M{dbField: unionId}) {
		// 账号不存在,创建
		createRsp, errCreate := account.CreateThird(request, dbField, unionId, ip)
		if errCreate != nil {
			response.Code = constant.ErrCodeCreateInternal
			response.Errmsg = fmt.Sprintf("thrid plat %s login error: %s", request.ThirdPlat, errAuth)
			return response, errCreate
		}

		rsp, ok := createRsp.(internal.Rsp)
		if !ok {
			response.Code = constant.ErrCodeCreateInternal
			response.Errmsg = fmt.Sprintf("thrid plat %s login error: %s", request.ThirdPlat, rsp)
			return response, nil
		}

		log.Infow("third login success 1", " request.ThirdPlat", request.ThirdPlat, "unionId", unionId)
		response.Code = constant.ErrCodeOk
		response.DhAccount = rsp.AccountData.ID
		response.SmId = smID
		response.DhToken = util.GenDHToken(rsp.AccountData.ID)
		// 防沉迷返回
		r2, ok := rsp.AntiRsp.(*anti_authentication.StateQueryResponse)
		if ok {
			response.ExtendData.Authentication = (*glogin.StateQueryResponse)(r2)
		}
		// ExtendData
		response.ExtendData.Nick = authRsp.Nick
		if request.ThirdPlat == "yedun" {
			response.ExtendData.Nick = util.HideStar(unionId)
		}
		response.Errmsg = "success"
		return response, nil
	} else {
		// 账号存在, 直接登录
		loginRsp, errLogin := account.LoginThird(request, dbField, unionId, ip)
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
		//	log.Infow("thrid login success", " request.ThirdPlat", request.ThirdPlat, "unionId", unionId)
		response.Code = constant.ErrCodeOk
		response.DhAccount = rsp.AccountData.ID
		response.SmId = smID
		response.DhToken = util.GenDHToken(rsp.AccountData.ID)
		// 防沉迷返回
		r2, ok := rsp.AntiRsp.(*anti_authentication.StateQueryResponse)
		if ok {
			response.ExtendData.Authentication = (*glogin.StateQueryResponse)(r2)
		}
		response.ExtendData.Nick = authRsp.Nick
		if request.ThirdPlat == "yedun" {
			response.ExtendData.Nick = util.HideStar(unionId)
		}
		response.Errmsg = "success"
		log.Infow("third login success 2", "response", response, "unionId", unionId)
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
		createRsp, errCreate := account.CreateVisitor(req, visitorId, ip)
		if errCreate != nil {
			log.Infow("visitor fast login create account error ", "visitor", visitorId)
			rsp.Code = constant.ErrCodeCreateInternal
			rsp.Errmsg = fmt.Sprintf("visitor  %s fast login create account error: %s", visitorId, errCreate)
			return rsp, errCreate
		}

		dcRsp, ok := createRsp.(internal.Rsp)
		if !ok {
			rsp.Code = constant.ErrCodeCreateInternal
			rsp.Errmsg = fmt.Sprintf("thrid plat %s login error: %s", "visitor", rsp)
			return rsp, nil
		}
		rsp.Code = constant.ErrCodeOk
		rsp.DhToken = util.GenDHToken(dcRsp.AccountData.ID)
		rsp.SmId = smId
		rsp.DhAccount = dcRsp.AccountData.ID
		rsp.Visitor = visitorId
		// 防沉迷返回
		r2, ok := dcRsp.AntiRsp.(anti_authentication.StateQueryResponse)
		if ok {
			rsp.ExtendData.Authentication = (*glogin.StateQueryResponse)(&r2)
		}
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
		// 防沉迷返回
		r2, ok := value.AntiRsp.(anti_authentication.StateQueryResponse)
		if ok {
			rsp.ExtendData.Authentication = (*glogin.StateQueryResponse)(&r2)
		}
		return rsp, nil
	}
	return rsp, nil
}

func (l Login) FastEx(request *glogin.FastLoginReq, ctx *gin.Context) (response *glogin.FastLoginRsp, err error) {
	l.Ctx = ctx
	return l.Fast(request)
}

func (l Login) Fast(request *glogin.FastLoginReq) (response *glogin.FastLoginRsp, err error) {
	log.Infow("fast login request", "request", request)
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
	// 防沉迷返回
	r2, ok := value.AntiRsp.(*anti_authentication.StateQueryResponse)
	if ok {
		response.ExtendData.Authentication = (*glogin.StateQueryResponse)(r2)
	}

	log.Infow("fast login success", "response", response)
	return response, nil
}

func ThirdAuth(request *glogin.ThirdLoginReq) (info *plat.AuthRsp, dbField string, err error) {
	thirdPlat := request.ThirdPlat
	if third, ok := plat.ThirdList[thirdPlat]; ok {
		info, err = third.Auth(request)
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
