package moss

import (
	"fmt"
	"glogin/pbs/glogin"
	"testing"
)

func TestGmp(t *testing.T) {
	//request := glogin.QueryRequest{
	//	//Account:  "notspecified",
	//	Account:  "784732974",
	//	PageNum:  1,
	//	PageSize: 100,
	//	//LoginType: "phone",
	//}
	gmp := Gmp{}
	//rsp, err := gmp.LoadAccountInfo(&request)
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(rsp)
	//}
	//accounts := []string{"627131848", "308240456", "118847550", "309408419", "311090152"}
	////accounts := []string{"627131848"}
	//request2 := glogin.QueryReq{
	//	Accounts: accounts,
	//	PageNum:  1,
	//	PageSize: 100,
	//	//LoginType: "phone",
	//}
	//rsp2, err2 := gmp.QueryAccount(&request2)
	//if err2 != nil {
	//	fmt.Println(err2)
	//} else {
	//	fmt.Println(rsp2)
	//}

	account := "553114659"
	//accounts := []string{"627131848"}
	request2 := &glogin.ChangeBindReq{
		Account: account,
		Phone:   "2234567890",
		Plat:    "visitor",
	}
	rsp2, err2 := gmp.ChangeBind(request2)
	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println(rsp2)
	}
}
