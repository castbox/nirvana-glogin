package cgi

import (
	"glogin/constant"
	glogin2 "glogin/pbs/glogin"
)

type Bind struct {
}

func (Bind) BindThird(request *glogin2.VistorBindThridReq) (response *glogin2.VistorBindThridRsp, err error) {
	response = &glogin2.VistorBindThridRsp{
		Code:   constant.ErrMsgOk,
		Errmsg: constant.ErrMsgOk,
	}
	return response, nil
}
