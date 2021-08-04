package plat

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"glogin/pbs/glogin"
	"hash"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"time"
)

var IOS ios

type ios struct{}

// Auth 登录返回第三方账号tokenId openId 错误信息
func (i ios) Auth(request *glogin.ThirdLoginReq) (*AuthRsp, error) {
	if request.ThirdToken == "" {
		return nil, ErrInvalidIdentityToken
	}
	appleToken, err := parseToken(request.ThirdToken)
	if err != nil {
		return nil, err
	}
	key, err := fetchKeysFromApple(appleToken.header.Kid)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, ErrFetchKeysFail
	}

	pubKey, err := generatePubKey(key.N, key.E)
	if err != nil {
		return nil, err
	}

	//利用获取到的公钥解密token中的签名数据
	sig, err := decodeSegment(appleToken.sign)
	if err != nil {
		return nil, err
	}

	//苹果使用的是SHA256
	var h hash.Hash
	switch appleToken.header.Alg {
	case "RS256":
		h = crypto.SHA256.New()
	case "RS384":
		h = crypto.SHA384.New()
	case "RS512":
		h = crypto.SHA512.New()
	}
	if h == nil {
		return nil, ErrInvalidHashType
	}

	h.Write([]byte(appleToken.headerStr + "." + appleToken.claimsStr))

	if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, h.Sum(nil), sig); err != nil {
		return nil, err
	}

	if ok, err := appleToken.IsValid(); !ok || err != nil {
		return nil, err
	}
	unionId := appleToken.claims.Sub + "_cnofficial"

	return &AuthRsp{
		Uid:     unionId,
		UnionId: unionId,
	}, nil
}

func (i ios) String() string {
	return "ios"
}
func (i ios) DbFieldName() string {
	return "ios"
}

func parseToken(token string) (*appleToken, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidTokenFormat
	}
	//header
	var apToken = &appleToken{
		headerStr: parts[0],
		claimsStr: parts[1],
		sign:      parts[2],
	}
	var headerBytes []byte
	var err error
	if headerBytes, err = decodeSegment(parts[0]); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(headerBytes, &apToken.header); err != nil {
		return nil, err
	}

	//claims
	var claimBytes []byte
	if claimBytes, err = decodeSegment(parts[1]); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(claimBytes, &apToken.claims); err != nil {
		return nil, err
	}
	return apToken, nil
}

func fetchKeysFromApple(kid string) (*appleKey, error) {
	rsp, err := http.Get("https://appleid.apple.com/auth/keys")
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching keys from apple server fail: %d", rsp.StatusCode)
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	type Keys struct {
		Keys []*appleKey `json:"keys"`
	}

	var ks *Keys
	var result *appleKey
	if err = json.Unmarshal(data, &ks); err != nil {
		return nil, err
	}
	for _, k := range ks.Keys {
		if k == nil {
			continue
		}
		if k.Kid == kid {
			result = k
			break
		}
	}
	return result, nil
}

func generatePubKey(nStr, eStr string) (*rsa.PublicKey, error) {
	nBytes, err := decodeBase64String(nStr)
	if err != nil {
		return nil, err
	}
	eBytes, err := decodeBase64String(eStr)
	if err != nil {
		return nil, err
	}

	n := &big.Int{}
	n.SetBytes(nBytes)
	e := &big.Int{}
	e.SetBytes(eBytes)

	var pub = rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}
	return &pub, nil
}

func decodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}

func decodeBase64String(src string) ([]byte, error) {
	var isRaw = !strings.HasSuffix(src, "=")
	if strings.Contains(src, "+/") {
		if isRaw {
			return base64.RawStdEncoding.DecodeString(src)
		}
		return base64.StdEncoding.DecodeString(src)
	}
	if isRaw {
		return base64.RawURLEncoding.DecodeString(src)
	}
	return base64.URLEncoding.DecodeString(src)
}

var (
	ErrInvalidHashType      = fmt.Errorf("invalid hash type")
	ErrInvalidTokenFormat   = fmt.Errorf("invalid token")
	ErrFetchKeysFail        = fmt.Errorf("invalid rsa public key")
	ErrInvalidClientID      = fmt.Errorf("invalid client_id")
	ErrInvalidClientSecret  = fmt.Errorf("invalid client_secret")
	ErrInvalidRedirectURI   = fmt.Errorf("invalid redirect_uri")
	ErrTokenExpired         = fmt.Errorf("token expired")
	ErrInvalidIssValue      = fmt.Errorf("invalid iss value")
	ErrInvalidRefreshToken  = fmt.Errorf("invalid refresh token")
	ErrInvalidIdentityCode  = fmt.Errorf("invalid identity code")
	ErrInvalidIdentityToken = fmt.Errorf("invalid identity token")
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"` //固定值: bearer
}

