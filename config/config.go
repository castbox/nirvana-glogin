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
	WebPort           string `json:"web_port"`
	UTLog             string `json:"utlog"`
	TestingOpen       bool   `json:"testing_open"`
	SmsUrl            string `json:"sms_url"`
	SmsSecret         string `json:"sms_secret"`
	SmsContent        string `json:"sms_content"`
	SmsAppid          string `json:"sms_appid"`
	MongoUrl          string `json:"mongo_url"`
	MongoOldGpDb      string `json:"mongo_old_gpdb"	`
	MongoDb           string `json:"mongo_db"`
	JwtSecret         string `json:"jwt_secret"`
	HawkEyeOpen       bool   `json:"hawkeye_open"`
	HawkEyeFilter     string `json:"hawkeye_filter"`
	HawkEyeDc         string `json:"hawkeye_dc"`
	AntiAddictionOpen bool   `json:"anti_addiction_open"`
	FacebookGraphUrl  string `json:"facebook_graphurl"`
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

}
func GetAll() *All { return &staticConfig }

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

func BundleInfo(bundleId string) string {
	bundles := Field("bundles").String()
	return bundles
}
