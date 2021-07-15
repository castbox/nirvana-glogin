package cgi

import (
	"fmt"
	"git.dhgames.cn/svr_comm/gmoss/v2"
	_ "github.com/gin-gonic/gin"
	"glogin/constant"
	"glogin/internal/plat"
	glogin2 "glogin/pbs/glogin"
	"glogin/pbs/hawkeye_login"
)

var (
	PlatIsWrong = fmt.Errorf("third plat is wrong")
)

type Login struct {
}

func (Login) SMS(request *glogin2.SmsLoginReq) (response *glogin2.SmsLoginRsp, err error) {
	response = &glogin2.SmsLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
	}
	return response, nil
}

func (Login) Third(request *glogin2.ThirdLoginReq) (response *glogin2.ThridLoginRsp, err error) {
	response = &glogin2.ThridLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
	}
	_, errAuth := ThirdAuth(request.ThirdPlat, request.Game.BundleId, request.ThirdToken)
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

	req := &hawkeye_login.Login{
		GameCd:     "1004",
		DeviceInfo: nil,
		UserInfo:   nil,
	}
	//call hawkeye
	cluster := gmoss.MossWithClusterServiceIndex("yanghaitao_dev", "hawkeye", 1)
	_, err = hawkeye_login.HawkeyeLogin(cluster, req, gmoss.Call, gmoss.DefaultCallOption())
	if err != nil {
		return nil, err
	}
	//
	//if account.Exist(bson.M{loginInfo.ThirdPlat: uid, "bundle_id": loginInfo.Game.BundleId}) {
	//	// 账号存在，登录
	//	accountId, errLogin := account.LoginThird(loginInfo.ThirdPlat, uid, loginInfo.Game.BundleId, loginInfo.Client.DHId, c.ClientIP(), loginInfo.Game)
	//	if errLogin != nil {
	//		// 登录内部错误
	//		c.JSON(200, gin.H{
	//			"code":   3,
	//			"errmsg": fmt.Sprintf("plat %s, uid:%s, login error: %s", loginInfo.ThirdPlat, uid, errLogin),
	//		})
	//		return
	//	} else {
	//		logger.Infof("plat %s, uid:%s, accountId:%d, ip:%s, third login success.", loginInfo.ThirdPlat, uid, accountId, c.ClientIP())
	//		c.JSON(200, gin.H{
	//			"code":     0,
	//			"dh_token": account.GenDHToken(accountId),
	//		})
	//	}
	//} else {
	//	// 账号不存在，注册
	//	accountId, errCreate := account.ThirdCreate(loginInfo.ThirdPlat, uid, loginInfo.Game.BundleId, loginInfo.Client.DHId, c.ClientIP(), loginInfo.Game)
	//	if errCreate != nil {
	//		// 登录内部错误
	//		c.JSON(200, gin.H{
	//			"code":   2,
	//			"errmsg": fmt.Sprintf("plat %s, uid:%s, creat account error: %s", loginInfo.ThirdPlat, uid, errCreate),
	//		})
	//	} else {
	//		logger.Infof("plat %s, uid:%s, create accountId:%d, ip:%s, third first login success.", loginInfo.ThirdPlat, uid, accountId, c.ClientIP())
	//		httputil.ReplaySuccess(c, accountId, "")
	//	}
	//}
	return response, nil
}

func (Login) Visitor(request *glogin2.VisitorLoginReq) (response *glogin2.VisitorLoginRsp, err error) {
	response = &glogin2.VisitorLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
	}
	return response, nil
}

func ThirdAuth(thirdPlat string, bundleId string, thirdToken string) (uid string, err error) {
	if third, ok := plat.ThirdList[thirdPlat]; ok {
		uid, err = third.Auth(bundleId, thirdToken)
	} else {
		err = PlatIsWrong
	}
	return
}