type appleKey struct {
	Kid string `json:"kid"` //公钥ID
	Alg string `json:"alg"` //签名算法
	Kty string `json:"kty"` //加密算法
	E   string `json:"e"`   //RSA公钥指数值
	N   string `json:"n"`   //RSA公钥模数值
	Use string `json:"use"` //
}

type appleHeader struct {
	Kid string `json:"kid"` //apple公钥的密钥ID
	Alg string `json:"alg"` //签名token的算法
}

type appleToken struct {
	header    *appleHeader //header
	headerStr string
	claims    *appleClaim //claims
	claimsStr string
	sign      string //签名
}

func (t *appleToken) Kid() string {
	if t == nil || t.claims == nil {
		return ""
	}
	return t.header.Kid
}

func (t *appleToken) Alg() string {
	if t == nil || t.claims == nil {
		return ""
	}
	return t.header.Alg
}

func (t *appleToken) Iss() string {
	if t == nil || t.claims == nil {
		return ""
	}
	return t.claims.Iss
}

func (t *appleToken) Aud() string {
	if t == nil || t.claims == nil {
		return ""
	}
	return t.claims.Aud
}

func (t *appleToken) Exp() int64 {
	if t == nil || t.claims == nil {
		return 0
	}
	return t.claims.Exp
}

func (t *appleToken) Iat() int64 {
	if t == nil || t.claims == nil {
		return 0
	}
	return t.claims.Iat
}

func (t *appleToken) Sub() string {
	if t == nil || t.claims == nil {
		return ""
	}
	return t.claims.Sub
}

func (t *appleToken) CHash() string {
	if t == nil || t.claims == nil {
		return ""
	}
	return t.claims.CHash
}

func (t *appleToken) AuthTime() int64 {
	if t == nil || t.claims == nil {
		return 0
	}
	return t.claims.AuthTime
}

func (t *appleToken) Email() string {
	if t == nil || t.claims == nil {
		return ""
	}
	return t.claims.Email
}

func (t *appleToken) EmailVerified() bool {
	if t == nil || t.claims == nil {
		return false
	}
	return t.claims.EmailVerified
}

func (t *appleToken) NonceSupported() bool {
	if t == nil || t.claims == nil {
		return false
	}
	return t.claims.NonceSupported
}

func (t *appleToken) IsPrivateEmail() bool {
	if t == nil || t.claims == nil {
		return false
	}
	return t.claims.IsPrivateEmail
}

func (t *appleToken) RealUserStatus() int {
	if t == nil || t.claims == nil {
		return 0
	}
	return t.claims.RealUserStatus
}

func (t *appleToken) Nonce() string {
	if t == nil || t.claims == nil {
		return ""
	}
	return t.claims.Nonce
}

func (t *appleToken) IsValid() (bool, error) {
	if t == nil || t.claims == nil {
		return false, ErrInvalidTokenFormat
	}
	if t.claims.Iss != "https://appleid.apple.com" {
		return false, ErrInvalidIssValue
	}
	var now = time.Now().Unix()
	if t.claims.Exp < now {
		return false, ErrTokenExpired
	}
	if t.claims.Iat > now {
		return false, ErrTokenExpired
	}
	return true, nil
}

func (t *appleToken) String() string {
	var hStr, cStr string
	if t.header != nil {
		hStr = fmt.Sprintf("%+v", *t.header)
	}
	if t.claims != nil {
		cStr = fmt.Sprintf("%+v", *t.claims)
	}
	return fmt.Sprintf("Header: [%s], Claims: [%s], Sign: [%s]\n", hStr, cStr, t.sign)
}

type appleClaim struct {
	Iss            string `json:"iss"`   //签发者，固定值: https://appleid.apple.com
	Sub            string `json:"sub"`   //用户唯一标识
	Aud            string `json:"aud"`   //App ID
	Iat            int64  `json:"iat"`   //token生成时间
	Exp            int64  `json:"exp"`   //token过期时间
	Nonce          string `json:"nonce"` //客户端设置的随机值
	NonceSupported bool   `json:"nonce_supported"`
	Email          string `json:"email"` //邮件
	EmailVerified  bool   `json:"email_verified"`
	IsPrivateEmail bool   `json:"is_private_email"`
	RealUserStatus int    `json:"real_user_status"`
	CHash          string `json:"c_hash"`    //
	AuthTime       int64  `json:"auth_time"` //验证时间
}
