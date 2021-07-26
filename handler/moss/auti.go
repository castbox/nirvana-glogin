package moss

import anti_authentication "glogin/pbs/authentication"

type Anti struct {
}

func (Anti) Query(request *anti_authentication.StateQueryRequest) (response *anti_authentication.StateQueryResponse, err error) {
	return response, nil
}

func (Anti) Check(request *anti_authentication.CheckRequest) (response *anti_authentication.CheckResponse, err error) {
	response = &anti_authentication.CheckResponse{}
	return response, nil
}
