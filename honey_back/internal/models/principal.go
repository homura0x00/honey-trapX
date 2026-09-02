package models

// CtxPrincipal 中间件写入 gin.Context 的键名
const CtxPrincipal = "principal"

// Principal 登录态用户（无敏感字段），auth 写入、各 feature 消费，避免跨包耦合
type Principal struct {
	ID          int64  `json:"id"`
	UserAccount string `json:"user_account"`
	UserName    string `json:"user_name"`
	UserRole    string `json:"user_role"`
}
