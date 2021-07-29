package cgi

import (
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/gin-gonic/gin"
	"glogin/constant"
	"glogin/internal"
	"glogin/internal/anti"
	anti_authentication "glogin/pbs/authentication"
	"glogin/pbs/glogin"
	"glogin/util"
)

const (
	AuthKey = "0b0ab8e6f04190958e90dfbaef180f7e"
)

// SDK实名信息认证
type AutiCheckRequest struct {
	GameCd  string `json:"game_cd" binding:"required"`
	Account string `json:"account"`
	Pid     string `json:"pid" binding:"required"`
	Name    string `json:"name"`
}

// 实名信息认证返回
type AutiCheckResponse struct {
	ErrCode        string                     `json:"err_code"`
	ErrMsg         string                     `json:"err_msg"`
	Authentication *glogin.StateQueryResponse `json:"authentication"`
}

func AutiHandler(ctx *gin.Context) {
	checkReq := &AutiCheckRequest{}
	err := ctx.Bind(checkReq)
	log.Infow("new query AutiHandler checkReq", "request", checkReq)
	if err != nil {
		ParseRequestError(ctx, 500, err)
		return
	}
	plaintext := checkReq.Name + "@" + checkReq.Pid
	playerInfo := util.Enr(plaintext, AuthKey)
	checkIn := &anti_authentication.CheckRequest{
		GameCd:     checkReq.GameCd,  //game_cd
		Id:         checkReq.Account, //账号ID
		PlayerInfo: playerInfo,
	}
	// 实名认证
	checkRsp := &AutiCheckResponse{
		ErrCode:        "0",
		Authentication: &glogin.StateQueryResponse{},
	}
	antiCheckRsp, errCheck := anti.Check(checkIn)
	if errCheck != nil {
		checkRsp.ErrCode = constant.ErrCodeStrAutiRpc
		checkRsp.ErrMsg = fmt.Sprintf(" AutiCheck error: %v", checkRsp)
		ctx.JSON(500, checkRsp)
		return
	}
	checkRsp.ErrMsg = antiCheckRsp.ErrMsg
	checkRsp.ErrCode = antiCheckRsp.ErrCode
	if antiCheckRsp.ErrCode != constant.ErrCodeStrOk {
		ctx.JSON(200, checkRsp)
		return
	}
	//  查询实名信息
	reqState := internal.Req{
		Account: checkReq.Account,
		GameCd:  checkReq.GameCd,
	}
	log.Infow("anti.StateQuery req", "req", reqState)
	antiQueryRsp, antiErr := anti.StateQuery(reqState)
	if antiErr != nil {
		ctx.JSON(200, checkRsp)
		return
	}
	log.Infow("anti.StateQuery rsp", "rsp", antiQueryRsp)
	r2, ok := antiQueryRsp.(*anti_authentication.StateQueryResponse)
	if ok {
		checkRsp.Authentication = (*glogin.StateQueryResponse)(r2)
	}
	ctx.JSON(200, checkRsp)
}
