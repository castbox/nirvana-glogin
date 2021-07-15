package cgi

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ReplaySuccess(c *gin.Context, accountId int32, platString string) {
	c.JSON(200, gin.H{
		"code":       0,
		"dh_token":   accountId,
		"third_plat": platString,
	})
}

func ParseRequestError(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"code":   500,
		"errmsg": fmt.Sprintf("parse json err:%v", err),
	})
}
