package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	log "github.com/castbox/nirvana-gcore/glog"
	"github.com/dgrijalva/jwt-go"
	"glogin/config"
	"io"
	"regexp"
	"strings"
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
		log.Warnw("ValidDHToken error", "DHToken", err)
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
		log.Warnw("ValidDHToken error", "DHToken", err)
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

// 匹配 手机号,邮箱,中文,身份证等等 进行脱敏处理
func HideStar(str string) (result string) {
	if str == "" {
		return "***"
	}
	if strings.Contains(str, "@") {
		// 邮箱
		res := strings.Split(str, "@")
		if len(res[0]) < 3 {
			resString := "***"
			result = resString + "@" + res[1]
		} else {
			res2 := Substr2(str, 0, 3)
			resString := res2 + "***"
			result = resString + "@" + res[1]
		}
		return result
	} else {
		reg := `^1[0-9]\d{9}$`
		rgx := regexp.MustCompile(reg)
		mobileMatch := rgx.MatchString(str)
		if mobileMatch {
			// 手机号
			result = Substr2(str, 0, 3) + "****" + Substr2(str, 7, 11)
		} else {
			nameRune := []rune(str)
			lens := len(nameRune)
			if lens <= 1 {
				result = "***"
			} else if lens == 2 {
				result = string(nameRune[:1]) + "*"
			} else if lens == 3 {
				result = string(nameRune[:1]) + "*" + string(nameRune[2:3])
			} else if lens == 4 {
				result = string(nameRune[:1]) + "**" + string(nameRune[lens-1:lens])
			} else if lens > 4 {
				result = string(nameRune[:2]) + "***" + string(nameRune[lens-2:lens])
			}
		}
		return
	}
}

func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	return string(rs[start:end])
}
