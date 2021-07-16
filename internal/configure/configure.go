package configure

import (
	"encoding/json"
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"git.dhgames.cn/svr_comm/gmoss/v2"
	"glogin/utils"
	"strings"
	"sync"
)

var pubConfig sync.Map

func GetCfg(appName string, vsn string, ip string) *ClusterCfg {
	v, ok := pubConfig.Load(appName)
	if !ok {
		log.Infow("appName GetCfg error", "appName", appName)
		return nil
	} else {
		fmt.Println(v)
		p := v.(PubCfg)
		if len(p.WhiteList) > 0 {
			for _, wIp := range p.WhiteList {
				if ip == wIp {
					return &p.ConnCfg.Test
				}
			}
		}
		if p.CheckVsn == vsn {
			return &p.ConnCfg.Check
		}
		return &p.ConnCfg.Stable
	}
	return nil
}

func WatchPubCfg() {
	cluster, _, _, _ := gmoss.WhoIAm()
	//const dynamicPubUrl = "kv/app_dynamic_cfg/lwk_dev/pub_cfg/?recurse=true"
	//http://127.0.0.1:2000/v1/kv/app_static_cfg/lwk_dev/glogin/pub_cfg/?recurse=true
	dynamicPubUrl := fmt.Sprintf("kv/app_dynamic_cfg/%s/pub_cfg/?recurse=true", cluster)
	gmoss.DynamicCfgCustom(dynamicPubUrl, WatchPubCfgCallBack)
}

