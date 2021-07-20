package sms

//{
//"code": "SUCCESS",
//"data": [
//{
//"mobile": "15538850000",
//"smsId": "20170392833833891100",
//"customSmsId": "20170392833833891100"
//},
//{
//"mobile": "15538850001",
//"smsId": "20170392833833892100",
//"customSmsId": "20170392833833891100"
//}
//]
//}
// 发送短信返回
type SendSmsRsp struct {
	Code string        `json:"code" binding:"required"`
	Data []interface{} `json:"data" binding:"required"`
}
