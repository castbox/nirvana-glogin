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
