package cgi

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ParseRequestError(c *gin.Context, code int32, err error) {
	c.JSON(500, gin.H{
		"code":   code,
		"errmsg": fmt.Sprintf("parse json err:%v", err),
	})
}

//func antiConvertClient(antiRsp *pb_obsession.CheckStateQueryResponse) *glogin.StateQueryResponse {
//	rspToClient := &glogin.StateQueryResponse{}
//	defer func() {
//		if err := recover(); err != nil {
//			log.Errorw("got panic", "err", err)
//		}
//	}()
//	rspBytes, err := protojson.Marshal(antiRsp)
//	if err != nil {
//		return rspToClient
//	}
//	err = protojson.Unmarshal(rspBytes, rspToClient)
//	return rspToClient

	//return &glogin.StateQueryResponse{
	//	RequestId:            antiRsp.RequestId,
	//	ErrCode:              antiRsp.ErrCode,
	//	ErrMsg:               antiRsp.ErrMsg,
	//	AuthenticationStatus: antiRsp.AuthenticationStatus,
	//	IsHoliday:            antiRsp.IsHoliday,
	//	LeftGameTime:         antiRsp.LeftGameTime,
	//	EachPayAmount:        antiRsp.EachPayAmount,
	//	LeftPayAmount:        antiRsp.LeftPayAmount,
	//	LoginCode:            antiRsp.LoginCode,
	//	LoginMessage:         antiRsp.LoginMessage,
	//}
//}
