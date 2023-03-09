package cgi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	log "github.com/castbox/nirvana-gcore/glog"
	"glogin/constant"
	"glogin/internal"
	"glogin/internal/account"
	"glogin/internal/bilog"
	"glogin/internal/plat"
	"glogin/internal/smfpcrypto"
	"glogin/internal/sms"
	glogin "glogin/pbs/glogin"
	"glogin/util"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	PlatIsWrong = fmt.Errorf("third plat is wrong")
	VisitorMutex sync.Mutex
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
	ip := l.Ctx.ClientIP()
	reqId := uuid.New()
	log.Infow("SMS login request", "reqId", reqId, "request", request, "ip", ip)
	response = &glogin.SmsLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
		ExtendData: &glogin.ExtendData{
			Authentication: &glogin.StateQueryResponse{},
		},
	}
	before := time.Now().UnixNano()
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("got panic", "err", err)
		}
		log.Infow("SMS login rsp", "reqId", reqId, "request", request, "response", response, "time_cost", (time.Now().UnixNano()-before)/1000000)
	}()

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
		bCheck, checkCode, checkMsg := smsParamCheck(request)
		if !bCheck {
			response.Code = checkCode
			response.Errmsg = checkMsg
			return response, nil
		}
		request.Client.Ip = ip
		// 短信认证
		_, errCk := sms.CheckPhone(request.Phone, request.Verifycode)
		if errCk != nil {
			response.Code = constant.ErrCodeSMSCheckVerifyFail
			response.Errmsg = fmt.Sprintf("check verify faild phone %s ,verify %s", request.Phone, request.Verifycode)
			return response, nil
		}
		// 数美ID解析
		smID := smfpcrypto.ParseSMID(request.Client.Dhid)
		request.GetClient().Dhid = smID
		response.SmId = smID
		if account.CheckNotExist(bson.M{"phone": request.Phone}) {
			// 账号不存在,创建
			createRsp, errCreate := account.CreatePhone(request, ip)
			if errCreate != nil {
				response.Code = constant.ErrCodeCreateInternal
				response.Errmsg = fmt.Sprintf("sms login1 %s create error: %s", request.Phone, errCreate)
				return response, errCreate
			}
			// 返回InternalRsp
			smsResponse(response, createRsp)
			response.ExtendData.Nick = util.HideStar(request.Phone)
			log.Infow("sms login success", " request.phone", request.Phone, "rsp", response)
			bilog.SmsLogin(request, util.Int642String(int64(response.DhAccount)))
			return response, nil
		} else {
			// 账号存在, 直接登录
			loginRsp, errLogin := account.LoginPhone(request, ip)
			if errLogin != nil {
				response.Code = constant.ErrCodeLoginInternal
				response.Errmsg = fmt.Sprintf("sms login3 %s  error: %s", request.Phone, errLogin)
				return response, nil
			}
			// 返回InternalRsp
			smsResponse(response, loginRsp)
			response.ExtendData.Nick = util.HideStar(request.Phone)
			log.Infow("sms login success", "request.phone", request.Phone, "rsp", response)
			bilog.SmsLogin(request, util.Int642String(int64(response.DhAccount)))
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
	reqId := uuid.New()
	log.Infow("Third login request", "reqId", reqId, "request", request, "ip", ip)
	response = &glogin.ThridLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
		ExtendData: &glogin.ExtendData{
			Authentication: &glogin.StateQueryResponse{},
		},
	}
	before := time.Now().UnixNano()
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("got panic", "err", err)
		}
		log.Infow("Third login rsp", "reqId", reqId, "response", response, "time_cost", (time.Now().UnixNano()-before)/1000000)
	}()

	// 必传参数验证
	bCheck, checkCode, checkMsg := thirdParamCheck(request)
	if !bCheck {
		response.Code = checkCode
		response.Errmsg = checkMsg
		return response, nil
	}

	request.Client.Ip = ip

	if request.ThirdToken == "editor_test_token"{
		// 游客验证
		visitorId := VisitorID(request.Client.Dhid)
		log.Infow("visitor Id","dhid: ", request.Client.Dhid, "visitorId: ", visitorId)
		if visitorId == "" {
			response.Code = constant.ErrCodeCreateVisitorIdFail
			response.Errmsg = fmt.Sprintf("visitor is agrs smid error: %s", request.Client.Dhid)
			return response, nil
		}

		if account.CheckNotExist(bson.M{"visitor": visitorId}) {
			// 账号不存在,创建
			createRsp, errCreate := account.CreateThird(request, "visitor", visitorId, ip)
			if errCreate != nil {
				log.Infow("visitor fast login create account error ", "visitor", visitorId)
				response.Code = constant.ErrCodeCreateInternal
				response.Errmsg = fmt.Sprintf("visitor  %s fast login create account error: %s", visitorId, errCreate)
				return response, errCreate
			}
			response.SmId = visitorId
			thirdResponse(response, createRsp)
			// token存储
			account.SetToken(response.DhAccount, response.DhToken)
			log.Infow("visitor fast login success ", "response", response)
			return response, nil
		} else {
			// 账号存在, 直接登录
			loginRsp, errLogin := account.LoginThird(request, "visitor", visitorId, ip)
			if errLogin != nil {
				log.Infow("visitor fast login 1 error", "visitor", visitorId)
				response.Code = constant.ErrCodeLoginInternal
				response.Errmsg = fmt.Sprintf("visitor %s fast login error: %s", request.Client.Dhid, errLogin)
				return response, nil
			}
			// 返回
			response.SmId = visitorId
			thirdResponse(response, loginRsp)
			// token存储
			account.SetToken(response.DhAccount, response.DhToken)
			log.Infow("visitor login success", "visitor", visitorId)
			return response, nil
		}
	}else{
		// 第三方认证
		authRsp, dbField, errAuth := ThirdAuth(request)
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
		response.SmId = smID

		/*
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
					response.SmId = smID
					response.ExtendData.Nick = authRsp.Nick
					thirdResponse(response, loginRsp)
					log.Infow("third bundle account login success", "response", response, "uid", uid)
					bilog.ThirdLogin(request, util.Int642String(int64(response.DhAccount)))
					return response, nil
				}
			}
		*/

		// 卓杭通行证
		if account.CheckNotExist(bson.M{dbField: unionId}) {
			// 账号不存在,创建
			createRsp, errCreate := account.CreateThird(request, dbField, unionId, ip)
			if errCreate != nil {
				response.Code = constant.ErrCodeCreateInternal
				response.Errmsg = fmt.Sprintf("thrid plat 1 %s login error: %s", request.ThirdPlat, errCreate)
				return response, errCreate
			}
			thirdResponse(response, createRsp)
			response.ExtendData.Nick = authRsp.Nick
			if request.ThirdPlat == "yedun" {
				response.ExtendData.Nick = util.HideStar(unionId)
			}
			log.Infow("third login success 1", " request.ThirdPlat", request.ThirdPlat, "unionId", unionId)
			//bilog.ThirdLogin(request, util.Int642String(int64(response.DhAccount)))
			return response, nil
		} else {
			// 账号存在, 直接登录
			loginRsp, errLogin := account.LoginThird(request, dbField, unionId, ip)
			if errLogin != nil {
				response.Code = constant.ErrCodeLoginInternal
				response.Errmsg = fmt.Sprintf("thrid plat 3 %s login error: %s", request.ThirdPlat, errLogin)
				return response, nil
			}
			// 返回
			thirdResponse(response, loginRsp)
			response.ExtendData.Nick = authRsp.Nick
			if request.ThirdPlat == "yedun" {
				response.ExtendData.Nick = util.HideStar(unionId)
			}
			log.Infow("third login success 2", "response", response, "unionId", unionId)
			//bilog.ThirdLogin(request, util.Int642String(int64(response.DhAccount)))
			return response, nil
		}
	}
}

