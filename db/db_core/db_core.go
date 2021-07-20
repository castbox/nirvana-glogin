package db_core

type AccountData struct {
	ID        int32      `json:"id" bson:"_id"`                // id DH_account
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
