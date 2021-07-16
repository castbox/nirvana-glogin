package cgi

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ParseRequestError(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"code":   500,
		"errmsg": fmt.Sprintf("parse json err:%v", err),
	})
}
