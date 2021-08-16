package db_core

// account table
type AccountData struct {
	ID        int32       `json:"_id" bson:"_id"`               // id DH_account
	BundleID  string      `json:"bundle_id" bson:"bundle_id"`   // 包名
	Create    CreateData  `json:"create" bson:"create"`         // create创建信息
	Google    string      `json:"google" bson:"google"`         // google
	Facebook  string      `json:"facebook" bson:"facebook"`     // facebook unionId
	IOS       interface{} `json:"ios" bson:"ios"`               // IOS 调整为interface兼容老数据 (老array 新string)
	Visitor   string      `json:"visitor" bson:"visitor"`       // visitor
	Phone     string      `json:"phone" bson:"phone"`           // phone
	WeChat    string      `json:"we_chat" bson:"we_chat"`       // wechat
	LastLogin int64       `json:"last_login" bson:"last_login"` // last_login
}

type CreateData struct {
	Ip       interface{} `json:"ip" bson:"ip"`               // Ip
	Time     int64       `json:"time" bson:"time"`           // Time
	SmId     string      `json:"sm_id" bson:"sm_id"`         // SmId
	BundleId string      `json:"bundle_id" bson:"bundle_id"` // BundleId
}

// VerifyCodeData
type VerifyCodeData struct {
	ID         string      `json:"id" bson:"_id"`                  // _id
	SendTime   int64       `json:"send_time" bson:"send_time"`     // send_time
	Expire     interface{} `json:"expire" bson:"expire"`           // expire
	Phone      string      `json:"phone" bson:"phone"`             // phone
	VerifyCode string      `json:"verify_code" bson:"verify_code"` // verify_code
}

// TokenForBusinessData
//{
//"_id": ObjectId("60f02a59c68c824948000001"),
//"bundle_id": "com.droidhang.aod.google",
//"facebook_token": "191614856085993",
//"token_for_business": "AbyrrMHgDlKgdieY"
//}
type TokenForBusinessData struct {
	ID               interface{} `json:"_id" bson:"_id"`                               // id 60f02a59c68c824948000001
	BundleID         string      `json:"bundle_id" bson:"bundle_id"`                   // 包名
	FacebookToken    string      `json:"facebook_token" bson:"facebook_token"`         // facebook_token
	TokenForBusiness string      `json:"token_for_business" bson:"token_for_business"` // token_for_business
}
