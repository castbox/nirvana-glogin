package bilog

import (
	"encoding/json"
	"glogin/pbs/glogin"
	"strconv"
	"time"
)

type DeviceInfo struct {
	Adid        string `protobuf:"bytes,1,opt,name=adid,proto3" json:"adid,omitempty"`
	Idfv        string `protobuf:"bytes,2,opt,name=idfv,proto3" json:"idfv,omitempty"`
	SmId        string `protobuf:"bytes,3,opt,name=sm_id,json=smId,proto3" json:"sm_id,omitempty"`
	Imei        string `protobuf:"bytes,4,opt,name=imei,proto3" json:"imei,omitempty"`
	AndroidId   string `protobuf:"bytes,5,opt,name=android_id,json=androidId,proto3" json:"android_id,omitempty"`
	AppsflyerId string `protobuf:"bytes,6,opt,name=appsflyer_id,json=appsflyerId,proto3" json:"appsflyer_id,omitempty"`
	DeviceToken string `protobuf:"bytes,8,opt,name=device_token,json=device_token,proto3" json:"device_token,omitempty"`
	MacAddress  string `protobuf:"bytes,8,opt,name=mac_address,json=mac_address,proto3" json:"mac_address,omitempty"`
	DeviceModel string `protobuf:"bytes,8,opt,name=device_model,json=device_model,proto3" json:"device_model,omitempty"`
	DeviceName  string `protobuf:"bytes,8,opt,name=device_name,json=device_name,proto3" json:"device_name,omitempty"`
	OsVersion   string `protobuf:"bytes,8,opt,name=os_version,json=os_version,proto3" json:"os_version,omitempty"`
	Language    string `protobuf:"bytes,8,opt,name=language,json=language,proto3" json:"language,omitempty"`
	NetworkType string `protobuf:"bytes,8,opt,name=network_type,json=network_type,proto3" json:"network_type,omitempty"`
	AppVersion  string `protobuf:"bytes,8,opt,name=app_version,json=app_version,proto3" json:"app_version,omitempty"`
	Ip          string `protobuf:"bytes,7,opt,name=ip,proto3" json:"ip,omitempty"`
	Oaid        string `protobuf:"bytes,8,opt,name=oaid,json=oaid,proto3" json:"oaid,omitempty"`
}

type UserInfo struct {
	BundleId string `protobuf:"bytes,1,opt,name=bundle_id,json=bundleId,proto3" json:"bundle_id,omitempty"`
	ServerId string `protobuf:"bytes,2,opt,name=server_id,json=serverId,proto3" json:"server_id,omitempty"`
	UserId   string `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Account  string `protobuf:"bytes,4,opt,name=account,proto3" json:"account,omitempty"`
	Lv       string `protobuf:"bytes,5,opt,name=lv,proto3" json:"lv,omitempty"`
	Vip      string `protobuf:"bytes,6,opt,name=vip,proto3" json:"vip,omitempty"`
}

type EventValueLogin struct {
	DeviceInfo *DeviceInfo `json:"device_info"`
	UserInfo   *UserInfo   `json:"user_info"`
	EventInfo  interface{} `json:"event_info"`
}

type LoginEvent struct {
	LoginType string `json:"login_type"`
}

func SmsLogin(req *glogin.SmsLoginReq, dhAccount string) {
	login(LoginTypeSMS, dhAccount, req.Game, req.Client)
}

func FastLogin(req *glogin.FastLoginReq, dhAccount string) {
	login(LoginTypeFast, dhAccount, req.Game, req.Client)
}

func ThirdLogin(req *glogin.ThirdLoginReq, dhAccount string) {
	loginType := GetLoginType(req.ThirdPlat)
	login(loginType, dhAccount, req.Game, req.Client)
}

func GetLoginType(thirdPlat string) string {
	if thirdPlat == "" {
		return LoginTypeNone
	}
	return thirdPlat
}

func login(loginType string, dhAccount string, reqGame *glogin.LoginGame, reqClient *glogin.LoginClient) {
	value := EventValueLogin{
		DeviceInfo: &DeviceInfo{
			Adid:        reqGame.Adid,
			Idfv:        reqGame.Idfv,
			SmId:        reqClient.Dhid,
			Imei:        reqClient.Imei,
			AndroidId:   reqClient.AndroidId,
			AppsflyerId: reqGame.AppsflyerId,
			DeviceToken: reqGame.DeviceToken,
			MacAddress:  reqClient.MacAddress,
			DeviceModel: reqClient.DeviceModel,
			DeviceName:  reqClient.DviceName,
			OsVersion:   reqClient.OsVersion,
			Language:    reqGame.Language,
			NetworkType: reqClient.NetworkType,
			AppVersion:  reqGame.AppVersion,
			Ip:          reqClient.Ip,
			Oaid:        "",
		},
		UserInfo: &UserInfo{
			BundleId: reqGame.BundleId,
			Account:  dhAccount,
		},
		EventInfo: &LoginEvent{
			loginType,
		}}
	data, _ := json.Marshal(value)
	log := Log{
		EventType:  EventTypeLogin,
		EventCode:  EventCodeLogin,
		EventName:  EventNameLogin,
		GameCd:     reqGame.GameCd,
		CreateTs:   strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		EventValue: string(data),
	}
	go log.Push()
}
