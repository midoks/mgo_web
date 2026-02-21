package form

type LogClean struct {
	Clean string `form:"clean"`
	Day   int64  `form:"day"`
}
