package res

type ErrCode int

// TODO 完善错误状态码
const (
	ParamCode    = 40100
	UserNotLogin = 40200
	Permission   = 40300
	SystemError  = 50000
)

func (err ErrCode) Error() string {
	switch err {
	case ParamCode:
		return "Parameter Code Error"
	case UserNotLogin:
		return "User Not Login"
	case Permission:
		return "Permission Error"
	case SystemError:
		return "System Error"
	default:
		return "Unknown Error"
	}
}
