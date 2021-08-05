package config

import (
	"encoding/json"
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"git.dhgames.cn/svr_comm/gmoss/v2"
	"git.dhgames.cn/svr_comm/gmoss/v2/consul"
	"github.com/tidwall/gjson"
)

type All struct {
	WebPort           string                 `json:"web_port"`
	UTLog             string                 `json:"utlog"`
	SmsUrl            string                 `json:"sms_url"`
	SmsSecret         string                 `json:"sms_secret"`
	SmsContent        string                 `json:"sms_content"`
	SmsAppid          string                 `json:"sms_appid"`
	MongoUrl          string                 `json:"mongo_url"`
	MongoOldGpDb      string                 `json:"mongo_old_gpdb"	`
	MongoDb           string                 `json:"mongo_db"`
	JwtSecret         string                 `json:"jwt_secret"`
	HawkEyeOpen       bool                   `json:"hawkeye_open"`
	HawkEyeFilter     string                 `json:"hawkeye_filter"`
	HawkEyeDc         string                 `json:"hawkeye_dc"`
	AntiAddictionOpen bool                   `json:"anti_addiction_open"`
	FacebookGraphUrl  string                 `json:"facebook_graphurl"`
	FacebookInfo      map[string]string      `json:"facebook_infos"`
	Packages          map[string]interface{} `json:"packages"`
}

var staticConfig All

func Init() {
	httpPort := Field("web_port").Int()
	log.Infow("config init ok", "config", httpPort)
	if err := consul.Watch(consul.StaticCfgUrl(), Reload, "kv"); err != nil {
		log.Fatalw("failed to watch config", "err", err)
		panic(err)
	}
}

func Reload(urlI interface{}, configI interface{}) {
	defer func() {
		log.Infow("reload config", "config", staticConfig)
	}()

	configBytes := configI.([]byte)
	if err := json.Unmarshal(configBytes, &staticConfig); err != nil {
		log.Fatalw("failed to init config", "err", err)
		panic(err)
	}
	v := staticConfig
	fmt.Println(v)

	//PackageParam("com.droidhang.aod.cnofficial", "yedun_secret_id")

}
func GetAll() *All { return &staticConfig }

func Packages() interface{} { return staticConfig.Packages }

func FacebookInfos() map[string]string { return staticConfig.FacebookInfo }

// Field 获取静态json数据
func Field(field string) gjson.Result {
	return gjson.GetBytes(gmoss.StaticCfg(), field)
}

func MongoUrl() string {
	mongoUrl := Field("mongo_url").String()
	if mongoUrl == "" {
		mongoUrl = "mongodb://WpxU:WpxU63@10.0.240.19:20294,10.0.240.19:24771/admin?replicaSet=dev-ulogin-db&maxPoolSize=10"
	}
	return mongoUrl
}

func MongoDb() string {
	MongoDb := Field("mongo_db").String()
	if MongoDb == "" {
		MongoDb = "glogin_account"
	}
	return MongoDb
}

func PackageParamRst(bundleId string, key string) gjson.Result {
	packages := Packages()
	mapData := packages.(map[string]interface{})
	if bundleData, ok := mapData[bundleId]; ok {
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
	mapData := packages.(map[string]interface{})
	if bundleData, ok := mapData[bundleId]; ok {
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
