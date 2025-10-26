package vo

// UserVO 脱敏后的用户数据
type UserVO struct {
	Id          int64  `json:"id"`
	UserAccount string `json:"user_account"`
	Username    string `json:"user_name"`
	UserRole    string `json:"user_role"`
	CreatedAt   int64  `json:"created_at"`
}
