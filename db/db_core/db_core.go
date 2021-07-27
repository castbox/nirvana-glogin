package db_core

// account table
type AccountData struct {
	ID        int32      `json:"_id" bson:"_id"`               // id DH_account
	BundleID  string     `json:"bundle_id" bson:"bundle_id"`   // 包名
	Create    CreateData `json:"create" bson:"create"`         // create创建信息
	Google    string     `json:"google" bson:"google"`         // google
	Facebook  string     `json:"facebook" bson:"facebook"`     // facebook unionId
	IOS       string     `json:"ios" bson:"ios"`               // IOS
	Visitor   string     `json:"visitor" bson:"visitor"`       // visitor
	Phone     string     `json:"phone" bson:"phone"`           // phone
	WeChat    string     `json:"we_chat" bson:"we_chat"`       // wechat
	LastLogin int64      `json:"last_login" bson:"last_login"` // last_login
}

type CreateData struct {
	Ip       string `json:"ip" bson:"ip"`               // Ip
	Time     int64  `json:"time" bson:"time"`           // Time
	SmId     string `json:"sm_id" bson:"sm_id"`         // SmId
	BundleId string `json:"bundle_id" bson:"bundle_id"` // BundleId
}

type VerifyCodeData struct {
	ID         string      `json:"id" bson:"_id"`                  // _id
	SendTime   int64       `json:"send_time" bson:"send_time"`     // send_time
	Expire     interface{} `json:"expire" bson:"expire"`           // expire
	Phone      string      `json:"phone" bson:"phone"`             // phone
	VerifyCode string      `json:"verify_code" bson:"verify_code"` // verify_code
}
