package form

type AdminAdd struct {
	ID        string `form:"id"`
	Username  string `form:"username"`
	Password  string `form:"password"`
	Password2 string `form:"password2"`
	FullName  string `form:"full_name"`
}
