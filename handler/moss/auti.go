package moss

import (
	"git.dhgames.cn/svr_comm/anti_obsession/pbs/pb_obsession"
)

type Anti struct {
}

func (Anti) Query(request *pb_obsession.CheckStateQueryRequest) (response *pb_obsession.CheckStateQueryResponse, err error) {
	return response, nil
}

func (Anti) Check(request *pb_obsession.CheckRequest) (response *pb_obsession.CheckResponse, err error) {
	response = &pb_obsession.CheckResponse{}
	return response, nil
}
