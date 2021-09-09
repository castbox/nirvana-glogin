package appsflyer

import (
	"encoding/json"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"glogin/config"
	"glogin/constant"
	"glogin/internal"
	"glogin/util"
	"time"
)

func SendAppsFlyer(req internal.Req) (interface{}, error) {
	bAppsFlyerOpen := config.PackageParamRst(req.Game.BundleId, "appsflyer_open").Bool()
	if !bAppsFlyerOpen {
		return nil, nil
	}
	if req.GameCd == "" {
		req.GameCd = req.Game.GameCd
	}
	urlBase := ""
	if req.Game.Platform == constant.ANDROID {
		urlBase = constant.AppsFlyerANDROID + req.Game.BundleId
	} else {
		//先从package参数去拿，如果没有 去game参数去拿，如果还没有就填默认
		appsFlyerIosId := "id1153461915"
		appsFlyerIosId = config.PackageParamRst(req.Game.BundleId, "appsflyer_ios_id").String()
		if appsFlyerIosId == "" {
			appsFlyerIosId = config.GameParamRst(req.GameCd, "appsflyer_ios_id").String()
		}
		urlBase = constant.AppsFlyerIOS + appsFlyerIosId
	}

	body := map[string]string{
		"appsflyer_id":   req.Game.AppsflyerId,
		"advertising_id": req.Game.AdvertisingId,
		"bundle_id":      req.Game.BundleId,
		"eventName":      "registration",
		"eventCurrency":  "USD",
		"ip":             req.IP,
		"eventTime":      util.FormatDate(time.Now(), util.YYYY_MM_DD_HH_MM_SS_SSS),
		"af_events_api":  "true",
	}
	appsFlyerAuthentication := "36FfNk244xi9BCxEURqa5n"
	appsflyerRegistrationId := 12

	appsFlyerAuthentication = config.PackageParamRst(req.Game.BundleId, "appsflyer_Authentication").String()
	if appsFlyerAuthentication == "" {
		appsFlyerAuthentication = config.GameParamRst(req.GameCd, "appsflyer_Authentication").String()
	}
	appsflyerRegistrationId = int(config.PackageParamRst(req.Game.BundleId, "appsflyer_registrationId").Int())
	if appsflyerRegistrationId == 0 {
		appsflyerRegistrationId = int(config.GameParamRst(req.GameCd, "appsflyer_registrationId").Int())
	}
	bm := map[string]interface{}{
		"id": appsflyerRegistrationId,
	}
	eventVal, err := json.Marshal(bm)
	if err == nil {
		//eventVal := jiffy:encode(#{'id' => get_EventVal(BundleId)})
		body["eventValue"] = string(eventVal)
	}
	bodyJson, err2 := json.Marshal(body)
	if err2 != nil {
		log.Warnw("SendAppsFlyer marshal bodyJson err", "urlBase", urlBase, "body", body, "err", err)
	}
	strBodyJson := string(bodyJson)
	log.Infow("SendAppsFlyer Info", "urlBase", urlBase, "strBodyJson", strBodyJson)
	// send to  utlog
	return util.HttpTo3rd(util.HttpOption{
		Method: "post",
		Header: map[string]string{
			"Content-Type":   "application/json; charset=utf-8",
			"Authentication": appsFlyerAuthentication,
		},
		URL:  urlBase,
		Body: strBodyJson,
	})
}
