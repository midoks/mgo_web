package form

type LogClean struct {
	Clean string `form:"clean"`
	Day   int64  `form:"day"`
}

type LogSetting struct {
	AllowedManualDelete bool   `json:"allowed_manual_delete"`  // allowed manual delete
	AllowedManual       bool   `json:"allowed_manual"`         // allowed manual
	SaveDay             int64  `json:"save_day"`               // save day
	MaxCapacityLimit    string `json:"max_capacity_limit"`     // max capacity limit
	AllowModClearConfig bool   `json:"allow_mod_clear_config"` // allow mod clear config
}
