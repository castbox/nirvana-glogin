package anti

import (
	"git.dhgames.cn/svr_comm/anti_obsession/api"
	"git.dhgames.cn/svr_comm/anti_obsession/pbs/pb_obsession"
	log "git.dhgames.cn/svr_comm/gcore/glog"
	"git.dhgames.cn/svr_comm/kite"
	"glogin/config"
	"glogin/constant"
	"glogin/internal"
	"strings"
	"time"
)

func Check(req *pb_obsession.CheckRequest) (*pb_obsession.CheckResponse, error) {
	autiDcCluster := config.Field(constant.AutiDcCluster).String()
	cfgDc := strings.Split(autiDcCluster, "|")
	rsp, err := api.Check(cfgDc[0], cfgDc[1], req, kite.Timeout(time.Second*constant.TimeOut))
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
	queryIn := &pb_obsession.CheckStateQueryRequest{
		GameCd: req.GameCd,
		Id:     req.Account,
		//DeviceId: req.Client.Dhid,
	}
	log.Infow("api.CheckAuditQuery req", "queryIn", queryIn)
	cfgDc := strings.Split(config.Field(constant.AutiDcCluster).String(), "|")
	rsp, err := api.CheckAuditQuery(cfgDc[0], cfgDc[1], queryIn, kite.Timeout(time.Second*constant.TimeOut))
	log.Infow("api.CheckAuditQuery Rsp", "rsp", rsp, "err", err)
	if err != nil {
		return rsp, err
	} else {
		return rsp, nil
	}
	return nil, nil
}
