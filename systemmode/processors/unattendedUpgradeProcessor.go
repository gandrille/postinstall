package processors

import (
	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

type unattendedUpgrade struct {
}

func (e unattendedUpgrade) Key() string {
	return "configure-unattendedUpgrade"
}

func (e unattendedUpgrade) Describe(args []string) string {
	return "Enables unattended-upgrade for standards updates"
}

func (e unattendedUpgrade) Run(args []string) result.Result {
	configFile := "/etc/apt/apt.conf.d/50unattended-upgrades"
	return filesystem.UpdateLineInFile(configFile, "//	\"${distro_id}:${distro_codename}-updates\";", "	\"${distro_id}:${distro_codename}-updates\";", false)
}
