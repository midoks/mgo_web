package conf


import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"net/url"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"

	"mgo/embed"

)

var File *ini.File


// creates a default configuration file if it doesn't exist
func autoMakeCustomConf(customConf string) error {
	if IsExist(customConf) {
		return nil
	}

	// Create default configuration
	cfg := ini.Empty()
	if IsFile(customConf) {
		if err := cfg.Append(customConf); err != nil {
			return errors.Wrap(err, "append existing config")
		}
	}

	// Set default values
	cfg.Section("").Key("app_name").SetValue("dztasks")
	cfg.Section("").Key("run_mode").SetValue("prod")

	cfg.Section("web").Key("http_port").SetValue("9921")
	cfg.Section("session").Key("provider").SetValue("memory")

	cfg.Section("database").Key("type").SetValue("sqlite3")
	cfg.Section("database").Key("path").SetValue("data/mgo.db")

	cfg.Section("plugins").Key("path").SetValue("plugins")
	cfg.Section("plugins").Key("show_error").SetValue("true")
	cfg.Section("plugins").Key("show_cmd").SetValue("true")

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

// Init initializes the configuration system
func Init(customConf string) error {
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
	if IsFile(customConf) {
		if err = File.Append(customConf); err != nil {
			return errors.Wrapf(err, "append %q", customConf)
		}
	} else {
		log.Printf("Custom config %s not found. Ignore this warning if you're running for the first time", customConf)
	}

	File.NameMapper = ini.TitleUnderscore

	// Map default section to App struct
	if err = File.Section(ini.DefaultSection).MapTo(&App); err != nil {
		return errors.Wrap(err, "mapping default section")
	}

	// ***************************
	// ----- Log settings -----
	// ***************************
	if err = File.Section("log").MapTo(&Log); err != nil {
		return errors.Wrap(err, "mapping [log] section")
	}

	// ****************************
	// ----- Web settings -----
	// ****************************

	if err = File.Section("web").MapTo(&Web); err != nil {
		return errors.Wrap(err, "mapping [web] section")
	}

	Web.AppDataPath = ensureAbs(Web.AppDataPath)

	if !strings.HasSuffix(Web.ExternalURL, "/") {
		Web.ExternalURL += "/"
	}
	Web.URL, err = url.Parse(Web.ExternalURL)
	if err != nil {
		return errors.Wrapf(err, "parse '[server] EXTERNAL_URL' %q", err)
	}

	// Subpath should start with '/' and end without '/', i.e. '/{subpath}'.
	Web.Subpath = strings.TrimRight(Web.URL.Path, "/")
	Web.SubpathDepth = strings.Count(Web.Subpath, "/")

	unixSocketMode, err := strconv.ParseUint(Web.UnixSocketPermission, 8, 32)
	if err != nil {
		return errors.Wrapf(err, "parse '[server] unix_socket_permission' %q", Web.UnixSocketPermission)
	}
	if unixSocketMode > 0777 {
		unixSocketMode = 0666
	}
	Web.UnixSocketMode = os.FileMode(unixSocketMode)

	// ****************************
	// ----- Session settings -----
	// ****************************
	if err = File.Section("session").MapTo(&Session); err != nil {
		return errors.Wrap(err, "mapping [session] section")
	}

	// ***************************
	// ----- Security settings -----
	// ***************************
	if err := File.Section("database").MapTo(&Database); err != nil {
		return errors.Wrap(err, "mapping [database] section")
	}

	// *****************************
	// ----- Security settings -----
	// *****************************
	if err = File.Section("security").MapTo(&Security); err != nil {
		return errors.Wrap(err, "mapping [security] section")
	}

	// Check run user when the install is locked.
	if Security.InstallLock {
		currentUser, match := CheckRunUser(App.RunUser)
		if !match {
			return fmt.Errorf("user configured to run imail is %q, but the current user is %q", App.RunUser, currentUser)
		}
	}

	return nil
}
