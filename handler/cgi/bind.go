package cgi

import (
	"glogin/constant"
	glogin2 "glogin/pbs/glogin"
)

type Bind struct {
}

func (Bind) Third(request *glogin2.SmsLoginReq) (response *glogin2.SmsLoginRsp, err error) {
	response = &glogin2.SmsLoginRsp{
		Code:   constant.ErrCodeOk,
		Errmsg: constant.ErrMsgOk,
	}
	return response, nil
}
