package moss

import (
	"fmt"
	"glogin/pbs/glogin"
	"testing"
)

func TestGmp(t *testing.T) {
	request := glogin.QueryRequest{
		Account:   "notspecified",
		LoginType: "phone",
	}
	gmp := Gmp{}
	rsp, err := gmp.LoadAccountInfo(&request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rsp)
	}
}
