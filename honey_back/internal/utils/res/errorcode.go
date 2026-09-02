package res

// 统一业务错误码。约定：0 成功；其余非 0。
const (
	CodeOK       = 0
	ParamCode    = 40100
	UserNotLogin = 40200
	BadLogin     = 40201
	Permission   = 40300
	NotFound     = 40400
	SystemError  = 50000
)
