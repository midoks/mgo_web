package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

var Install = cli.Command{
	Name:   "install",
	Usage:  "install mgo as a systemd service (Linux only)",
	Action: runInstall,
	Flags: []cli.Flag{
		stringFlag("name, n", "mgo", "service name"),
		stringFlag("user, u", "root", "system user to run the service"),
		stringFlag("binary, b", "", "path to mgo binary"),
	},
}

var Uninstall = cli.Command{
	Name:   "uninstall",
	Usage:  "uninstall mgo from systemd service (Linux only)",
	Action: runUninstall,
	Flags: []cli.Flag{
		stringFlag("name, n", "mgo", "service name"),
	},
}

const serviceTemplate = `[Unit]
Description=MGO Service
After=network.target

[Service]
Type=simple
User={{.User}}
ExecStart={{.Binary}} web
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
`

type serviceConfig struct {
	Name   string
	User   string
	Binary string
}

func runInstall(c *cli.Context) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("install command is only supported on Linux")
	}

	name := c.String("name")
	user := c.String("user")
	binary := c.String("binary")

	if binary == "" {
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("cannot get executable path")
		}
		binary = exePath
	}

	absPath, err := filepath.Abs(binary)
	if err != nil {
		return fmt.Errorf("cannot resolve binary path")
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("binary not found at: %s", absPath)
	}

	cfg := serviceConfig{
		Name:   name,
		User:   user,
		Binary: absPath,
	}

	t, _ := template.New("service").Parse(serviceTemplate)
	var sb strings.Builder
	t.Execute(&sb, cfg)

	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", name)
	if err := os.WriteFile(servicePath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("write service file failed")
	}

	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("systemctl daemon-reload failed")
	}

	if err := exec.Command("systemctl", "enable", name).Run(); err != nil {
		return fmt.Errorf("systemctl enable failed")
	}

	fmt.Printf("Installed: %s\n", name)
	fmt.Printf("Binary: %s\n", absPath)
	fmt.Printf("Start: systemctl start %s\n", name)
	return nil
}

func runUninstall(c *cli.Context) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("uninstall command is only supported on Linux")
	}

	name := c.String("name")
	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", name)

	if _, err := os.Stat(servicePath); os.IsNotExist(err) {
		return fmt.Errorf("service '%s' not installed", name)
	}

	exec.Command("systemctl", "stop", name).Run()
	exec.Command("systemctl", "disable", name).Run()
	os.Remove(servicePath)
	exec.Command("systemctl", "daemon-reload").Run()

	fmt.Printf("Uninstalled: %s\n", name)
	return nil
}
