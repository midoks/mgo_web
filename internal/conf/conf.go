package conf

import (
	"fmt"
	"log"
	// "net/url"
	"os"
	"path/filepath"
	// "strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"

	"mgo/embed"
)

var File *ini.File

func ReadConf() (*ini.File, error) {
	cfg := ini.Empty()
	data, err := embed.Conf.ReadFile("conf/app.conf")
	if err != nil {
		return cfg, errors.Wrap(err, "read file 'conf/app.conf'")
	}

	File, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, data)

	File.NameMapper = ini.TitleUnderscore
	if err != nil {
		return cfg, errors.Wrap(err, "parse 'conf/app.conf'")
	}

	return cfg, nil
}

// creates a default configuration file if it doesn't exist
func autoMakeCustomConf(customConf string) error {
	if isExist(customConf) {
		return nil
	}

	// Create default configuration
	cfg := ini.Empty()
	if isFile(customConf) {
		if err := cfg.Append(customConf); err != nil {
			return errors.Wrap(err, "append existing config")
		}
	}

	// Set default values
	cfg.Section("").Key("app_name").SetValue("mgo")
	cfg.Section("").Key("run_mode").SetValue("prod")

	cfg.Section("web").Key("http_port").SetValue("9999")
	cfg.Section("web").Key("admin_path").SetValue("admin")

	cfg.Section("session").Key("provider").SetValue("memory")

	cfg.Section("database").Key("type").SetValue("sqlite3")
	cfg.Section("database").Key("path").SetValue("data/mgo.db")

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(customConf), os.ModePerm); err != nil {
		return errors.Wrap(err, "create config directory")
	}

	// Save configuration file
	if err := cfg.SaveTo(customConf); err != nil {
		return errors.Wrap(err, "save config file")
	}

	return nil
}

func InstallConf(data map[string]string) error {
	File, err := ReadConf()
	if err != nil {
		return err
	}

	err = renderSection(File)
	if err != nil {
		return err
	}

	customConf := filepath.Join(CustomDir(), "conf", "app.conf")

	if !isExist(filepath.Dir(customConf)) {
		os.MkdirAll(filepath.Dir(customConf), os.ModePerm)
	}

	File.Section("").Key("app_name").SetValue(App.Name)
	File.Section("").Key("brand_name").SetValue(App.BrandName)
	File.Section("").Key("run_user").SetValue(App.RunUser)
	File.Section("").Key("run_mode").SetValue("prod")

	// File.Section("log").Key("format").SetValue(Log.Format)
	File.Section("log").Key("root_path").SetValue(Log.RootPath)

	File.Section("web").Key("port").SetValue("9999")
	admin_path := fmt.Sprintf("/mgo_%s", randString(6))
	File.Section("web").Key("admin_path").SetValue(admin_path)

	if strings.EqualFold(data["type"], "mysql") {
		File.Section("database").Key("type").SetValue("mysql")
		File.Section("database").Key("hostname").SetValue(data["hostname"])
		File.Section("database").Key("hostport").SetValue(data["hostport"])
		File.Section("database").Key("name").SetValue(data["dbname"])
		File.Section("database").Key("user").SetValue(data["username"])
		File.Section("database").Key("password").SetValue(data["password"])
		File.Section("database").Key("table_prefix").SetValue(data["table_prefix"])
	} else if strings.EqualFold(data["type"], "sqlite3") {
		File.Section("database").Key("type").SetValue("sqlite3")
		File.Section("database").Key("path").SetValue(data["dbpath"])
	}

	File.Section("security").Key("install_lock").SetValue("true")
	File.Section("security").Key("secret_key").SetValue(randString(10))

	if err := File.SaveTo(customConf); err != nil {
		return err
	}

	// write custom configuration file, rewrite initialization read
	err = InitConf("")
	if err != nil {
		return err
	}
	return nil
}

// Init initializes the configuration system
func InitConf(customConf string) error {
	data, err := embed.Conf.ReadFile("conf/app.conf")
	if err != nil {
		return errors.Wrap(err, "read embedded config")
	}

	// Load embedded configuration
	File, err = ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, data)
	if err != nil {
		return errors.Wrap(err, "parse 'conf/app.conf'")
	}

	// Determine custom config path
	if customConf == "" {
		customConf = filepath.Join(CustomDir(), "conf", "app.conf")
		if err := autoMakeCustomConf(customConf); err != nil {
			return errors.Wrap(err, "create default config")
		}
	} else {
		customConf, err = filepath.Abs(customConf)
		if err != nil {
			return errors.Wrap(err, "get absolute path")
		}
	}
	CustomConf = customConf

	// Append custom configuration if exists
	if isFile(customConf) {
		if err = File.Append(customConf); err != nil {
			return errors.Wrapf(err, "append %q", customConf)
		}
	} else {
		log.Printf("Custom config %s not found. Ignore this warning if you're running for the first time", customConf)
	}

	File.NameMapper = ini.TitleUnderscore

	// Check run user when the install is locked.
	if Security.InstallLock {
		currentUser, match := CheckRunUser(App.RunUser)
		if !match {
			return fmt.Errorf("user configured to run imail is %q, but the current user is %q", App.RunUser, currentUser)
		}
	}

	err = renderSection(File)
	if err != nil {
		return err
	}

	return nil
}

func renderSection(File *ini.File) error {
	// Map default section to App struct
	if err := File.Section(ini.DefaultSection).MapTo(&App); err != nil {
		return errors.Wrap(err, "mapping default section")
	}

	// ****************************
	// ----- Web settings -----
	// ****************************

	if err := File.Section("web").MapTo(&Web); err != nil {
		return errors.Wrap(err, "mapping [web] section")
	}

	// ***************************
	// ----- Log settings -----
	// ***************************
	if err := File.Section("log").MapTo(&Log); err != nil {
		return errors.Wrap(err, "mapping [log] section")
	}

	// ***************************
	// ----- Security settings -----
	// ***************************
	if err := File.Section("database").MapTo(&Database); err != nil {
		return errors.Wrap(err, "mapping [database] section")
	}

	// ***************************
	// ----- Security settings -----
	// ***************************
	if err := File.Section("security").MapTo(&Security); err != nil {
		return errors.Wrap(err, "mapping [security] section")
	}
	return nil
}
