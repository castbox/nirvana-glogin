package smfpcrypto

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("args != 3")
	}

	//"parse_smid_private_key": "MIIBPAIBAAJBAOfPLQ993UR8qJoCVJsj00/BcPDbKIjEDYnqMjgUAiQkMgYf9O4L7+WNhhtA+kIllsHpEAYJuSdl04wP05Pk0TkCAwEAAQJBAMzdJOafBrDjNqI9UwZ0x+ihfa3vEcik844iItW6oRXMFIo+P+6YHjgiiyLXeSu+60WQ4IfWdRZNdbHMhr1IIN0CIQD6GzKUls0YXxASUmdcSTUFeqXcedkhcLafHTk8jqcX8wIhAO1FmW+cyx1gm4msyhgXN1Fb7frFHniaP5L89zc6NwkjAiEAmA0e5A0GJVHt8GWepxFupaUZ3v9JDTZ8ICHhITrMxRcCIFI92Z0yP8UDA2aJGdOX2Hi+4JIXWSR8cqTEQfxGlWT5AiEA5TqIC6znNIGzeAeuz3Hdj4srmAEP/VG9EkDvdgMT6Tg=",
	//"parse_smid_private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIBPAIBAAJBAOfPLQ993UR8qJoCVJsj00/BcPDbKIjEDYnqMjgUAiQkMgYf9O4L7+WNhhtA+kIllsHpEAYJuSdl04wP05Pk0TkCAwEAAQJBAMzdJOafBrDjNqI9UwZ0x+ihfa3vEcik844iItW6oRXMFIo+P+6YHjgiiyLXeSu+60WQ4IfWdRZNdbHMhr1IIN0CIQD6GzKUls0YXxASUmdcSTUFeqXcedkhcLafHTk8jqcX8wIhAO1FmW+cyx1gm4msyhgXN1Fb7frFHniaP5L89zc6NwkjAiEAmA0e5A0GJVHt8GWepxFupaUZ3v9JDTZ8ICHhITrMxRcCIFI92Z0yP8UDA2aJGdOX2Hi+4JIXWSR8cqTEQfxGlWT5AiEA5TqIC6znNIGzeAeuz3Hdj4srmAEP/VG9EkDvdgMT6Tg=\n-----END RSA PRIVATE KEY-----",

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
