package form

type Page struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type ID struct {
	ID int64 `form:"id"`
}
