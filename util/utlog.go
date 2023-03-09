package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/castbox/nirvana-gcore/glog"
	"github.com/bitly/go-simplejson"
	"glogin/config"
	"io/ioutil"
	"net/http"
)

type HttpOption struct {
	Method string            `json:"method"`
	URL    string            `json:"url" `
	Header map[string]string `json:"header" `
	Body   string            `json:"body" `
}

func HttpTo3rd(option HttpOption) (interface{}, error) {
	url := config.Field("utlog").String()
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(option)
	log.Infow("HttpTo3rd req", "url", url, "option", option)
	resp, err := http.Post(url, "application/json; charset=utf-8", requestBody)
	if err != nil {
		log.Warnw("HttpTo3rd error ", "err", err)
		return nil, nil
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resErr := fmt.Errorf("failed reading from metadata server: %w", err)
		log.Warnw("HttpTo3rd error ", "resErr", resErr)
		return nil, nil
	}
	result, _ := simplejson.NewJson(contents)
	log.Infow("HttpTo3rd rsp", "contents", contents, "result", result)
	return nil, nil
}
