package cgi

import (
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	glogin2 "glogin/pbs/glogin"
	"reflect"
)

func ServiceHandler(ctx *gin.Context) {
	service := ctx.Param("service")
	action := ctx.Param("action")
	response := schedule(service, action, ctx)
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	m := jsonpb.Marshaler{EmitDefaults: true, OrigName: true}
	if err := m.Marshal(ctx.Writer, response); err != nil {
		log.Warnw("ServiceHandler http返回结果，json解析报错", "response", response)
	}
}

func schedule(service string, action string, ctx *gin.Context) (result proto.Message) {
	request := reflect.ValueOf(routes[service][action][1]).Call(nil)
	log.Infow("schedule 2", "req", request)
	if err := ctx.BindJSON(request[0].Interface()); err != nil {
		panic(err)
	}
	function := reflect.ValueOf(routes[service][action][0])
	result = function.Call(request)[0].Interface().(proto.Message)
	return
}

var loginFuncs = map[string][]interface{}{
	"third": {
		(&Login{}).Third,
		glogin2.NewThirdLoginReq,
	},
	"sms": {
		(&Login{}).SMS,
		glogin2.NewSmsLoginReq,
	},
	"visitor": {
		(&Login{}).Visitor,
		glogin2.NewVisitorLoginReq,
	},
	"fast": {
		(&Login{}).Fast,
		glogin2.NewFastLoginReq,
	},
}

var bindFuncs = map[string][]interface{}{
	"third": {
		(&Bind{}).BindThird,
		glogin2.NewVistorBindThridReq,
	},
}

var antiFuncs = map[string][]interface{}{
	"antiFuncs": {
		(&Bind{}).BindThird,
		glogin2.NewVistorBindThridReq,
	},
}

var routes = map[string]map[string][]interface{}{
	"login": loginFuncs,
	"bind":  bindFuncs,
	"anti":  antiFuncs,
}
