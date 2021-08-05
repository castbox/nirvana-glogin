package plat

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/bitly/go-simplejson"
	_ "github.com/gogf/gf/encoding/gjson"
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

const (
	//apiUrl     = "https://ye.dun.163yun.com/v1/oneclick/check" //本机认证服务身份证实人认证在线检测接口地址
	version = "v1"
	//secretId   = "a20a4fd6a0ac8a32a2b8d01042433778" //产品密钥ID，产品标识
	//secretKey  = "945b23b071ae712e21e1722bc967b753" //产品私有密钥，服务端生成签名信息使用，请严格保管，避免泄露
	//businessId = "efedd541fba94b82a9854363975f16e0" //业务ID，易盾根据产品业务特点分配
)

var YeDun yedun

type yedun struct{}

// Auth 登录返回第三方账号tokenId openId 错误信息
func (y yedun) Auth(request *glogin.ThirdLoginReq) (*AuthRsp, error) {
	log.Infow("yedun_check auth", "request", request)
	apiUrl := config.Field("yedun_oauth_url").String()
	//yedunParam := config.PackageParam(request.Game.BundleId, "yedun_param")
	//if yedunParam == "" {
	//	resErr := fmt.Errorf("failed reading from metadata server: %s", request.Game.BundleId)
	//	return nil, resErr
	//}
	//paramArr := strings.Split(yedunParam, "|")
	secretId := config.PackageParam(request.Game.BundleId, "yedun_secret_id")
	secretKey := config.PackageParam(request.Game.BundleId, "yedun_secret_key")
	businessId := config.PackageParam(request.Game.BundleId, "yedun_business_id")
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
	params["signature"] = []string{genSignature(params, secretKey)}

	log.Infow("yedun_check auth url info", "apiUrl", apiUrl, "body", params)
	resp, err := http.Post(apiUrl, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		log.Errorw("yedun auth error ", "err", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}
	contents, err := ioutil.ReadAll(resp.Body)
	//{"code":450,"msg":"wrong token"}
	//{"code":200,"data":{"phone":"19181732997","resultCode":"0"},"msg":"ok"}}
	if err != nil {
		resErr := fmt.Errorf("failed reading from metadata server: %w", err)
		log.Errorw("yedun auth error ", "resErr", resErr)
		return nil, resErr
	}
	result, _ := simplejson.NewJson(contents)
	log.Infow("yedun auth rsp", "contents", contents, "result", result)
	code, _ := result.Get("code").Int()
	if code == 200 {
		data, _ := result.Get("data").Map()
		phone, _ := data["phone"].(string)
		if len(phone) != 0 {
			//fmt.Printf("取号成功, 执行登录等流程!")
			log.Infow("yedun auth get phonecode", "phone", phone)
			return &AuthRsp{
				Uid:     phone,
				UnionId: phone,
			}, nil
		} else {
			resultCode, _ := data["resultCode"].(string)
			log.Errorw("yedun auth get phonenum error", "resultCode", resultCode)
			resultErr := fmt.Errorf("yedun auth get phonenum error: %v", resultCode)
			return nil, resultErr
			//fmt.Printf("取号失败,建议进行二次验证,例如短信验证码。运营商返回码resultCode为: %s", resultCode)
		}
	} else {
		//fmt.Printf("降级走短信！")
		msg, _ := result.Get("msg").String()
		resultErr := fmt.Errorf("yedun auth get phonenum error code: %v msg: %s", code, msg)
		return nil, resultErr
	}

}

//生成签名信息
func genSignature(params url.Values, secretKey string) string {
	var paramStr string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		paramStr += key + params[key][0]
	}
	//secretKey := config.Field("yedun_secret_key").String()
	//secretKey := config.PackageParam(bundleId, "yedun_secret_key")
	paramStr += secretKey
	md5Reader := md5.New()
	md5Reader.Write([]byte(paramStr))
	return hex.EncodeToString(md5Reader.Sum(nil))
}

func (y yedun) String() string {
	return "yedun"
}

func (y yedun) DbFieldName() string {
	return "phone"
}
