package cgi

import (
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"glogin/config"
	"time"
)

// AgreementCheckReq 用户协议检查请求
type AgreementCheckReq struct {
	BundleID string `json:"bundle_id" binding:"required"`
	Version  int32  `json:"version"`
	Language string `json:"language" binding:"required"`
}

// AgreementCheckRsp 用户协议检查返回
type AgreementCheckRsp struct {
	Version int32  `json:"version"`
	Tips    string `json:"tips"`
	PopUp   bool   `json:"pop_up"`
}

func CheckHandler(ctx *gin.Context) {
	checkReq := &AgreementCheckReq{}
	reqId := uuid.New()
	before := time.Now().UnixNano()
	checkRsp := AgreementCheckRsp{}
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("got panic", "err", err)
		}
		log.Infow("CheckHandler response", "reqId", reqId, "checkRsp", checkRsp, "time_cost", (time.Now().UnixNano()-before)/1000000)
	}()
	err := ctx.Bind(checkReq)
	log.Infow("CheckHandler request", "reqId", reqId, "request", checkReq)
	if err != nil {
		ParseRequestError(ctx, 500, err)
		return
	}
	//todo 验证逻辑
	agreementMap := config.PackageParamRst(checkReq.BundleID, "user_agreement").Map()
	if len(agreementMap) > 0 {
		v, ok := agreementMap[checkReq.Language]
		if ok {
			CurVsn := v.Get("vsn").Int()
			if int64(checkReq.Version) < CurVsn {
				checkRsp.PopUp = true
				checkRsp.Tips = v.Get("tips").String()
			}
			checkRsp.Version = int32(CurVsn)
		}
	}
	ctx.JSON(200, checkRsp)
}
