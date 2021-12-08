package config

import (
	"encoding/json"
	"git.dhgames.cn/svr_comm/gcore/consul"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"github.com/tidwall/gjson"
	"strconv"
)

type (
	MongoCfg struct {
		Uri    string `json:"uri"`
		DBName string `json:"dbname"`
	}
	GameCfg struct {
		AppsFlyerIosId          string `json:"appsflyer_ios_id"`
		AppsFlyerAuthentication string `json:"appsflyer_Authentication"`
		AppsFlyerRegistrationId int    `json:"appsflyer_registrationId"`
	}
	PushLogCfg struct {
		Url       string `json:"url"`
		Salt      string `json:"salt"`
		TopicCode int64  `json:"topic_code"`
	}
	LoginConfig struct {
		Ports map[string]int `json:"port"`
		//UTLog             string                 `json:"utlog"`
		PushLog PushLogCfg `json:"push_log"`
		//SmsUrl            string                 `json:"sms_url"`
		//SmsSecret         string                 `json:"sms_secret"`
		//SmsContent        string                 `json:"sms_content"`
		//SmsAppid          string                 `json:"sms_appid"`
		//JwtSecret         string                 `json:"jwt_secret"`
		//HawkEyeOpen       bool                   `json:"hawkeye_open"`
		//HawkEyeFilter     string                 `json:"hawkeye_filter"`
		//HawkEyeDc         string                 `json:"hawkeye_dc"`
		//AntiAddictionOpen bool                   `json:"anti_addiction_open"`
		//FacebookGraphUrl  string                 `json:"facebook_graphurl"`
		FacebookInfo map[string]string      `json:"facebook_infos"`
		Packages     map[string]interface{} `json:"packages"`
		Mongo        map[string]MongoCfg    `json:"mongo"`
		Games        map[string]interface{} `json:"games"`
		Pubs         map[string]interface{} `json:"pubs"`

		PubsByte []byte
	}
)

var staticConfig = LoginConfig{}

func (c *LoginConfig) Reload() {
	log.Infow("reload", "config", c)
	byteData, err := json.Marshal(c.Pubs)
	if err != nil {
		log.Warnw("pubs  marshal err", "pubs", c.Pubs, "err", err)
	}
	c.PubsByte = byteData

	biHost := Field("uds").String()
	log.Infow("reload finish", "biHost", biHost)
}

func WatchStaticConfig() {
	// 本地环境
	serInfo := consul.ServiceInfo{
		Cluster: "lwk_dev",
		Service: "glogin",
		Index:   2,
	}
	err := consul.WatchConfigByService(&serInfo, &staticConfig)
	// 部署环境
	//err := consul.WatchConfig(&staticConfig)
	if err != nil {
		log.Fatalw("failed to wath config", "err", err)
	} else {
		log.Infow("succeed to wath config", "data", staticConfig)
	}

	//log.Infow("config init ok", "config", config.GetAll())
}

func Packages() map[string]interface{} { return staticConfig.Packages }

func PushLog() PushLogCfg { return staticConfig.PushLog }

func Ports() map[string]int { return staticConfig.Ports }

func FacebookInfos() map[string]string { return staticConfig.FacebookInfo }

func Mongo() map[string]MongoCfg { return staticConfig.Mongo }

func Games() map[string]interface{} { return staticConfig.Games }

func Pubs() map[string]interface{} { return staticConfig.Pubs }

// Field 获取静态 pubs json数据
func Field(field string) gjson.Result {
	pubs := Pubs()
	bD, err := json.Marshal(pubs)
	if err != nil {
		log.Warnw("pubs  marshal err", "pubs", pubs, "err", err)
		return gjson.Result{}
	}
	return gjson.GetBytes(bD, field)
}

// WebPort 运维新配置调整
//"mongo_url": "mongodb://WpxU:WpxU63@10.0.240.19:20294,10.0.240.19:24771/admin?replicaSet=dev-ulogin-db&maxPoolSize=10",
//"mongo_db": "ulogin-account-wai",
func WebPort() string {
	ports := Ports()
	if valueData, ok := ports["web_port"]; ok {
		return "0.0.0.0:" + strconv.Itoa(int(valueData))
	}
	return "0.0.0.0:8081"
}

func MongoUrl() string {
	mongoInfo := Mongo()
	mongoUrl := "mongodb://WpxU:WpxU63@10.0.240.19:20294,10.0.240.19:24771/admin?replicaSet=dev-ulogin-db&maxPoolSize=10"
	if VData, ok := mongoInfo["login"]; ok {
		mongoUrl = VData.Uri
	}
	return mongoUrl
}

func MongoDb() string {
	mongoInfo := Mongo()
	mongoDb := "glogin_account"
	if VData, ok := mongoInfo["login"]; ok {
		mongoDb = VData.DBName
	}
	return mongoDb
}

func MongoOldDb() string {
	mongoInfo := Mongo()
	mongoOldDb := "glogin_account_old"
	if VData, ok := mongoInfo["login2"]; ok {
		mongoOldDb = VData.DBName
	}
	return mongoOldDb
}

func PackageParamRst(bundleId string, key string) gjson.Result {
	packages := Packages()
	if bundleData, ok := packages[bundleId]; ok {
		bD, err := json.Marshal(bundleData)
		if err != nil {
			log.Warnw("PackageParamRst  marshal err", "bundleData", bundleData, "err", err)
			return gjson.Result{}
		}
		mapResult := gjson.ParseBytes(bD).Map()
		return mapResult[key]
	}
	return gjson.Result{}
}

func PackageParam(bundleId string, key string) string {
	packages := Packages()
	if bundleData, ok := packages[bundleId]; ok {
		if value, ok := bundleData.(map[string]interface{})[key]; ok {
			return value.(string)
		}
	}
	return ""
}

func FacebookParam(bundleKey string) string {
	facebookInfos := FacebookInfos()
	if valueData, ok := facebookInfos[bundleKey]; ok {
		return valueData
	}
	return ""
}

func GameParamRst(gameCd string, key string) gjson.Result {
	games := Games()
	if gameData, ok := games[gameCd]; ok {
		bD, err := json.Marshal(gameData)
		if err != nil {
			log.Warnw("gameDataRst  marshal err", "gameData", gameData, "err", err)
			return gjson.Result{}
		}
		mapResult := gjson.ParseBytes(bD).Map()
		return mapResult[key]
	}
	return gjson.Result{}
}
