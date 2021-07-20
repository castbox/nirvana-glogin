package sms

import (
	"glogin/pbs/glogin"
	"math/rand"
)

func SmsVerify(req *glogin.SmsLoginReq) (bool, error) {
	return false, nil
}

func canSendSmsVerify(phone string) (bool, error) {
	return false, nil
}

func sendSmsVerify(phone string, verifyCode string) (bool, error) {
	return false, nil
}

func createVerifyCode() string {
	var Min int32 = 100000
	var Max int32 = 999999
	verifyCode := rand.Int31n(Max-Min) + Min
	return string(verifyCode)
}

func smsSign() string {
	return ""
}

func CheckPhone(phone string, verifyCode string) (bool, error) {
	return false, nil
}
