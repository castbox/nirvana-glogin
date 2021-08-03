package smfpcrypto

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"glogin/config"
	"glogin/internal/xhttp"
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

// 数美ID处理相关函数 -去除时间戳后缀
func DealSMID(srcId string) string {
	log.Infow("dealSMID", "src", srcId)
	lastPos := strings.LastIndex(srcId, "-")
	if lastPos == -1 {
		return srcId
	}
	desId := srcId[0:lastPos]
	log.Infow("dealSMID", "des", desId)
	return desId
}

func ParseSMID(smId string) string {
	log.Infow("ParseSMID", "src", smId)
	srcId := DealSMID(smId)
	var des string
	switch srcId[0] {
	case '2':
		des = srcId
		break
	case 'B':
		prikey := config.Field("parse_smid_private_key").String()
		if prikey == "" {
			prikey = "-----BEGIN RSA PRIVATE KEY-----\nMIIBPAIBAAJBAOfPLQ993UR8qJoCVJsj00/BcPDbKIjEDYnqMjgUAiQkMgYf9O4L7+WNhhtA+kIllsHpEAYJuSdl04wP05Pk0TkCAwEAAQJBAMzdJOafBrDjNqI9UwZ0x+ihfa3vEcik844iItW6oRXMFIo+P+6YHjgiiyLXeSu+60WQ4IfWdRZNdbHMhr1IIN0CIQD6GzKUls0YXxASUmdcSTUFeqXcedkhcLafHTk8jqcX8wIhAO1FmW+cyx1gm4msyhgXN1Fb7frFHniaP5L89zc6NwkjAiEAmA0e5A0GJVHt8GWepxFupaUZ3v9JDTZ8ICHhITrMxRcCIFI92Z0yP8UDA2aJGdOX2Hi+4JIXWSR8cqTEQfxGlWT5AiEA5TqIC6znNIGzeAeuz3Hdj4srmAEP/VG9EkDvdgMT6Tg=\n-----END RSA PRIVATE KEY-----"
		}
		parseDes, err := ParseBoxId(srcId, prikey)
		if err != nil {
			des = srcId
		} else {
			des = parseDes
		}
		break
	case 'D':
		accessKey := config.Field("parse_smid_access_key").String()
		if accessKey == "" {
			accessKey = "QbDRypciANq6gOected1"
		}
		des = ParseBoxData(srcId, accessKey)
		break
	default:
		des = srcId
		break
	}
	log.Infow("ParseSMID", "des", des)
	return des
}

type DeviceLabels struct {
	Id                     string      `json:"id"`
	FakeDevice             interface{} `json:"fake_device"`
	DeviceSuspiciousLabels interface{} `json:"device_suspicious_labels"`
	MonkeyDevice           interface{} `json:"monkey_device"`
}

type ParseRsp struct {
	Code         int          `json:"code"`
	Message      string       `json:"message"`
	RequestId    string       `json:"requestId"`
	DeviceLabels DeviceLabels `json:"deviceLabels"`
}

func ParseBoxData(srcBoxId string, accessKey string) string {
	if accessKey == "" {
		accessKey = "QbDRypciANq6gOected1"
	}
	FullUrl := config.Field("parse_smid_url").String()
	if FullUrl == "" {
		FullUrl = "http://api-tianxiang-bj.fengkongcloud.com/tianxiang/v4"
	}
	httpClient := xhttp.NewClient().Type(xhttp.TypeJSON)
	req := xhttp.BodyMap{}
	req.Set("accessKey", accessKey)
	data := xhttp.BodyMap{}
	data.Set("deviceId", srcBoxId)
	data.Set("tokenId", "")
	req.Set("data", data)
	res, bs, errs := httpClient.Post(FullUrl).SendBodyMap(req).EndBytes()
	if len(errs) > 0 {
		log.Errorw("ParseSMID ParseBoxData HTTP Request Error1", "errs", errs[0])
		return srcBoxId
	}
	if res.StatusCode != 200 {
		log.Errorw("ParseSMID ParseBoxData HTTP Request Error2,", "StatusCode", res.StatusCode)
		return srcBoxId
	}
	log.Infow("ParseSMID ParseBoxData HTTP Rsp,", "string(bs)", string(bs))
	smRsp := new(ParseRsp)
	if err := json.Unmarshal(bs, smRsp); err != nil {
		log.Infow("ParseSMID ParseBoxData HTTP Request Error3,", "StatusCode", res.StatusCode)
		return srcBoxId
	}
	if smRsp.Code == 1100 {
		return smRsp.DeviceLabels.Id
	} else {
		strErr := fmt.Errorf("ParseSMID Rsp%s", string(bs))
		log.Infow("ParseSMID ParseBoxData HTTP Request Error4,", "strErr", strErr)
		return srcBoxId
	}
	return smRsp.Message
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
