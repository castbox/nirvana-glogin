package cgi

import (
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"git.dhgames.cn/svr_comm/gmoss/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	glogin2 "glogin/pbs/glogin"
	"reflect"
)

func ServiceHandler(ctx *gin.Context) {
	service := ctx.Param("service")
	action := ctx.Param("action")
	log.Infow("new http request", "service", service)
	response := schedule(service, action, ctx)
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	m := jsonpb.Marshaler{EmitDefaults: true, OrigName: true}
	if err := m.Marshal(ctx.Writer, response); err != nil {
		log.Warnw("ServiceHandler http返回结果，json解析报错", "response", response)
	}
}

func schedule(service string, action string, ctx *gin.Context) (result proto.Message) {
	request := reflect.ValueOf(routes[service][action][1]).Call(nil)
	req := request[0].Interface()
	if err := ctx.BindJSON(req); err != nil {
		gmoss.Error("%v", err)
		panic(err)
	}
	log.Infow("schedule 2", "req", req)
	in := []reflect.Value{reflect.ValueOf(req)}
	bCtx := routes[service][action][2]
	log.Infow("routes[service][action][2] 2", "bAddCtx", bCtx)
	if bCtx == true {
		in = append(in, reflect.ValueOf(ctx))
	}
	f := reflect.ValueOf(routes[service][action][0])
	result = f.Call(in)[0].Interface().(proto.Message)
	return
}

var loginFuncs = map[string][]interface{}{
	"third": {
		(&Login{}).ThirdEx,
		glogin2.NewThirdLoginReq,
		true,
	},
	"sms": {
		(&Login{}).SMSEx,
		glogin2.NewSmsLoginReq,
		true,
	},
	"visitor": {
		(&Login{}).VisitorEx,
		glogin2.NewVisitorLoginReq,
		true,
	},
	"fast": {
		(&Login{}).FastEx,
		glogin2.NewFastLoginReq,
		true,
	},
}

var bindFuncs = map[string][]interface{}{
	"third": {
		(&Bind{}).BindThird,
		glogin2.NewVistorBindThridReq,
		false,
	},
}

var antiFuncs = map[string][]interface{}{
	"antiFuncs": {
		(&Bind{}).BindThird,
		glogin2.NewVistorBindThridReq,
		false,
	},
}

var routes = map[string]map[string][]interface{}{
	"login": loginFuncs,
	"bind":  bindFuncs,
	"anti":  antiFuncs,
}
