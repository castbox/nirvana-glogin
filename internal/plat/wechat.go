package plat

import (
	"fmt"
	log "gitlab.degames.cn/svr_comm/gcore/glog"
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

type CheckAccessTokenRsp struct {
	Errcode int    `json:"errcode,omitempty"` // 错误码
	Errmsg  string `json:"errmsg,omitempty"`  // 错误信息
}

// 微信开放平台用户信息
type Oauth2UserInfo struct {
	Openid     string   `json:"openid,omitempty"`     // 普通用户的标识，对当前开发者帐号唯一
	Nickname   string   `json:"nickname,omitempty"`   // 普通用户昵称
	Sex        int      `json:"sex,omitempty"`        // 普通用户性别，1为男性，2为女性
	City       string   `json:"city,omitempty"`       // 普通用户个人资料填写的城市
	Province   string   `json:"province,omitempty"`   // 普通用户个人资料填写的省份
	Country    string   `json:"country,omitempty"`    // 国家，如中国为CN
	Headimgurl string   `json:"headimgurl,omitempty"` // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	Privilege  []string `json:"privilege,omitempty"`  // 用户特权信息，json数组，如微信沃卡用户为（chinaunicom）
	Unionid    string   `json:"unionid,omitempty"`    // 用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。
}

var WeChat wechat

type wechat struct{}

// Auth 登录返回第三方账号tokenId openId 错误信息
func (w wechat) Auth(request *glogin.ThirdLoginReq) (*AuthRsp, error) {
	log.Infow("wechat auth req", "req", request)
	appId := config.PackageParam(request.Game.BundleId, "wx_app_id")
	appSecret := config.PackageParam(request.Game.BundleId, "wx_app_secret")
	code := request.ThirdToken
	//log.Infow("wechat auth appInfo", "appId", appId, "appSecret", appSecret, "code", code)
	accessToken, err := GetOauth2AccessToken(appId, appSecret, code)
	if err != nil {
		log.Warnw("wechat auth error ", "err", err)
		return nil, err
	}
	log.Infow("wechat auth rsp", "rsp", accessToken)
	if accessToken.Errcode != 0 {
		resErr := fmt.Errorf("wechat auth error code: %d, errmsg: %s", accessToken.Errcode, accessToken.Errmsg)
		return nil, resErr
	}
	if accessToken.Unionid == "" {
		log.Infow("wechat auth unionid is nil", "Openid", accessToken.Openid)
		accessToken.Unionid = accessToken.Openid
		if accessToken.Unionid == "" {
			return nil, fmt.Errorf("wechat auth error unionid and opendid is nil")
		}
	}
	// Nick
	oauth2UserInfo, errUserInfo := GetOauth2UserInfo(accessToken.AccessToken, accessToken.Openid)
	if errUserInfo != nil {
		log.Warnw("wechat auth GetOauth2UserInfo error ", "errUserInfo", errUserInfo)
		return nil, errUserInfo
	}
	log.Infow("wechat auth GetOauth2UserInfo rsp", "rsp", errUserInfo)
	return &AuthRsp{
		Uid:     accessToken.Openid,
		UnionId: accessToken.Unionid,
		Nick:    oauth2UserInfo.Nickname,
	}, nil
}

func (w wechat) String() string {
	return "wechat"
}

func (w wechat) DbFieldName() string {
	return "wechat"
}

//  GetOauth2AccessToken 微信第三方登录，code 换取 access_token
//	appId：应用唯一标识，在微信开放平台提交应用审核通过后获得
//	appSecret：应用密钥AppSecret，在微信开放平台提交应用审核通过后获得
//	code：App用户换取access_token的code
//	文档：https://developers.weixin.qq.com/doc/oplatform/Mobile_App/WeChat_Login/Development_Guide.html
func GetOauth2AccessToken(appId, appSecret, code string) (accessToken *Oauth2AccessToken, err error) {
	accessToken = new(Oauth2AccessToken)
	url := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + appId + "&secret=" + appSecret + "&code=" + code + "&grant_type=authorization_code"
	log.Infow("wechat auth url", "url", url)
	_, errs := xhttp.NewClient().Get(url).EndStruct(accessToken)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return accessToken, nil
}

//  RefreshOauth2AccessToken 刷新微信第三方登录后，获取到的 access_token
//	appId：应用唯一标识，在微信开放平台提交应用审核通过后获得
//	refreshToken：填写通过获取 access_token 获取到的 refresh_token 参数
//	文档：https://developers.weixin.qq.com/doc/oplatform/Mobile_App/WeChat_Login/Development_Guide.html
func RefreshOauth2AccessToken(appId, refreshToken string) (accessToken *Oauth2AccessToken, err error) {
	accessToken = new(Oauth2AccessToken)
	url := "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=" + appId + "&grant_type=refresh_token&refresh_token=" + refreshToken

	_, errs := xhttp.NewClient().Get(url).EndStruct(accessToken)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return accessToken, nil
}

// CheckOauth2AccessToken 检验授权凭证（access_token）是否有效
//	accessToken：调用接口凭证
//	openid：普通用户标识，对该公众帐号唯一，获取 access_token 是获取的
//	文档：https://developers.weixin.qq.com/doc/oplatform/Mobile_App/WeChat_Login/Authorized_API_call_UnionID.html
func CheckOauth2AccessToken(accessToken, openid string) (result *CheckAccessTokenRsp, err error) {
	result = new(CheckAccessTokenRsp)
	url := "https://api.weixin.qq.com/sns/auth?access_token=" + accessToken + "&openid=" + openid

	_, errs := xhttp.NewClient().Get(url).EndStruct(result)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return result, nil
}

// GetOauth2UserInfo 微信开放平台：获取用户个人信息
//	accessToken：接口调用凭据
//	openId：用户的OpenID
//	lang:默认为 zh_CN ，可选填 zh_CN 简体，zh_TW 繁体，en 英语
//	文档：https://developers.weixin.qq.com/doc/oplatform/Mobile_App/WeChat_Login/Authorized_API_call_UnionID.html
func GetOauth2UserInfo(accessToken, openId string, lang ...string) (userInfo *Oauth2UserInfo, err error) {
	userInfo = new(Oauth2UserInfo)
	url := "https://api.weixin.qq.com/sns/userinfo?access_token=" + accessToken + "&openid=" + openId
	if len(lang) == 1 {
		url += "&lang=" + lang[0]
	}
	_, errs := xhttp.NewClient().Get(url).EndStruct(userInfo)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return userInfo, nil
}
