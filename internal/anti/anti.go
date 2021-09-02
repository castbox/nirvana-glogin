package anti

import (
	"git.dhgames.cn/svr_comm/gmoss/v3"
	"git.dhgames.cn/svr_comm/gmoss/v3/global"
	"glogin/config"
	"glogin/constant"
	"glogin/internal"
	anti_authentication "glogin/pbs/authentication"
	"strings"
	"time"
)

func Check(req *anti_authentication.CheckRequest) (*anti_authentication.CheckResponse, error) {
	autiDcCluster := config.Field(constant.AutiDcCluster).String()
	cfgDc := strings.Split(autiDcCluster, "|")
	rsp, err := anti_authentication.Check(req, global.WithCluster(cfgDc[0], cfgDc[1], constant.AutiService).WithTimeout(time.Second*9))
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
	cfgDc := strings.Split(config.Field(constant.AutiDcCluster).String(), "|")
	service := gmoss.MossWithDcClusterService(cfgDc[0], cfgDc[1], constant.AutiService)
	rsp, err := anti_authentication.AuditQuery(queryIn, &global.CallOption{
		Cluster: service,
	})
	if err != nil {
		return rsp, err
	} else {
		return rsp, nil
	}
	return nil, nil
}
