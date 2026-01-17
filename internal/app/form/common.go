package form

type Page struct {
	Name string `form:"name"`
}

type ID struct {
	ID int64 `form:"id"`
}
