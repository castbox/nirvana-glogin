package ids

import (
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/tidwall/gjson"
	"glogin/config"
	"io/ioutil"
	"net/http"
	"strings"
)

// 获得应用token
//"https://graph.facebook.com/oauth/access_token?client_id={your-app-id}&client_secret={your-app-secret} &grant_type=client_credentials"
func GetAccessToken(bundleID string) (string, error) {
	faceBookInfos := strings.Split(config.FacebookParam(bundleID), "|")
	appId := faceBookInfos[0]
	appSecret := faceBookInfos[1]
	fullURL := config.Field("facebook_graphurl").String() + "oauth/access_token?client_id=" + appId + "&client_secret=" + appSecret + "&grant_type=client_credentials"
	log.Infow("getAccessToken fullURL ", "fullURL", fullURL)
	resp, err := http.Get(fullURL)
	if err != nil {
		resErr := fmt.Errorf("failed communicating with server: %v", err)
		log.Warnw("elkAlarm http", "status", "error", "url", fullURL)
		return "", resErr
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Warnw("elkAlarm http", "status", resp.Status, "url", fullURL)
		return "", fmt.Errorf(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resErr := fmt.Errorf("failed reading from metadata server: %v", err)
		log.Warnw("elkAlarm http", "status", "error", "url", fullURL)
		return "", resErr
	}

	accessToken := gjson.GetBytes(body, "access_token").String()
	log.Infow("getAccessToken access_token ", "accessToken", accessToken)
	return accessToken, nil
}

// 获得关联商户其他appID对应Facebook UID
// curl -X GET https://graph.facebook.com/488372655904650/ids_for_apps?access_token=262137525262855|OnXXyzEqcH2FHLpJkNQGu_7aPNM
func GetIds(faceBookId string, oldBundleId string, accessToken string) (string, error) {
	log.Infow("getAccessToken GetIds ", "faceBookId", faceBookId, "oldBundleId", oldBundleId, "accessToken", accessToken)
	fullURL := config.Field("facebook_graphurl").String() + faceBookId + "/ids_for_apps?access_token=" + accessToken
	log.Infow("GetIds fullURL ", "fullURL", fullURL)
	resp, err := http.Get(fullURL)
	if err != nil {
		resErr := fmt.Errorf("failed communicating with server: %v", err)
		log.Warnw("elkAlarm http", "status", "error", "url", fullURL)
		return "", resErr
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Warnw("elkAlarm http", "status", resp.Status, "url", fullURL)
		return "", fmt.Errorf(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resErr := fmt.Errorf("failed reading from metadata server: %v", err)
		log.Warnw("elkAlarm http", "status", "error", "url", fullURL)
		return "", resErr
	}

	data := gjson.GetBytes(body, "data").Array()
	log.Infow("getAccessToken GetIds ", "data", data)
	if len(data) == 0 {
		resErr := fmt.Errorf("failed reading data from metadata server: %v", err)
		log.Warnw("elkAlarm http", "status", "error", "url", fullURL)
		return "", resErr
	}

	faceBookInfos := strings.Split(config.FacebookParam(oldBundleId), "|")
	appId := faceBookInfos[0]
	//appSecret := faceBookInfos[1]
	for _, v := range data {
		log.Infow("GetIds fullURL 33333", "v", v)
		appInfo := v.Get("app")
		id := appInfo.Get("id").String()
		if appId == id {
			log.Infow("GetIds fullURL 44444", "appInfo2", appInfo)
			return v.Get("id").String(), nil
		}
	}
	return "", nil
}
