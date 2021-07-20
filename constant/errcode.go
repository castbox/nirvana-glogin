package constant

const (
	ErrCodeOk                  = 0
	ErrCodePlatIsWrong         = 1
	ErrCodeCreateAccount       = 2
	ErrCodeInternal            = 3
	ErrCodeSMIDError           = 4
	ErrCodeTokenError          = 9
	ErrCodeFastTokenExpired    = 10
	ErrCodeFastTokenVaild      = 11
	ErrCodeSMSGetVerifyFaild   = 31
	ErrCodeSMSCheckVerifyFaild = 32

	// 兼容erlang版本的返回值
	ErrGLoginVerifyCode  = -20010113
	ErrGLoginBindType    = -20010150
	ErrGLoginSmsInterval = -20010200
	ErrGLoginSmsCount    = -20010201
	ErrGLoginSmsFail     = -20010202
)

const (
	ErrMsgOk = "ok"
)
