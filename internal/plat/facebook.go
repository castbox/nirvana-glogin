package plat

import (
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/tidwall/gjson"
	"glogin/config"
	"glogin/db"
	"glogin/db/db_core"
	"glogin/pbs/glogin"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
)

const (
	facebookAuthKey = "facebook_oauth_url"
)

var Facebook facebook

type facebook struct{}

// Auth 登录返回第三方账号tokenId openId 错误信息
func (f facebook) Auth(request *glogin.ThirdLoginReq) (string, string, error) {
	//baseUrl := authURL(request.Game.BundleId, facebookAuthKey)
	log.Infow("facebook auth", "request", request)
	baseUrl := config.PackageParam(request.Game.BundleId, facebookAuthKey)
	url := baseUrl + request.ThirdToken
	resp, err := http.Get(url)
	if err != nil {
		resErr := fmt.Errorf("failed communicating with server: %v", err)
		elkAlarm("error", url, resErr)
		return "", "", resErr
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		elkAlarm(resp.Status, url, "")
		return "", "", fmt.Errorf(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resErr := fmt.Errorf("failed reading from metadata server: %v", err)
		elkAlarm("error", url, resErr)
		return "", "", resErr
	}

	data := gjson.GetBytes(body, "data").Map()
	if len(data) == 0 {
		resErr := fmt.Errorf("failed reading data from metadata server: %v", err)
		elkAlarm("error", url, resErr)
		return "", "", resErr
	}

	isValid := data["is_valid"]
	uid := data["user_id"]
	if isValid.Exists() == false || uid.Exists() == false {
		resErr := fmt.Errorf("failed reading isValid from metadata server: %v", err)
		elkAlarm("error", url, resErr)
		return "", "", resErr
	}

	if isValid.Bool() == false {
		return "", "", fmt.Errorf("Invalid Auth Access Token")
	}
	facebookId := uid.String()
	errB, unionId := tokenForBusiness(facebookId, request.ThirdToken, request.Game.BundleId)
	if errB != nil {
		return facebookId, facebookId, nil
	}
	return facebookId, unionId, nil
}

// 获得关联商户 unionId
// https://graph.facebook.com/v9.0/me?fields=token_for_business,name&access_token=EAADuaaYTWgcBAIZC4sOFMAVxuc6LNcM8rGG8iyXZBhZA6LYuAaUM7z8oaoZAozAsrBTBrmm9w1XH4ZA9UUdtUtDpnJ9HBoRZBbBPwLPPbxozbtGKVDU4AgcVf8N2rfV2h6SBkjsvZBjtZA9cR9a9ZALgmPU1vgdPbC7D9ZC3OCDsU1PVf0I1pojT9ZBH1ZAmo0ksgIEyGgpQjtps9aSvZBcHQcgQqXifl1CcZAxJYZD
func tokenForBusiness(faceBookId string, token string, bundleId string) (error, string) {
	doc := db_core.TokenForBusinessData{}
	filter := bson.M{
		"facebook_token": faceBookId,
	}
	errLoad := db.LoadOne(filter, &doc, db.TokenForBusinessTable)
	if errLoad == nil {
		if doc.TokenForBusiness != "" {
			return nil, doc.TokenForBusiness
		}
	}
	baseUrl := config.Field("facebook_for_business_url").String()
	url := baseUrl + token
	resp, err := http.Get(url)
	if err != nil {
		resErr := fmt.Errorf("failed communicating with server: %v", err)
		elkAlarm("error", url, resErr)
		return err, ""
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		elkAlarm(resp.Status, url, "")
		return fmt.Errorf(resp.Status), ""
	}
	//#{<<"token_for_business">> := OpenId, <<"id">> := Id}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resErr := fmt.Errorf("failed reading from metadata server: %v", err)
		elkAlarm("error", url, resErr)
		return resErr, ""
	}
	unionId := gjson.GetBytes(body, "token_for_business").String()
	_, _ = db.AddFbTokenForBusiness(faceBookId, unionId, bundleId)
	return nil, unionId
}

func (f facebook) String() string {
	return "facebook"
}

func (f facebook) DbFieldName() string {
	return "facebook"
}
