package smfpcrypto

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	//	"fmt"
	//	"os"
	"errors"
	"strings"
)

func IsBoxId(data string) bool {
	if len(data) == 89 && data[0] == 'B' {
		return true
	}
	return false
}

func rsaDecrypt(ciphertext []byte, privateKey []byte) ([]byte, error) {
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(privateKey)
	if err != nil {
		return nil, errors.New("x509 parse ERROR")
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func base64Decode(src string) ([]byte, error) {
	by, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return by, err
	}
	return by, nil
}

func getMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func getCheckSerial(smid string) string {
	trueId := smid[0:48]
	var platform string
	platformTag := smid[48]
	switch platformTag {
	case '0':
		platform = "android"
		break
	case '1':
		platform = "ios"
		break
	case '2':
		platform = "web"
		break
	case '3':
		platform = "webapp"
		break
	case '4':
		platform = "quickapp"
		break
	default:
		platform = "unknown"
		break
	}

	spellbound := "shumei_" + platform + "_sec_key_" + trueId
	hash := getMd5String(spellbound)[0:14]
	return trueId + strings.ToLower(hash)
}

func ParseSMID(srcId string) string {
	var des string
	switch srcId[0] {
	case 'B':
		prikey := "QbDRypciANq6gOected1"
		des, _ = ParseBoxId(srcId, prikey)
		break
	case 'D':
		accessKey := "QbDRypciANq6gOected1"
		des = ParseBoxData(srcId, accessKey)
		break
	default:
		des = srcId
		break
	}
	return des
}

func ParseBoxData(boxId string, accessKey string) string {
	return boxId
}

func ParseBoxId(boxdata string, prikey string) (string, error) {
	if !IsBoxId(boxdata) {
		return "", errors.New("boxId is illegal")
	}

	mdata := boxdata[1:]

	data, err := base64Decode(mdata)
	if err != nil {
		return "", errors.New("data base64 decode error")
	}

	pri, err := base64Decode(prikey)
	if err != nil {
		return "", errors.New("prikey base64 decode error")
	}

	plainText, err := rsaDecrypt(data, pri)
	if err != nil {
		return "", errors.New("not boxId")
	}

	return getCheckSerial(string(plainText)), nil
}
