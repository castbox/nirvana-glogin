package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	apiUrl     = "https://ye.dun.163yun.com/v1/oneclick/check" //本机认证服务身份证实人认证在线检测接口地址
	version    = "v1"
	secretId   = "a20a4fd6a0ac8a32a2b8d01042433778" //产品密钥ID，产品标识
	secretKey  = "945b23b071ae712e21e1722bc967b753" //产品私有密钥，服务端生成签名信息使用，请严格保管，避免泄露
	businessId = "efedd541fba94b82a9854363975f16e0" //业务ID，易盾根据产品业务特点分配
)

//请求易盾接口
func check(params url.Values) *simplejson.Json {
	params["secretId"] = []string{secretId}
	params["businessId"] = []string{businessId}
	params["version"] = []string{version}
	params["timestamp"] = []string{strconv.FormatInt(time.Now().UnixNano()/1000000, 10)}
	params["nonce"] = []string{string(make([]byte, 32))} //32位随机字符串
	params["signature"] = []string{gen_signature(params)}

	resp, err := http.Post(apiUrl, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))

	if err != nil {
		fmt.Println("调用API接口失败:", err)
		return nil
	}

	defer resp.Body.Close()

	contents, _ := ioutil.ReadAll(resp.Body)
	result, _ := simplejson.NewJson(contents)
	return result
}

//生成签名信息
func gen_signature(params url.Values) string {
	var paramStr string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		paramStr += key + params[key][0]
	}
	paramStr += secretKey
	md5Reader := md5.New()
	md5Reader.Write([]byte(paramStr))
	return hex.EncodeToString(md5Reader.Sum(nil))
}

func main() {
	params := url.Values{
		//token为易盾返回的token
		"token": []string{"f701918e05124ce6aa90f76cfff5f534"},
		//accessToken为运营商预取号获取到的token
		"accessToken": []string{"7b22616363657373546f6b656e223a226e6d6633373430303638636230373465333539303664366336356461363837333366222c22677741757468223a2231353836227d"},
	}
	ret := check(params)

	code, _ := ret.Get("code").Int()

	if code == 200 {
		data, _ := ret.Get("data").Map()
		phone, _ := data["phone"].(string)
		if len(phone) != 0 {
			fmt.Printf("取号成功, 执行登录等流程!")
		} else {
			resultCode, _ := data["resultCode"].(string)
			fmt.Printf("取号失败,建议进行二次验证,例如短信验证码。运营商返回码resultCode为: %s", resultCode)
		}
	} else {
		fmt.Printf("降级走短信！")
	}
}
