package plat

import (
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"glogin/config"
	"glogin/internal/xhttp"
	"glogin/pbs/glogin"
)

// 获取开放平台，access_token 返回结构体
type Oauth2AccessToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Openid       string `json:"openid,omitempty"`
	Scope        string `json:"scope,omitempty"`
	Unionid      string `json:"unionid,omitempty"`
	Errcode      int    `json:"errcode,omitempty"` // 错误码
	Errmsg       string `json:"errmsg,omitempty"`  // 错误信息
}

var WeChat wechat

type wechat struct{}

// Auth 登录返回第三方账号tokenId openId 错误信息
func (w wechat) Auth(request *glogin.ThirdLoginReq) (string, string, error) {
	log.Infow("wechat Auth", "request", request)
	appId := config.PackageParam(request.Game.BundleId, "wx_app_id")
	appKey := config.PackageParam(request.Game.BundleId, "wx_app_key")
	code := request.AccessToken
	accessToken, err := GetOauth2AccessToken(appId, appKey, code)
	if err != nil {
		log.Errorw("wechat auth error ", "err", err)
		return "", "", err
	}
	if accessToken.Unionid == "" {
		log.Infow("wechat Auth Unionid is nil", "Openid", accessToken.Openid)
		accessToken.Unionid = accessToken.Openid
	}
	return accessToken.Openid, accessToken.Unionid, nil
}

func (w wechat) String() string {
	return "wechat"
}

func (w wechat) DbFieldName() string {
	return "we_chat"
}

//  GetOauth2AccessToken 微信第三方登录，code 换取 access_token
//	appId：应用唯一标识，在微信开放平台提交应用审核通过后获得
//	appSecret：应用密钥AppSecret，在微信开放平台提交应用审核通过后获得
//	code：App用户换取access_token的code
//	文档：https://developers.weixin.qq.com/doc/oplatform/Mobile_App/WeChat_Login/Development_Guide.html
func GetOauth2AccessToken(appId, appSecret, code string) (accessToken *Oauth2AccessToken, err error) {
	accessToken = new(Oauth2AccessToken)
	url := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + appId + "&secret=" + appSecret + "&code=" + code + "&grant_type=authorization_code"

	_, errs := xhttp.NewClient().Get(url).EndStruct(accessToken)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return accessToken, nil
}
