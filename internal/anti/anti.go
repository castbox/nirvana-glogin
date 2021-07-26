package anti

import (
	"git.dhgames.cn/svr_comm/gmoss/v2"
	"glogin/config"
	"glogin/constant"
	"glogin/internal"
	anti_authentication "glogin/pbs/authentication"
	"strings"
)

func Check(req *anti_authentication.CheckRequest) (*anti_authentication.CheckResponse, error) {
	autiDc := config.Field(constant.AutiDc).String()
	cfgDc := strings.Split(autiDc, "|")
	service := gmoss.MossWithDcClusterService(cfgDc[0], cfgDc[1], constant.AutiService)
	rsp, err := anti_authentication.Check(service, req, gmoss.Call, gmoss.DefaultCallOption())
	if err != nil {
		return rsp, err
	} else {
		return rsp, nil
	}
	return nil, nil
}

func StateQuery(req internal.Req) (interface{}, error) {
	if !config.Field("auti_open").Bool() {
		return nil, nil
	}

	if req.GameCd == "" {
		req.GameCd = req.Game.GameCd
	}
	queryIn := &anti_authentication.StateQueryRequest{
		GameCd: req.GameCd,
		Id:     req.Account,
		//DeviceId: req.Client.Dhid,
	}
	cfgDc := strings.Split(config.Field(constant.AutiDc).String(), "|")
	service := gmoss.MossWithDcClusterService(cfgDc[0], cfgDc[1], constant.AutiService)
	rsp, err := anti_authentication.StateQuery(service, queryIn, gmoss.Call, gmoss.DefaultCallOption())
	if err != nil {
		return rsp, err
	} else {
		return rsp, nil
	}
	return nil, nil
}
