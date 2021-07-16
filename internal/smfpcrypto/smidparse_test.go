package smfpcrypto

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("args != 3")
	}

	boxId := os.Args[1]
	pri := os.Args[2]

	smid, err := ParseBoxData(boxId, pri)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(smid)
	}
}