func (l Login) VisitorEx(request *glogin.VisitorLoginReq, ctx *gin.Context) (response *glogin.VisitorLoginRsp, err error) {
	l.Ctx = ctx
	return l.Visitor(request)
}

func (l Login) Visitor(req *glogin.VisitorLoginReq) (rsp *glogin.VisitorLoginRsp, err error) {
	VisitorMutex.Lock()
	defer VisitorMutex.Unlock()

	rsp = &glogin.VisitorLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
		ExtendData: &glogin.ExtendData{
			Authentication: &glogin.StateQueryResponse{},
		},
	}
	reqId := uuid.New()
	log.Infow("Visitor login", "reqId", reqId, "request", req)
	before := time.Now().UnixNano()
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("got panic", "err", err)
		}
		log.Infow("Visitor login rsp", "reqId", reqId, "response", rsp, "time_cost", (time.Now().UnixNano()-before)/1000000)
	}()

	// 必传参数验证
	bCheck, checkCode, checkMsg := visitorParamCheck(req)
	if !bCheck {
		rsp.Code = checkCode
		rsp.Errmsg = checkMsg
		return rsp, nil
	}
	ip := l.Ctx.ClientIP()
	req.Client.Ip = ip
	// 游客验证
	visitorId := VisitorID(req.Dhid)
	log.Infow("visitor Id","dhid: ", req.Dhid, "visitorId: ", visitorId)
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
		rsp.SmId = visitorId
		visitorResponse(rsp, createRsp)
		// token存储
		account.SetToken(rsp.DhAccount, rsp.DhToken)
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
		rsp.SmId = visitorId
		visitorResponse(rsp, loginRsp)
		// token存储
		account.SetToken(rsp.DhAccount, rsp.DhToken)
		log.Infow("visitor login success", "visitor", visitorId)
		return rsp, nil
	}
	return rsp, nil
}

