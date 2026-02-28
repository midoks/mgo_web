package form

type SettingProfile struct {
	Name string `form:"name"`
}

type SettingLogin struct {
	Name      string `form:"name"`
	Password  string `form:"password"`
	Password2 string `form:"password2"`
}
