package conf

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"mgo/embed"
)

var File *yaml.Node

func ReadConf() (*yaml.Node, error) {
	data, err := embed.Conf.ReadFile("conf/app.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "read file 'conf/app.yaml'")
	}

	File = &yaml.Node{}
	err = yaml.Unmarshal(data, File)
	if err != nil {
		return nil, errors.Wrap(err, "parse 'conf/app.yaml'")
	}

	return File, nil
}

func InstallConf(data map[string]string) error {
	cfg := &Config{
		AppName:   App.Name,
		BrandName: App.BrandName,
		RunUser:   App.RunUser,
		RunMode:   "prod",
		General: GeneralConfig{
			MenuFile: "menu.json",
		},
		Log: LogConfig{
			Format:   "text",
			RootPath: "logs",
		},
		Session: SessionConfig{
			Provider:       "memory",
			ProviderConfig: "data/sessions",
			CookieName:     "mgo_web",
			CookieSecure:   false,
			GCInterval:     3600,
			MaxLifeTime:    86400,
			CSRFCookieName: "_csrf",
		},
		Web: WebConfig{
			HTTPPort:  9999,
			AdminPath: "mgo",
		},
		Security: SecurityConfig{
			InstallLock:             true,
			SecretKey:               randString(10),
			LoginRememberDays:       7,
			CookieRememberName:      "mgo_incredible",
			CookieUsername:          "mgo_awesome",
			CookieSecure:            false,
			EnableLoginStatusCookie: true,
			LoginStatusCookieName:   "login_status",
		},
	}

	if strings.EqualFold(data["type"], "mysql") {
		cfg.Database = DatabaseConfig{
			Type:        "mysql",
			Hostname:    data["hostname"],
			Hostport:    int64(parseInt(data["hostport"])),
			Name:        data["dbname"],
			User:        data["username"],
			Password:    data["password"],
			TablePrefix: data["table_prefix"],
		}
	} else if strings.EqualFold(data["type"], "sqlite3") {
		cfg.Database = DatabaseConfig{
			Type:        "sqlite3",
			Path:        data["dbpath"],
			TablePrefix: data["table_prefix"],
		}
	}

	customConf := filepath.Join(CustomDir(), "conf", "app.yaml")

	if !isExist(filepath.Dir(customConf)) {
		os.MkdirAll(filepath.Dir(customConf), os.ModePerm)
	}

	return SaveConfData(cfg, customConf)
}

func InitConf(customConf string) error {
	data, err := embed.Conf.ReadFile("conf/app.yaml")
	if err != nil {
		return errors.Wrap(err, "read embedded config")
	}

	File = &yaml.Node{}
	err = yaml.Unmarshal(data, File)
	if err != nil {
		return errors.Wrap(err, "parse 'conf/app.yaml'")
	}

	if customConf == "" {
		customConf = filepath.Join(CustomDir(), "conf", "app.yaml")
	} else {
		customConf, err = filepath.Abs(customConf)
		if err != nil {
			return errors.Wrap(err, "get absolute path")
		}
	}
	CustomConf = customConf

	if isFile(customConf) {
		customData, err := os.ReadFile(customConf)
		if err != nil {
			return errors.Wrapf(err, "read %q", customConf)
		}
		customNode := &yaml.Node{}
		if err = yaml.Unmarshal(customData, customNode); err != nil {
			return errors.Wrapf(err, "parse %q", customConf)
		}
		mergeYaml(File, customNode)
	}

	err = renderSection(File)
	if err != nil {
		return err
	}

	return nil
}