func (l Login) FastEx(request *glogin.FastLoginReq, ctx *gin.Context) (response *glogin.FastLoginRsp, err error) {
	l.Ctx = ctx
	return l.Fast(request)
}

func (l Login) Fast(request *glogin.FastLoginReq) (response *glogin.FastLoginRsp, err error) {
	reqId := uuid.New()
	log.Infow("fast login request", "reqId", reqId, "request", request)
	response = &glogin.FastLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
		ExtendData: &glogin.ExtendData{
			Authentication: &glogin.StateQueryResponse{},
		},
	}
	before := time.Now().UnixNano()
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("got panic", "err", err)
		}
		log.Infow("fast login rsp", "reqId", reqId, "response", response, "time_cost", (time.Now().UnixNano()-before)/1000000)
	}()

	// 必传参数验证
	bCheck, checkCode, checkMsg := fastParamCheck(request)
	if !bCheck {
		response.Code = checkCode
		response.Errmsg = checkMsg
		return response, nil
	}
	ip := l.Ctx.ClientIP()
	request.Client.Ip = ip
	// 校验token
	dhAccount, errToken := util.ValidDHToken(request.DhToken)
	if errToken != nil {
		if errToken == util.ExpiredToken {
			// 1:快速登录失败，token 过期 30天
			response.Code = constant.ErrCodeFastTokenExpired
			response.Errmsg = fmt.Sprintf("fast login3: %s", errToken)
			return response, nil
		} else {
			response.Code = constant.ErrCodeFastTokenVaild
			response.Errmsg = fmt.Sprintf("fast login4: %s", errToken)
			return response, nil
		}
	}

	// 数美ID解析
	smId := smfpcrypto.ParseSMID(request.Client.Dhid)
	response.SmId = smId

	// 登录验证
	loginRsp, errLogin := account.LoginFast(request, dhAccount, ip)
	if errLogin != nil {
		response.Code = constant.ErrCodeLoginInternal
		response.Errmsg = fmt.Sprintf("fast login1 %s auth error: %s", request.DhToken, errLogin)
		log.Warnw("Fast Login error ", "errLogin", errLogin)
		return response, nil
	}

	// 返回赋值
	fastResponse(response, loginRsp)
	log.Infow("fast login success", "response", response)
	//bilog.FastLogin(request)
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

// 获取visitor id
func VisitorID(srcVisitor string) (ret string) {
	if srcVisitor == "0" {
		Count := account.AccountCount(bson.M{"_id": bson.M{"$gt": 0}})
		if Count == -1 {
			log.Warnw("VisitorID err: Count = -1")
			ret = ""
		}else {
			ret = strconv.FormatInt(Count, 10)
		}
	}
	return
}

func fastParamCheck(request *glogin.FastLoginReq) (bRsp bool, code int32, errMsg string) {
	if request.Game.GameCd == "" || request.Game.Channel == "" {
		code = constant.ErrCodeParamError
		errMsg = fmt.Sprintf("fast login req param error GameCd: %s Channel: %s", request.Game.GameCd, request.Game.Channel)
		return false, code, errMsg
	}

	if request.DhToken == "" {
		code = constant.ErrCodeFastTokenError
		errMsg = fmt.Sprintf("fast is token null game: %s client: %s", request.Game, request.Client)
		return false, code, errMsg
	}
	return true, code, errMsg
}

