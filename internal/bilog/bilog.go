package bilog

import (
	"encoding/json"
	"fmt"
	log "github.com/castbox/nirvana-gcore/glog"
	"glogin/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Log struct {
	EventCode  string `json:"event_code"`
	EventType  string `json:"event_type"`
	EventName  string `json:"event_name"`
	GameCd     string `json:"game_cd"`
	CreateTs   string `json:"create_ts"`
	EventValue string `json:"event_value,omitempty"`
}

func (l *Log) Push() {
	// 暂时注释楼上的BI打点
	return
	if err := l.push(); err != nil {
		log.Warnw("bilog push 失败", "err", err)
		return
	}
}

func (l *Log) push() error {
	before := time.Now().UnixNano()
	dataLog, err := json.Marshal(l)
	if err != nil {
		return err
	}
	bodyString := fmt.Sprintf("data=%s", url.QueryEscape(string(dataLog)))
	res, err := http.Post(config.PushLog().Url, "application/x-www-form-urlencoded", strings.NewReader(bodyString))
	// res, err := http.Post("http://ulog-inner-test.dhgames.cn:8180/inner/push_log", "application/x-www-form-urlencoded", strings.NewReader(bodyString))
	defer func() {
		if err == nil {
			res.Body.Close()
		}
		log.Infow("bilog_push", "time_cost", (time.Now().UnixNano()-before)/1000000)
	}()
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("log post error,status_code=%v", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	if result["error_msg"] != "ok" {
		//errRsp := fmt.Errorf("CheckSmsVerifyCode fail phone:%v ,dbFind:%v", phone, result)
		return fmt.Errorf("请求ulog失败，对方访问err_msg不为ok，error=%v", result)
	}

	f := string(dataLog)
	fmt.Println(string(f))
	//LocalLog.Debugw("succeed to push log", "content", string(dataLog))
	//fmt.Println("ok")
	return nil
}
