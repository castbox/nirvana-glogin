package hawkeye

import (
	"fmt"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	api "git.dhgames.cn/svr_comm/hawkeye/v2/api"
	"git.dhgames.cn/svr_comm/hawkeye/v2/pbs/pbgo"
	"git.dhgames.cn/svr_comm/kite"
	"glogin/config"
	"glogin/constant"
	"glogin/internal"
	"strings"
	"time"
)

// AOD特殊处理（AOD台湾审核和内网unity dev等需绕开鹰眼检查）
func channelFilter(channel string) bool {
	filterCfg := config.Field("hawkeye_filter").String()
	filters := strings.Split(filterCfg, "|")
	log.Infow("hawkeye_filter marks ", "channel", channel, "filters", filters)
	for _, value := range filters {
		if value == channel {
			return true
		}
	}
	return false
}

func CheckLogin(req internal.Req) (interface{}, error) {
	if !config.Field("hawkeye_open").Bool() {
		return nil, nil
	}
	if channelFilter(req.Game.Channel) {
		return nil, nil
	}

	loginReq := &pbgo.LoginRequest{
		GameCd: req.Game.GameCd,
		DeviceInfo: &pbgo.DeviceInfo{
			SmId:        req.Client.Dhid,
			Ip:          req.IP,
			AppsflyerId: req.Game.AppsflyerId,
		},
		UserInfo: &pbgo.UserInfo{
			BundleId: req.Game.BundleId,
			Account:  req.Account,
		},
	}
	cfgDc := strings.Split(config.Field(constant.HawkEyeDcCluster).String(), "|")
	log.Infow("HawkeyeLogin Req", "loginReq", loginReq)
	rsp, err := api.Login(cfgDc[0], cfgDc[1], loginReq, kite.Timeout(time.Second*constant.TimeOut))
	log.Infow("HawkeyeLogin Rsp", "loginRsp", rsp, "err", err)
	if err != nil {
		log.Errorw("HawkeyeLogin Service Error let pass error 2", "err", rsp, "err", err)
		return rsp, nil
	} else {
		if !rsp.Pass {
			return rsp, fmt.Errorf("CheckLogin HawkeyeLogin : %v", rsp)
		} else {
			return rsp, nil
		}
	}
	return rsp, nil
}

func CheckRegister(req internal.Req) (interface{}, error) {
	if !config.Field("hawkeye_open").Bool() {
		return nil, nil
	}
	if channelFilter(req.Game.Channel) {
		return nil, nil
	}

	registerReq := &pbgo.RegisterRequest{
		GameCd:  req.Game.GameCd,
		Subject: pbgo.RegisterRequest_Account,
		DeviceInfo: &pbgo.DeviceInfo{
			SmId:        req.Client.Dhid,
			Ip:          req.IP,
			AppsflyerId: req.Game.AppsflyerId,
		},
		UserInfo: &pbgo.UserInfo{
			BundleId: req.Game.BundleId,
		},
	}
	cfgDc := strings.Split(config.Field(constant.HawkEyeDcCluster).String(), "|")
	log.Infow("HawkeyeRegister Req", "Req", registerReq)
	rsp, err := api.Register(cfgDc[0], cfgDc[1], registerReq, kite.Timeout(time.Second*constant.TimeOut))
	log.Infow("HawkeyeRegister Rsp", "Rsp", rsp, "err", err)
	if err != nil {
		log.Infow("HawkeyeRegister Rsp", "Rsp", rsp, "Err", err)
		return nil, nil
	} else {
		if !rsp.Pass {
			return nil, fmt.Errorf("hawkeye_login HawkeyeLogin : %v", rsp)
		} else {
			return nil, nil
		}
	}
	return nil, nil
}
