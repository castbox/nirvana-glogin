package sms

import (
	"encoding/json"
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"glogin/config"
	"glogin/constant"
	"glogin/db"
	"glogin/internal/xhttp"
	"glogin/pbs/glogin"
	"glogin/util"
	"time"
)

func GetVerify(req *glogin.SmsLoginReq) (int32, error) {
	phone := req.Phone
	log.Infow("SmsVerify", "phone", phone)
	code, err := canSendVerify(phone)
	if err != nil {
		return code, err
	}
	verifyCode := CreateVerifyCode()
	bSend, err := sendVerify(phone, verifyCode)
	if false == bSend {
		return constant.ErrCodeSmsFail, err
	}
	_, errAdd := db.AddSmsVerify(phone, verifyCode)
	if errAdd != nil {
		log.Errorw("DbAddSmsVerify", "errAdd", errAdd)
		return constant.ErrCodeDB, errAdd
	}
	return constant.ErrCodeOk, nil
}

// 检查是否能发送
func canSendVerify(phone string) (int32, error) {
	bCheck, errCheck := db.CheckSmsInterval(phone)
	if !bCheck {
		log.Errorw("CheckSmsInterval false", "errCheck", errCheck)
		return constant.ErrCodeSmsInterval, errCheck
	}
	_, errCount := db.CheckSmsVerifyCount(phone)
	if errCount != nil {
		return constant.ErrCodeSmsCount, nil
	}
	return 0, nil
}

// 发送验证码
func sendVerify(phone string, verifyCode string) (bool, error) {
	smsUrl := config.Field("sms_url").String()
	timeStamp := util.FormatDate(time.Now(), util.YYYYMMDDHHMMSS)
	content := config.Field("sms_content").String() + verifyCode
	httpClient := xhttp.NewClient().Type(xhttp.TypeUrlencoded)
	reqBody := xhttp.BodyMap{}
	reqBody.Set("appId", config.Field("sms_appid").String())
	reqBody.Set("timestamp", timeStamp)
	reqBody.Set("sign", smsSign(timeStamp))
	reqBody.Set("mobiles", phone)
	reqBody.Set("content", content)
	log.Infow("begin sendSmsVerify", "sendSmsUrl", smsUrl, "reqBody", reqBody)
	res, bs, errs := httpClient.Post(smsUrl).SendBodyMap(reqBody).EndBytes()
	if len(errs) > 0 {
		log.Warnw("sendSmsVerify Error1", "errs", errs[0], "sendSmsUrl", smsUrl)
		return false, errs[0]
	}
	if res.StatusCode != 200 {
		log.Warnw("sendSmsVerify Error2,", "StatusCode", res.StatusCode)
		return false, fmt.Errorf("sendSmsVerify error statuscode: %v", res.StatusCode)
	}
	log.Infow("sendSmsVerify HTTP Rsp,", "string(bs)", string(bs))
	smsRsp := new(SendSmsRsp)
	if err := json.Unmarshal(bs, smsRsp); err != nil {
		log.Warnw("sendSmsVerify HTTP Rsp Error3,", "Unmarshal", string(bs))
		return false, fmt.Errorf("sendSmsVerify error Unmarshal: %v", string(bs))
	}
	if smsRsp.Code == "SUCCESS" {
		return true, nil
	} else {
		log.Warnw("ParseSMID ParseBoxData HTTP Request Error4,", "CodeErr", smsRsp.Code)
		return false, fmt.Errorf("sendSmsVerify error smsRsp.Code: %v", smsRsp.Code)
	}
	return false, fmt.Errorf("sendSmsVerify error 5")
}

// 生成验证码
func CreateVerifyCode() string {
	var Min int32 = 100000
	var Max int32 = 999999
	verifyCode := util.Rand32Num(Min, Max)
	strCode := fmt.Sprintf("%d", verifyCode)
	return strCode
}

func smsSign(timeStamp string) string {
	appId := config.Field("sms_appid").String()
	secret := config.Field("sms_secret").String()
	Data := appId + secret + timeStamp
	return util.Md5String(Data)
}

func CheckPhone(phone string, verifyCode string) (bool, error) {
	return db.CheckSmsVerifyCode(phone, verifyCode)
}
