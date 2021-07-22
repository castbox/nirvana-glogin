package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/dgrijalva/jwt-go"
	"glogin/config"
	"time"
)

const (
	TimeLayout   = "2006-01-02 15:04:05"
	TimeLayout_2 = "20060102150405"
	DateLayout   = "2006-01-02"
	NULL         = ""
)

type File struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

const (
	//TokenExpiredTime = 30 * 24 * 60 * 60
	TokenExpiredTime = 60 * 10
)

func Md5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

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

var (
	ExpiredToken = fmt.Errorf("token is expired")
	InvalidToken = fmt.Errorf("token parese is not valid")
)

// ValidDHToken 验证token是否有效
func ValidDHToken(tokenString string) (accountId int32, err error) {
	token, parseErr := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, KeyMethod)
	if parseErr != nil {
		err = fmt.Errorf("token parese err:%v", parseErr)
		log.Errorw("ValidDHToken error", "DHToken", err)
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		accountId = int32(claims["dhAccountId"].(float64))
		expired := int64(claims["expire"].(float64))
		if expired < time.Now().Unix() {
			return 0, ExpiredToken
		} else {
			return accountId, nil
		}
	} else {
		err = InvalidToken
		log.Errorw("ValidDHToken error", "DHToken", err)
		return 0, err
	}
}
