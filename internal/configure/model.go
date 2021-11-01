package configure

//公共配置请求
type CfgRequest struct {
	AppName string `json:"app_name" binding:"required"`
	Vsn     string `json:"vsn"`
}

// 地址配置
type ClusterCfg struct {
	ClusterType      int32  `json:"cluster_type" binding:"required"`
	UGateAddr        string `json:"ugate_addr" binding:"required"`
	VsnAddr          string `json:"vsn_addr" binding:"required"`
	UPayAddr         string `json:"upay_addr" binding:"required"`
	UChatAddr        string `json:"uchat_addr" binding:"required"`
	UChatWSAddr      string `json:"uchat_ws_addr" binding:"required"`
	CommunityWebAddr string `json:"community_web_addr" binding:"required"`
	CommunitySrvAddr string `json:"community_srv_addr" binding:"required"`
	AicsWsAddr       string `json:"aics_ws_addr" binding:"required"`
	AicsHttpAddr     string `json:"aics_http_addr" binding:"required"`
}

// 连接配置
type ConnCfg struct {
	Stable ClusterCfg `json:"stable" binding:"required"`
	Check  ClusterCfg `json:"check" binding:"required"`
	Test   ClusterCfg `json:"test"`
}

// pub_cfg consul配置
type PubCfg struct {
	ConnCfg   ConnCfg  `json:"conn_cfg" binding:"required"`
	CheckVsn  string   `json:"check_vsn"`
	WhiteList []string `json:"white_list"`
}

// 加密的pub返回结构
type PubCfgValue struct {
	LockIndex   int32
	Key         string
	Flags       int32
	Value       string
	CreateIndex int64
	ModifyIndex int64
}

type PubCfgList struct {
	AppList []PubCfgValue
}

type BundleInfo struct {
	FacebookOauthUrl        string `json:"facebook_oauth_url" binding:"required"`
	AppsFlyerANDROID        string `json:"appsflyer_ANDROID" binding:"required"`
	AppsFlyerOpen           bool   `json:"appsflyer_open" binding:"required"`
	AppsFlyerIOS            string `json:"appsflyer_IOS" binding:"required"`
	AppsFlyerAuthentication string `json:"appsflyer_Authentication" binding:"required"`
	AppsflyerRegistrationId int32  `json:"appsflyer_registrationId" binding:"required"`
}
