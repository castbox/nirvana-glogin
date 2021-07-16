package utils

import (
	"encoding/base64"
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/dgrijalva/jwt-go"
	"glogin/config"
	"time"
)

const (
	TokenExpiredTime = 30 * 24 * 60 * 60
)

func Base64Decode(src string) ([]byte, error) {
	by, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return by, err
	}
	return by, nil
}

func KeyMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("UnExpected signing method: %v", token.Header["alg"])
	}
	return []byte(config.Field("jwt_secret").String()), nil
}

// GenDHToken 根据账号id 生成token信息
func GenDHToken(accountId int32) string {
	expired := time.Now().Unix() + TokenExpiredTime
	// 将account、过期时间作为数据写入 token 中
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"dhAccountId": accountId,
		"expire":      expired,
	})
	// SecretKey 用于对用户数据进行签名
	res, err := token.SignedString([]byte(config.Field("jwt_secret").String()))
	if err != nil {
		log.Errorw("GenDHToken SignedString", "GenDHToken", err)
	}
	return res
}