func renderSection(node *yaml.Node) error {
	m := make(map[string]yaml.Node)
	if err := node.Decode(&m); err != nil {
		return errors.Wrap(err, "decode yaml node")
	}

	if app, ok := m["app_name"]; ok {
		App.Name = app.Value
	}
	if app, ok := m["brand_name"]; ok {
		App.BrandName = app.Value
	}
	if app, ok := m["run_user"]; ok {
		App.RunUser = app.Value
	}
	if app, ok := m["run_mode"]; ok {
		App.RunMode = app.Value
	}

	if general, ok := m["general"]; ok {
		if err := general.Decode(&General); err != nil {
			return errors.Wrap(err, "mapping general section")
		}
	}

	if log, ok := m["log"]; ok {
		if err := log.Decode(&Log); err != nil {
			return errors.Wrap(err, "mapping log section")
		}
	}

	if database, ok := m["database"]; ok {
		if err := database.Decode(&Database); err != nil {
			return errors.Wrap(err, "mapping database section")
		}
	}

	if session, ok := m["session"]; ok {
		if err := session.Decode(&Session); err != nil {
			return errors.Wrap(err, "mapping session section")
		}
	}

	if security, ok := m["security"]; ok {
		if err := security.Decode(&Security); err != nil {
			return errors.Wrap(err, "mapping security section")
		}
	}

	if web, ok := m["web"]; ok {
		if err := web.Decode(&Web); err != nil {
			return errors.Wrap(err, "mapping web section")
		}
	}

	return nil
}

func mergeYaml(base, overlay *yaml.Node) {
}

func SaveConf(node *yaml.Node, path string) error {
	data, err := yaml.Marshal(node)
	if err != nil {
		return errors.Wrap(err, "marshal yaml")
	}
	return os.WriteFile(path, data, 0644)
}

type Config struct {
	AppName   string         `yaml:"app_name"`
	BrandName string         `yaml:"brand_name"`
	RunUser   string         `yaml:"run_user"`
	RunMode   string         `yaml:"run_mode"`
	General   GeneralConfig  `yaml:"general"`
	Log       LogConfig      `yaml:"log"`
	Session   SessionConfig  `yaml:"session"`
	Web       WebConfig      `yaml:"web"`
	Database  DatabaseConfig `yaml:"database"`
	Security  SecurityConfig `yaml:"security"`
}

type GeneralConfig struct {
	MenuFile string `yaml:"menu_file"`
}

type LogConfig struct {
	Format   string `yaml:"format"`
	RootPath string `yaml:"root_path"`
}

type SessionConfig struct {
	Provider       string `yaml:"provider"`
	ProviderConfig string `yaml:"provider_config"`
	CookieName     string `yaml:"cookie_name"`
	CookieSecure   bool   `yaml:"cookie_secure"`
	GCInterval     int64  `yaml:"gc_interval"`
	MaxLifeTime    int64  `yaml:"max_life_time"`
	CSRFCookieName string `yaml:"csrf_cookie_name"`
}

type WebConfig struct {
	HTTPAddr   string `yaml:"http_addr"`
	HTTPPort   int    `yaml:"http_port"`
	AdminPath  string `yaml:"admin_path"`
	EnableGzip bool   `yaml:"enable_gzip"`
}

type DatabaseConfig struct {
	Type        string `yaml:"type"`
	Path        string `yaml:"path"`
	DSN         string `yaml:"dsn"`
	TablePrefix string `yaml:"table_prefix"`
	Hostname    string `yaml:"hostname"`
	Hostport    int64  `yaml:"hostport"`
	Name        string `yaml:"name"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	SSLMode     string `yaml:"ssl_mode"`
}

type SecurityConfig struct {
	InstallLock             bool   `yaml:"install_lock"`
	SecretKey               string `yaml:"secret_key"`
	LoginRememberDays       int    `yaml:"login_remember_days"`
	CookieRememberName      string `yaml:"cookie_remember_name"`
	CookieUsername          string `yaml:"cookie_username"`
	CookieSecure            bool   `yaml:"cookie_secure"`
	EnableLoginStatusCookie bool   `yaml:"enable_login_status_cookie"`
	LoginStatusCookieName   string `yaml:"login_status_cookie_name"`
}

func SaveConfData(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "marshal config")
	}
	return os.WriteFile(path, data, 0644)
}

func parseInt(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}
