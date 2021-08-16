package util

import (
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

// 字符串转Int
//    intStr：数字的字符串
func String2Int(intStr string) (intNum int) {
	intNum, _ = strconv.Atoi(intStr)
	return
}

// 字符串转Int64
//    intStr：数字的字符串
func String2Int64(intStr string) (int64Num int64) {
	intNum, _ := strconv.Atoi(intStr)
	int64Num = int64(intNum)
	return
}

// 字符串转Float64
//    floatStr：小数点数字的字符串
func String2Float64(floatStr string) (floatNum float64) {
	floatNum, _ = strconv.ParseFloat(floatStr, 64)
	return
}

// 字符串转Float32
//    floatStr：小数点数字的字符串
func String2Float32(floatStr string) (floatNum float32) {
	floatNum64, _ := strconv.ParseFloat(floatStr, 32)
	floatNum = float32(floatNum64)
	return
}

// Int转字符串
//    intNum：数字字符串
func Int2String(intNum int) (intStr string) {
	intStr = strconv.Itoa(intNum)
	return
}

// Int64转字符串
//    intNum：数字字符串
func Int642String(intNum int64) (int64Str string) {
	//10, 代表10进制
	int64Str = strconv.FormatInt(intNum, 10)
	return
}

// Float64转字符串
//    floatNum：float64数字
//    prec：精度位数（不传则默认float数字精度）
func Float64ToString(floatNum float64, prec ...int) (floatStr string) {
	if len(prec) > 0 {
		floatStr = strconv.FormatFloat(floatNum, 'f', prec[0], 64)
		return
	}
	floatStr = strconv.FormatFloat(floatNum, 'f', -1, 64)
	return
}

// Float32转字符串
//    floatNum：float32数字
//    prec：精度位数（不传则默认float数字精度）
func Float32ToString(floatNum float32, prec ...int) (floatStr string) {
	if len(prec) > 0 {
		floatStr = strconv.FormatFloat(float64(floatNum), 'f', prec[0], 32)
		return
	}
	floatStr = strconv.FormatFloat(float64(floatNum), 'f', -1, 32)
	return
}

// 二进制转10进制
func BinaryToDecimal(bit string) (num int) {
	fields := strings.Split(bit, "")
	lens := len(fields)
	var tempF float64 = 0
	for i := 0; i < lens; i++ {
		floatNum := String2Float64(fields[i])
		tempF += floatNum * math.Pow(2, float64(lens-i-1))
	}
	num = int(tempF)
	return
}

// BytesToString 0 拷贝转换 slice byte 为 string
func BytesToString(b []byte) (s string) {
	_bptr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	_sptr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	_sptr.Data = _bptr.Data
	_sptr.Len = _bptr.Len
	return s
}

func IntToByteArray(num int32) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

func ByteArrayToInt(arr []byte) int32 {
	val := int32(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

// 老数据BsonAToStr兼容
func BsonAToStr(src interface{}) string {
	if value, ok := src.(string); ok {
		return value
	} else if value, ok := src.(bson.A); ok {
		var bs []byte
		for i := 0; i < len(value); i++ {
			v := value[i].(int32)
			bs = append(bs, byte(v))
		}
		return bytes2str(bs)
	} else if value, ok := src.([]interface{}); ok {
		var bs []byte
		for i := 0; i < len(value); i++ {
			if v, ok := value[i].(int32); ok {
				bs = append(bs, byte(v))
			} else if v, ok := value[i].(uint8); ok {
				bs = append(bs, v)
			} else {
				return ""
			}
		}
		return bytes2str(bs)
	} else {
		return ""
	}
}

func IntArrayToInterfaceArray(src []int32) []interface{} {
	var rstArray []interface{}
	for i := 0; i < len(src); i++ {
		rstArray = append(rstArray, interface{}(src[i]))
	}
	return rstArray
}

func StrToInterfaceArray(src string) []interface{} {
	bs := str2bytes(src)
	var rstArray []interface{}
	for i := 0; i < len(bs); i++ {
		rstArray = append(rstArray, interface{}(bs[i]))
	}
	return rstArray
}

func StrToIntArray(src string) []int32 {
	bs := str2bytes(src)
	var rstArray []int32
	for i := 0; i < len(bs); i++ {
		rstArray = append(rstArray, int32(bs[i]))
	}
	return rstArray
}

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func bytes2str(b []byte) string {
	//strValue := string(bs[:])
	return *(*string)(unsafe.Pointer(&b))
}