func fastResponse(response *glogin.FastLoginRsp, loginRsp internal.Rsp) {
	response.Code = constant.ErrCodeOk
	response.Errmsg = "success"
	response.DhAccount = loginRsp.AccountData.ID
	response.DhToken = util.GenDHToken(loginRsp.AccountData.ID)
	// 防沉迷返回
	//r2, ok := loginRsp.AntiRsp.(pb_obsession.CheckStateQueryResponse)
	//if ok {
	//	response.ExtendData.Authentication = antiConvertClient(&r2)
	//}
	if loginRsp.AccountData.Phone != "" {
		response.ThirdPlat = "phone"
	} else {
		plat, errG := account.GetPlat(loginRsp.AccountData)
		if errG == nil {
			response.ThirdPlat = plat
		}
	}
	response.ExtendData.GameFirstLogin = loginRsp.GameRsp.FirstLogin
}

func visitorParamCheck(request *glogin.VisitorLoginReq) (bRsp bool, code int32, errMsg string) {
	if request.Dhid == "" || request.Game.GameCd == "" || request.Game.Channel == "" {
		code = constant.ErrCodeParamError
		errMsg = fmt.Sprintf("visitor is agrs error dhid: %s GameCd: %s Channel: %s", request.Dhid, request.Game.GameCd, request.Game.Channel)
		return false, code, errMsg
	}
	return true, code, errMsg
}

func visitorResponse(rsp *glogin.VisitorLoginRsp, dcRsp internal.Rsp) {
	rsp.Code = constant.ErrCodeOk
	rsp.DhToken = util.GenDHToken(dcRsp.AccountData.ID)
	rsp.DhAccount = dcRsp.AccountData.ID
	rsp.Visitor = dcRsp.AccountData.Visitor
	// 防沉迷返回
	//r2, ok := dcRsp.AntiRsp.(pb_obsession.CheckStateQueryResponse)
	//if ok {
	//	rsp.ExtendData.Authentication = antiConvertClient(&r2)
	//}
	rsp.ExtendData.Nick = ""
	rsp.ExtendData.GameFirstLogin = dcRsp.GameRsp.FirstLogin
	rsp.Errmsg = "success"
}

func thirdParamCheck(request *glogin.ThirdLoginReq) (bRsp bool, code int32, errMsg string) {
	if request.Game.GameCd == "" || request.Game.Channel == "" || request.Client.Dhid == "" {
		code = constant.ErrCodeParamError
		errMsg = fmt.Sprintf("third login req param error GameCd: %s Channel: %s Client.Dhid: %s", request.Game.GameCd, request.Game.Channel, request.Client.Dhid)
		return false, code, errMsg
	}
	return true, code, errMsg
}

func thirdResponse(response *glogin.ThridLoginRsp, dcRsp internal.Rsp) {
	response.Code = constant.ErrCodeOk
	response.DhAccount = dcRsp.AccountData.ID
	response.Errmsg = "success"
	response.DhToken = util.GenDHToken(dcRsp.AccountData.ID)
	//r2, ok := dcRsp.AntiRsp.(*pb_obsession.CheckStateQueryResponse)
	//if ok {
	//	response.ExtendData.Authentication = antiConvertClient(r2)
	//}
	response.ExtendData.GameFirstLogin = dcRsp.GameRsp.FirstLogin
}

func smsParamCheck(request *glogin.SmsLoginReq) (bRsp bool, code int32, errMsg string) {
	if request.Client == nil || request.Phone == "" || request.Game.GameCd == "" || request.Game.Channel == "" || request.Client.Dhid == "" {
		code = constant.ErrCodeParamError
		errMsg = fmt.Sprintf("sms login req param error GameCd: %s Channel: %s Client.Dhid: %s", request.Game.GameCd, request.Game.Channel, request.Client.Dhid)
		return false, code, errMsg
	}
	return true, code, errMsg
}

func smsResponse(response *glogin.SmsLoginRsp, dcRsp internal.Rsp) {
	response.Code = constant.ErrCodeOk
	response.DhAccount = dcRsp.AccountData.ID
	response.Errmsg = "success"
	response.DhToken = util.GenDHToken(dcRsp.AccountData.ID)
	//r2, ok := dcRsp.AntiRsp.(*pb_obsession.CheckStateQueryResponse)
	//if ok {
	//	response.ExtendData.Authentication = antiConvertClient(r2)
	//}
	response.ExtendData.GameFirstLogin = dcRsp.GameRsp.FirstLogin
}
