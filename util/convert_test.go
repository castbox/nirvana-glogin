package util

import "testing"

func TestBytesToString(t *testing.T) {
	//var b []byte = []byte("abc")
	//s := BytesToString(b)
	//if s != string(b) {
	//	t.Fatal("BytesToString error")
	//}
	//b[1] = 'c'
	//if s != "acc" {
	//	t.Fatal("BytesToString error")
	//}

	//
	var IOSString = "001286.b6c90ca4961f4dc39791e2e63ea9d134.0148_cnofficial"
	t.Logf("src %v", IOSString)
	intArray := StrToIntArray(IOSString)
	t.Logf("StrToIntArray %v", intArray)

	bsonArray := StrToInterfaceArray(IOSString)
	t.Logf("StrToInterfaceArray %v", bsonArray)

	strConvert1 := BsonAToStr(bsonArray)
	t.Logf("IntArrayToStr %v", strConvert1)
}
