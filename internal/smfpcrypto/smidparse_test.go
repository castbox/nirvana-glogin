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

	smid, err := ParseBoxId(boxId, pri)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(smid)
	}
}

//fmt.Printf("时间戳（秒）：%v;\n", time.Now().Unix())
//fmt.Printf("时间戳（纳秒）：%v;\n",time.Now().UnixNano())
//fmt.Printf("时间戳（毫秒）：%v;\n",time.Now().UnixNano() / 1e6)
//fmt.Printf("时间戳（纳秒转换为秒）：%v;\n",time.Now().UnixNano() / 1e9)