//[{"LockIndex":0,"Key":"app_dynamic_cfg/lwk_dev/pub_cfg/com.dh.bpc.gp","Flags":0,"Value":"ewogICAgImNvbm5fY2ZnIjp7CiAgICAgICAgInN0YWJsZSI6ewoJCQkJCQkiY2x1c3Rlcl90eXBlIiA6IDEsCiAgICAgICAgICAgICJ1Z2F0ZV9hZGRyIjoiMTAuMC4yNDAuMjM0OjE4ODg5IiwKICAgICAgICAgICAgInZzbl9hZGRyIjoiMTAuMC4yNDAuMjM0OjE5ODg5IiwKICAgICAgICAgICAgInVwYXlfYWRkciI6IjEwLjAuMjQwLjE5OjgwODgiLAogICAgICAgICAgICAidWNoYXRfYWRkciI6ICJhb2QtZGV2LXVjaGF0LmRoZ2FtZXMuY246MTg4ODciLAogICAgICAgICAgICAiY29tbXVuaXR5X3dlYl9hZGRyIjoiaHR0cDovLzEwLjAuMC4xOTo3NzcwL2FvZC92MS9pbmRleC5odG1sIy9ob21lIiwKICAgICAgICAgICAgImNvbW11bml0eV9zcnZfYWRkciI6Imh0dHA6Ly8xODIuMTUwLjIyLjYxOjI3Nzc4IiwKICAgICAgICAgIAkiYWljc193c19hZGRyIjogIndzczovL2FpY3MtY2xpLmRldi1kaC5jb20vZGhfd3Mvd3NfYXBwIiwKICAgICAgICAgICAgImFpY3NfaHR0cF9hZGRyIjogImh0dHBzOi8vYWljcy1jbGkuZGV2LWRoLmNvbSIKICAgICAgICB9LAogICAgICAgICJjaGVjayI6ewogICAgICAgICAgICAiY2x1c3Rlcl90eXBlIjogMiwKICAgICAgICAgICAgInVnYXRlX2FkZHIiOiIiLAogICAgICAgICAgICAidnNuX2FkZHIiOiAiIiwKICAgICAgICAgICAgInVwYXlfYWRkciI6IiIsCiAgICAgICAgICAgICJ1Y2hhdF9hZGRyIjogIiIKICAgICAgICB9CiAgICB9LAogICAgImNoZWNrX3ZzbiI6IiIKfQ==","CreateIndex":16921314,"ModifyIndex":16921314},{"LockIndex":0,"Key":"app_dynamic_cfg/lwk_dev/pub_cfg/com.droidhang.bpc.ios","Flags":0,"Value":"ewogICAgImNvbm5fY2ZnIjp7CiAgICAgICAgInN0YWJsZSI6ewoJCQkJCQkiY2x1c3Rlcl90eXBlIiA6IDEsCiAgICAgICAgICAgICJ1Z2F0ZV9hZGRyIjoiMTAuMC4yNDAuMjM0OjE4ODg5IiwKICAgICAgICAgICAgInZzbl9hZGRyIjoiMTAuMC4yNDAuMjM0OjE5ODg5IiwKICAgICAgICAgICAgInVwYXlfYWRkciI6Imh0dHA6Ly8xMC4wLjI0MC4xOTo4MDg4IiwKICAgICAgICAgICAgInVjaGF0X2FkZHIiOiAiYW9kLWRldi11Y2hhdC5kaGdhbWVzLmNuOjE4ODg3IiwKICAgICAgICAgICAgImNvbW11bml0eV93ZWJfYWRkciI6Imh0dHA6Ly8xMC4wLjAuMTk6Nzc3MC9hb2QvdjEvaW5kZXguaHRtbCMvaG9tZSIsCiAgICAgICAgICAgICJjb21tdW5pdHlfc3J2X2FkZHIiOiJodHRwOi8vMTgyLjE1MC4yMi42MToyNzc3OCIsCiAgICAgICAgICAJImFpY3Nfd3NfYWRkciI6ICJ3c3M6Ly9haWNzLWNsaS5kZXYtZGguY29tL2RoX3dzL3dzX2FwcCIsCiAgICAgICAgICAgICJhaWNzX2h0dHBfYWRkciI6ICJodHRwczovL2FpY3MtY2xpLmRldi1kaC5jb20iCiAgICAgICAgfSwKICAgICAgICAiY2hlY2siOnsKICAgICAgICAgICAgImNsdXN0ZXJfdHlwZSI6IDIsCiAgICAgICAgICAgICJ1Z2F0ZV9hZGRyIjoiIiwKICAgICAgICAgICAgInZzbl9hZGRyIjogIiIsCiAgICAgICAgICAgICJ1cGF5X2FkZHIiOiIiLAogICAgICAgICAgICAidWNoYXRfYWRkciI6ICIiCiAgICAgICAgfQogICAgfSwKICAgICJjaGVja192c24iOiIiCn0=","CreateIndex":16821432,"ModifyIndex":16821481}]
//	[
//	{
//		"LockIndex": 0,
//		"Key": "app_static_cfg/lwk_dev/glogin/pub_cfg/com.dh.bpc.gp",
//		"Flags": 0,
//		"Value": "ewogICAgICAgICAgICAiZmFjZWJvb2tfb2F1dGhfdXJsIjogImh0dHBzOi8vZ3JhcGguZmFjZWJvb2suY29tL2RlYnVnX3Rva2VuP2FjY2Vzc190b2tlbj02MjI0MDAyNjUzNjMyMDIlN0MyMjY4ODgxNDEwYjNmYTA5MzNjYjJiMGVjZWFhNTE2ZSZpbnB1dF90b2tlbj0iLAogICAgICAgICAgICAiYXBwc2ZseWVyX0FORFJPSUQiOiAiaHR0cHM6Ly9hcGkyLmFwcHNmbHllci5jb20vaW5hcHBldmVudC9jb20uZHJvaWRoYW5nLmFvZC5ncCIsCiAgICAgICAgICAgICJhcHBzZmx5ZXJfb3BlbiI6IHRydWUsCiAgICAgICAgICAgICJhcHBzZmx5ZXJfSU9TIjogImh0dHBzOi8vYXBpMi5hcHBzZmx5ZXIuY29tL2luYXBwZXZlbnQvaWQxMTUzNDYxOTE1IiwKICAgICAgICAgICAgImFwcHNmbHllcl9BdXRoZW50aWNhdGlvbiI6ICIzNkZmTmsyNDR4aTlCQ3hFVVJxYTVuIiwKICAgICAgICAgICAgImFwcHNmbHllcl9yZWdpc3RyYXRpb25JZCI6IDEyCiB9",
//		"CreateIndex": 21120571,
//		"ModifyIndex": 21120586
//	}
//]
func WatchPubCfgCallBack(path interface{}, value interface{}) {
	//p := value.(PubCfg)
	path2 := path.(string)
	switch val := value.(type) {
	case []byte:
		pList := PubCfgList{}.AppList
		err := json.Unmarshal(val, &pList)
		if err != nil {
			log.Errorw("WatchPubCfgCallBack ", "path2", err)
			return
		}
		//AppName = lists:last(binary:split(Key, <<"/">>, [global])),
		//Info = jiffy:decode(base64:decode(Value), [return_maps]),
		for _, baseValue := range pList {
			fmt.Println(baseValue.Key)
			keySlice := strings.Split(baseValue.Key, "/")
			appName := keySlice[len(keySlice)-1]
			fmt.Printf("appName : %vn", appName)
			jsonData, err := utils.Base64Decode(baseValue.Value)
			if err != nil {
				return
			}
			bInfo := PubCfg{}
			err2 := json.Unmarshal(jsonData, &bInfo)
			if err2 != nil {
				log.Errorw("WatchPubCfgCallBack ", "path2", err)
				return
			}
			fmt.Println(bInfo)
			pubConfig.Store(appName, bInfo)
		}
		pubConfig.Range(allBundles)
	}
	log.Infow("reload config", "path2", path2)
	return
}

func allBundles(key, value interface{}) bool {
	fmt.Println("appName =", key, "JsonData =", value)
	return true
}
