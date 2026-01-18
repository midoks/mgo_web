package form

type Page struct {
	Page  int64 `form:"page"`
	Limit int64 `form:"limit"`
}

type ID struct {
	ID int64 `form:"id"`
}
