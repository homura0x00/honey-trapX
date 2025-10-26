package dto

// UserRegisterReq 用户注册参数
type UserRegisterReq struct {
	UserAccount   string `json:"user_account"`
	UserPassword  string `json:"user_password"`
	UserName      string `json:"user_name"`
	CheckPassword string `json:"check_password"`
}

type LoginUserReq struct {
	UserAccount  string `json:"user_account"`
	UserPassword string `json:"user_password"`
}

type UpdateUserReq struct {
	UserId       int64  `json:"user_id"`
	UserAccount  string `json:"user_account"`
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_password"`
	UserRole     string `json:"user_role"`
}
