package plat

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	_ "github.com/gogf/gf/encoding/gjson"
	"github.com/tidwall/gjson"
	"glogin/config"
	"glogin/pbs/glogin"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	fullUrl    = "https://ye.dun.163yun.com/v1/oneclick/check" //本机认证服务身份证实人认证在线检测接口地址
	version    = "v1"
	secretId   = "your_secret_id"   //产品密钥ID，产品标识
	secretKey  = "your_secret_key"  //产品私有密钥，服务端生成签名信息使用，请严格保管，避免泄露
	businessId = "your_business_id" //业务ID，易盾根据产品业务特点分配
)

var YeDun yedun

type yedun struct{}

// Auth 登录返回第三方账号id 和 错误信息
func (y yedun) Auth(request *glogin.ThirdLoginReq) (string, string, error) {
	log.Infow("yedun_check auth", "request", request)
	fullUrl = config.Field("yedun_oauth_url").String()
	secretId = config.Field("yedun_secret_id").String()
	secretKey = config.Field("yedun_secret_key").String()
	businessId = config.Field("yedun_businessId").String()
	params := url.Values{
		//token为易盾返回的token
		"token": []string{request.ThirdToken},
		//accessToken为运营商预取号获取到的token
		"accessToken": []string{request.AccessToken},
	}
	params["secretId"] = []string{secretId}
	params["businessId"] = []string{businessId}
	params["version"] = []string{version}
	params["timestamp"] = []string{strconv.FormatInt(time.Now().UnixNano()/1000000, 10)}
	params["nonce"] = []string{string(make([]byte, 32))} //32位随机字符串
	params["signature"] = []string{genSignature(params)}
	resp, err := http.Post(fullUrl, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		log.Errorw("yedun auth error ", "err", err)
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf(resp.Status)
	}
	contents, err := ioutil.ReadAll(resp.Body)
	//{"code":450,"msg":"wrong token"}
	if err != nil {
		resErr := fmt.Errorf("failed reading from metadata server: %w", err)
		log.Errorw("yedun auth error ", "resErr", resErr)
		return "", "", resErr
	}
	code := gjson.GetBytes(contents, "code").Int()
	if code == 200 {
		phone := gjson.GetBytes(contents, "phone").String()
		if len(phone) != 0 {
			log.Infow("yedun auth get phonenum", "phone", phone)
			return phone, phone, nil
		} else {
			resultCode := gjson.GetBytes(contents, "resultCode").String()
			log.Errorw("yedun auth get phonenum error", "resultCode", resultCode)
			resultErr := fmt.Errorf("yedun auth get phonenum error: %v", resultCode)
			return "", "", resultErr
		}
	} else {
		msg := gjson.GetBytes(contents, "msg").String()
		resultErr := fmt.Errorf("yedun auth get phonenum error code: %v msg: %s", code, msg)
		return "", "", resultErr
	}

}

//生成签名信息
func genSignature(params url.Values) string {
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

func (y yedun) String() string {
	return "yedun"
}
