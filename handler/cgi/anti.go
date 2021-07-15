package cgi

import (
	"glogin/constant"
	glogin2 "glogin/pbs/glogin"
)

type Anti struct {
}

func (Anti) Query(request *glogin2.SmsLoginReq) (response *glogin2.SmsLoginRsp, err error) {
	response = &glogin2.SmsLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
	}
	return response, nil
}

func (Anti) Check(request *glogin2.SmsLoginReq) (response *glogin2.SmsLoginRsp, err error) {
	response = &glogin2.SmsLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
	}
	return response, nil
}
