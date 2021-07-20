package cgi

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"glogin/util"
)

type TokenReq struct {
	Uid   int32  `json:"uid" binding:"required"`
	Token string `json:"token"`
}

func TokenHandler(ctx *gin.Context) {
	verify := &TokenReq{}
	err := ctx.Bind(verify)
	if err != nil {
		ParseRequestError(ctx, err)
		return
	}
	token, parseErr := jwt.ParseWithClaims(verify.Token, jwt.MapClaims{}, util.KeyMethod)
	if parseErr != nil {
		err = fmt.Errorf("token parese err:%v", parseErr)
		ParseRequestError(ctx, err)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		accountId := int32(claims["dhAccountId"].(float64))
		//expired := int64(claims["expire"].(float64))
		if accountId == verify.Uid {
			ctx.JSON(200, gin.H{
				"errno": 0})
			return
		}
	} else {
		ctx.JSON(200, gin.H{
			"errno": 1})
		return
	}
	ctx.JSON(200, gin.H{
		"errno": 1})
}
