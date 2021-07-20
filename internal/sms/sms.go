package sms

import (
	"encoding/json"
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"glogin/config"
	"glogin/db"
	"glogin/internal/xhttp"
	"glogin/pbs/glogin"
	"glogin/util"
	"math/rand"
	"time"
)

func SmsVerify(req *glogin.SmsLoginReq) (bool, error) {
	phone := req.Phone
	log.Infow("SmsVerify", "phone", phone)
	_, err := canSendSmsVerify(phone)
	if err != nil {
		return false, err
	}
	verifyCode := createVerifyCode()
	bSend, err := sendSmsVerify(phone, verifyCode)
	if bSend {
		_, errAdd := db.AddSmsVerify(phone, verifyCode)
		if errAdd != nil {
			log.Errorw("AddSmsVerify", "errAdd", errAdd)
			return false, errAdd
		}
		return true, nil
	}
	return false, nil
}

// 检查是否能发送
func canSendSmsVerify(phone string) (bool, error) {
	bCheck, errCheck := db.CheckSmsInterval(phone)
	if !bCheck {
		log.Errorw("CheckSmsInterval false", "errCheck", errCheck)
		return false, errCheck
	}
	_, errCount := db.CheckSmsVerifyCount(phone)
	if errCount != nil {
		return false, nil
	}
	return true, nil
}

// 发送验证码
func sendSmsVerify(phone string, verifyCode string) (bool, error) {
	/*%% sms content
	Content = list_to_binary(binary_to_list(consul_config:sms_content()) ++ binary_to_list(VerifyCode)),
	Req = #{
		<<"appId">> => consul_config:sms_appid(),
		<<"timestamp">> => TimeStamp,
		<<"sign">> => sms_sign(TimeStamp),
		<<"mobiles">> => Phone,
		<<"content">> => Content
	},
	Body = map_to_url_param(Req),
		ibrowse:send_req(?SMS_URL, [{"Content-Type", "application/x-www-form-urlencoded"}], post, Body).
	*?
	*/
	smsUrl := config.Field("sms_url").String()
	timeStamp := util.FormatDate(time.Now(), util.YYYYMMDDHHMMSS)
	content := config.Field("sms_content").String() + verifyCode
	httpClient := xhttp.NewClient().Type(xhttp.TypeUrlencoded)
	reqBody := xhttp.BodyMap{}
	reqBody.Set("appId", config.Field("sms_appid"))
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
func createVerifyCode() string {
	var Min int32 = 100000
	var Max int32 = 999999
	verifyCode := rand.Int31n(Max-Min) + Min
	return string(verifyCode)
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
