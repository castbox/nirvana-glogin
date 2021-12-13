package cgi

import (
	"fmt"
	"git.dhgames.cn/svr_comm/anti_obsession/pbs/pb_obsession"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/gin-gonic/gin"
	"glogin/constant"
	"glogin/internal"
	"glogin/internal/anti"
	"glogin/pbs/glogin"
	"glogin/util"
	"time"
)

const (
	AuthKey = "0b0ab8e6f04190958e90dfbaef180f7e"
)

// SDK实名信息认证
type AutiCheckRequest struct {
	GameCd  string `json:"game_cd" binding:"required"` //游戏唯一表示
	Account string `json:"account"`                    //卓杭账号ID
	Pid     string `json:"pid" binding:"required"`     //身份PID
	Name    string `json:"name"`                       //名字
}

// 防沉迷返回
type StateQueryResponse struct {
	RequestId            string `json:"request_id"`            // 每次请求唯一标识
	ErrCode              string `json:"err_code"`              // 返回代码
	ErrMsg               string `json:"err_msg"`               // 错误信息
	AuthenticationStatus int32  `json:"authentication_status"` // 实名状态：0：成年，1：游客，2：0-8岁，3：8-16岁，4：16-18岁
	IsHoliday            bool   `json:"is_holiday"`            // 是否节假日
	LeftGameTime         int32  ` json:"left_game_time,"`      // 剩余游戏时间，已成年请忽略
	EachPayAmount        int32  `json:"each_pay_amount"`       // 单次可充值额度，已成年请忽略
	LeftPayAmount        int32  `json:"left_pay_amount"`       // 总剩充值额度，已成年请忽略
	LoginCode            int32  `json:"login_code"`            // 是否可以登陆的提示code
	LoginMessage         string `json:"login_message"`         // 如果login_code不为0，相应的提示字段。
}

// 实名信息认证返回
type AutiCheckResponse struct {
	ErrCode        string                     `json:"err_code"`       // 每次请求唯一标识,用作问题校验时查询
	ErrMsg         string                     `json:"err_msg"`        // 错误码
	CheckMsg       string                     `json:"check_msg"`      // 实名提示信息,如果认证成功，并且玩家是未成年，客户端需要显示此字段
	Authentication *glogin.StateQueryResponse `json:"authentication"` // 防沉迷状态
}

func AutiHandler(ctx *gin.Context) {
	checkReq := &AutiCheckRequest{}
	err := ctx.Bind(checkReq)
	log.Infow("new query AutiHandler checkReq", "request", checkReq)

	// 实名认证
	checkRsp := &AutiCheckResponse{
		ErrCode:        "0",
		Authentication: &glogin.StateQueryResponse{},
	}

	if err != nil {
		ParseRequestError(ctx, 500, err)
		return
	}
	plaintext := checkReq.Name + "@" + checkReq.Pid
	playerInfo := util.Enr(plaintext, AuthKey)
	checkIn := &pb_obsession.CheckRequest{
		GameCd:     checkReq.GameCd,  //game_cd
		Id:         checkReq.Account, //账号ID
		PlayerInfo: playerInfo,
	}

	before := time.Now().UnixNano()
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("got panic", "err", err)
		}
		log.Infow("new query AutiHandler rsp1", "request", checkReq, "rsp", checkRsp, "time_cost", (time.Now().UnixNano()-before)/1000000)
	}()

	log.Infow("anti.Check req", "req", checkIn)
	antiCheckRsp, errCheck := anti.Check(checkIn)
	if errCheck != nil {
		checkRsp.ErrCode = constant.ErrCodeStrAutiRpc
		checkRsp.ErrMsg = fmt.Sprintf(" AutiCheck error: %v", errCheck)
		log.Infow("anti.Check rsp 1", "rsp", antiCheckRsp)
		ctx.JSON(500, checkRsp)
		return
	}
	checkRsp.ErrMsg = antiCheckRsp.ErrMsg
	checkRsp.ErrCode = antiCheckRsp.ErrCode
	checkRsp.CheckMsg = antiCheckRsp.CheckMsg
	if antiCheckRsp.ErrCode != constant.ErrCodeStrOk {
		log.Infow("anti.Check rsp 2", "rsp", antiCheckRsp)
		ctx.JSON(200, checkRsp)
		return
	}
	log.Infow("anti.Check rsp", "rsp", antiCheckRsp)
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
	r2, ok := antiQueryRsp.(*pb_obsession.CheckStateQueryResponse)
	if ok {
		checkRsp.Authentication = antiConvertClient(r2)
		//checkRsp.Authentication = &StateQueryResponse{
		//	RequestId:            r2.RequestId,
		//	ErrCode:              r2.ErrCode,
		//	ErrMsg:               r2.ErrMsg,
		//	AuthenticationStatus: r2.AuthenticationStatus,
		//	IsHoliday:            r2.IsHoliday,
		//	LeftGameTime:         r2.LeftGameTime,
		//	EachPayAmount:        r2.EachPayAmount,
		//	LeftPayAmount:        r2.LeftPayAmount,
		//	LoginCode:            r2.LoginCode,
		//	LoginMessage:         r2.LoginMessage,
		//}
	}
	log.Infow("new query AutiHandler rsp", "rsp", checkRsp, "request", checkReq)
	ctx.JSON(200, checkRsp)
}
