package constant

const (
	ErrCodeOk = 0

	ErrCodeDB                  = -1002     //Login内部错误
	ErrGLoginBindType          = -20010150 //绑定类型错误
	ErrGLoginSmsInterval       = -20010200 //60s内只能发送一条短信
	ErrGLoginSmsCount          = -20010201 //10分钟内只能发送3条
	ErrGLoginSmsFail           = -20010202 //短信服务失败
	ErrGLoginThirdAuthFail     = -20010401 //第三方验证失败
	ErrGLoginPlatWrong         = -20010402 //平台标示错误
	ErrCodeParamError          = -20010403 //参数错误
	ErrCodeFastTokenError      = -20010404 //fast参数token错误
	ErrCodeFastTokenExpired    = -20010405 //token过期
	ErrCodeFastTokenVaild      = -20010406 //token错误无效
	ErrCodeSMSGetVerifyFail    = -20010407 //获得验证码错误
	ErrCodeSMSCheckVerifyFail  = -20010408 //验证码检查错误
	ErrCodeCreateInternal      = -20010409 //Create内部错误
	ErrCodeLoginInternal       = -20010410 //Login内部错误
	ErrCodeParsePbInternal     = -20010411 //内部pb解析错误
	ErrCodeCreateVisitorIdFail = -20010412 //创建visitorID错误
)

const (
	ErrMsgOk = "ok"
)
