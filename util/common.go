package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/dgrijalva/jwt-go"
	"glogin/config"
	"io"
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
	TokenExpiredTime = 30 * 24 * 60 * 60
	//TokenExpiredTime = 60 * 10
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

// 输入文本，key，返回经过base64(aes-128)加密的结果，并且会将nonce放入到前12个字节
func Enr(plaintext, key string) string {
	keyBytes, _ := hex.DecodeString(key)
	block, _ := aes.NewCipher(keyBytes)

	nonce := make([]byte, 12)
	io.ReadFull(rand.Reader, nonce)

	gcm, _ := cipher.NewGCM(block)
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	result := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(result)
}
