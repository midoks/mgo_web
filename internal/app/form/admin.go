package form

type AdminAdd struct {
	ID         int64             `form:"id"`
	Username   string            `form:"username"`
	Password   string            `form:"password"`
	Password2  string            `form:"password2"`
	FullName   string            `form:"full_name"`
	AllowLogin string            `form:"allow_login"`
	SuperAdmin string            `form:"super_admin"`
	Auth       map[string]string `form:"auth"`
}
