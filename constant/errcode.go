package constant

const (
	ErrCodeOk                  = 0
	ErrCodeDB                  = -1002     //Login内部Db错误
	ErrUnkown                  = -1004     //ErrUnkown
	ErrCodeBindType            = -20010150 //绑定类型错误
	ErrCodeSmsInterval         = -20010200 //60s内只能发送一条短信
	ErrCodeSmsCount            = -20010201 //10分钟内只能发送3条
	ErrCodeSmsFail             = -20010400 //短信服务失败
	ErrCodeThirdAuthFail       = -20010401 //第三方验证失败
	ErrCodePlatWrong           = -20010402 //平台标示错误
	ErrCodeParamError          = -20010403 //请求参数错误
	ErrCodeFastTokenError      = -20010404 //fast token错误
	ErrCodeFastTokenExpired    = -20010405 //fast token过期
	ErrCodeFastTokenVaild      = -20010406 //token错误无效
	ErrCodeSMSGetVerifyFail    = -20010407 //获得验证码错误
	ErrCodeSMSCheckVerifyFail  = -20010408 //验证码检查错误
	ErrCodeCreateInternal      = -20010409 //Create内部错误
	ErrCodeLoginInternal       = -20010410 //Login内部错误
	ErrCodeParsePbInternal     = -20010411 //内部pb解析错误
	ErrCodeCreateVisitorIdFail = -20010412 //创建visitorID错误
	ErrCodeBindVisitorNotExist = -20010413 //绑定游客账号不存在
	ErrCodeVisitorLoadErr      = -20010414 //游客账号加载失败
	ErrCodeThirdAlreadyBind    = -20010415 //第三方账号已经绑定
	ErrCodeThirdBindFail       = -20010416 //第三方账号绑定失败
	ErrCodePhoneAlreadyBind    = -20010417 //phone已经绑定
	ErrCodeNoAccount           = -20010418 //account不存在
	ErrCodeUpdateFail          = -20010419 //更新token失败
	ErrCodeNoThirdErr          = -20010420 //第三方错误

)

const (
	ErrMsgOk          = "success"
	ErrCodeStrOk      = "0"
	ErrCodeStrAutiRpc = "auti service rpc error"
	ErrMsgUnkwon      = "unkown"
)
