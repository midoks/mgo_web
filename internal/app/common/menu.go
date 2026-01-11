package common

import (
	// "fmt"

	"encoding/json"
	"mgo/embed"
	"sync"
)

type MenuConf struct {
	Code     string     `json:"code"`
	Name     string     `json:"name"`
	Icon     string     `json:"icon"`
	Path     string     `json:"path"`
	Perm     string     `json:"perm"`
	Children []MenuConf `json:"children,omitempty"`
}

var (
	menus    []MenuConf
	menuOnce sync.Once
)

func GetMenus() []MenuConf {
	menuOnce.Do(func() {
		content, err := embed.Conf.ReadFile("conf/menu.json")
		if err != nil {
			return
		}
		_ = json.Unmarshal(content, &menus)
	})
	return menus
}

func FindMenuCodeByPath(requestPath string, adminPath string) string {
	ms := GetMenus()
	return findMenuCodeRecursive(ms, requestPath, adminPath)
}

func findMenuCodeRecursive(menus []MenuConf, requestPath string, adminPath string) string {
	for _, m := range menus {
		// Construct full path for comparison
		// If adminPath is not empty, prepend it.
		// requestPath usually starts with /
		// m.Path usually starts with /
		// e.g. adminPath="admin", m.Path="/index" -> "/admin/index"
		// e.g. adminPath="", m.Path="/index" -> "/index"

		fullPath := ""
		if adminPath != "" {
			fullPath = "/" + adminPath + m.Path
		} else {
			fullPath = m.Path
		}

		// fmt.Println(m.Path, fullPath, requestPath)
		// fmt.Println(m.Children)
		// Check exact match
		if m.Path != "" && fullPath == requestPath {
			return m.Code
		}

		// Recursive check children
		if len(m.Children) > 0 {
			code := findMenuCodeRecursive(m.Children, requestPath, adminPath)
			if code != "" {
				return code
			}
		}
	}
	return ""
}
