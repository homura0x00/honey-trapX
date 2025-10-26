package enums

type UserRoleEnum struct {
	Admin string
	User  string
}

var (
	UserRoleUser  = "USER"
	UserRoleAdmin = "ADMIN"
)
