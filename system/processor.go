package system

import (
	"os/exec"
	"strings"
	"syscall"

	"github.com/gandrille/go-commons/ini"
	"github.com/gandrille/go-commons/systemctl"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/misc"
	"github.com/gandrille/go-commons/result"
)

// Function interface
type processor interface {
	key() string
	describe(args []string) string
	run(args []string) result.Result
}

func getProcessor(key string) processor {
	for _, p := range getProcessors() {
		if p.key() == key {
			return p
		}
	}
	return unknown{key}
}

func getProcessors() []processor {
	var list []processor
	list = append(list, aptUpdate{})
	list = append(list, aptUpgrade{})
	list = append(list, aptAutoremove{})
	list = append(list, aptInstall{})
	list = append(list, symlink{})
	list = append(list, inifile{})
	list = append(list, debconf{})
	list = append(list, fuse{})
	list = append(list, timezone{})
	list = append(list, systemdtimesyncd{})
	list = append(list, systemdlogind{})
	return list
}

// =================
// unknown processor
// =================

type unknown struct {
	str string
}

func (e unknown) key() string {
	return e.str
}

func (e unknown) describe(args []string) string {
	return "unknown command " + e.str + " " + strings.Join(args, " ")
}

func (e unknown) run(args []string) result.Result {
	return result.NewError("NOT executed: " + e.describe(args))
}

// ================
// update processor
// ================

type aptUpdate struct {
}

func (e aptUpdate) key() string {
	return "apt-update"
}

func (e aptUpdate) describe(args []string) string {
	return "apt update: updates packets database"
}

func (e aptUpdate) run(args []string) result.Result {
	cmd := exec.Command("/usr/bin/apt", "update", "--yes")
	return misc.RunCmd(cmd, e.describe(args))
}

// =================
// upgrade processor
// =================

type aptUpgrade struct {
}

func (e aptUpgrade) key() string {
	return "apt-upgrade"
}

func (e aptUpgrade) describe(args []string) string {
	return "apt upgrade: upgrades all the installed packets"
}

func (e aptUpgrade) run(args []string) result.Result {
	cmd := exec.Command("/usr/bin/apt", "upgrade", "--yes")
	return misc.RunCmd(cmd, e.describe(args))
}

// ====================
// autoremove processor
// ====================

type aptAutoremove struct {
}

func (e aptAutoremove) key() string {
	return "apt-autoremove"
}

func (e aptAutoremove) describe(args []string) string {
	return "apt autoremove: removes packages not necessary anymore"
}

func (e aptAutoremove) run(args []string) result.Result {
	cmd := exec.Command("/usr/bin/apt", "autoremove", "--yes")
	return misc.RunCmd(cmd, e.describe(args))
}

// =================
// apt get processor
// =================

type aptInstall struct {
}

func (e aptInstall) key() string {
	return "apt-install"
}

func (e aptInstall) describe(args []string) string {
	return "package installation: " + strings.Join(args, ",")
}

func (e aptInstall) run(args []string) result.Result {
	return install(args, e.describe(args))
}

// ==================
// sym-link processor
// ==================

type symlink struct {
}

func (e symlink) key() string {
	return "sym-link"
}

func (e symlink) describe(args []string) string {
	return "Symbolic link: " + strings.Join(args, " --> ")
}

func (e symlink) run(args []string) result.Result {
	return filesystem.UpdateOrCreateSymlink(args[0], args[1])
}

// ==================
// ini-file processor
// ==================

type inifile struct {
}

func (e inifile) key() string {
	return "ini-file"
}

func (e inifile) describe(args []string) string {
	return "INI file: " + args[0] + " [" + args[1] + "] " + args[2] + " = " + args[3]
}

func (e inifile) run(args []string) result.Result {
	updated, err := ini.SetValue(args[0], args[1], args[2], args[3], false, false)
	if err != nil {
		return result.NewError(e.describe(args) + " " + err.Error())
	}
	if updated {
		return result.NewUpdated(e.describe(args))
	} else {
		return result.NewUnchanged(e.describe(args))
	}
}

// =================
// debconf processor
// =================

type debconf struct {
}

func (e debconf) key() string {
	return "deb-conf"
}

func (e debconf) describe(args []string) string {
	return "debconf: setting " + strings.Join(args, " ")
}

func (e debconf) run(args []string) result.Result {
	packageName := args[0]
	keyName := args[1]
	typeName := args[2]
	value := args[3]
	return env.WriteDebconfKey(packageName, keyName, typeName, value)
}

// ==============
// fuse processor
// ==============

type fuse struct {
}

func (e fuse) key() string {
	return "configure-fuse"
}

func (e fuse) describe(args []string) string {
	return "fuse config: sets the user_allow_other option"
}

func (e fuse) run(args []string) result.Result {
	apended, err := filesystem.CreateOrAppendIfNotInFile("/etc/fuse", "user_allow_other")
	if err != nil {
		return result.NewError(e.describe(args) + " " + err.Error())
	}
	if apended {
		return result.NewUpdated("File /etc/fuse updated with user_allow_other option")
	} else {
		return result.NewUnchanged("File /etc/fuse already contains user_allow_other option")
	}
}

// ==================
// timezone processor
// ==================

type timezone struct {
}

func (e timezone) key() string {
	return "configure-timezone"
}

func (e timezone) describe(args []string) string {
	return "timezone config: sets the timezone with " + args[0]
}

func (e timezone) run(args []string) result.Result {
	res := filesystem.WriteStringFile("/etc/timezone", args[0]+"\n", true)

	if res.IsUnchanged() {
		return result.NewUnchanged("Timezone is already " + args[0])
	}
	if res.IsUpdated() || res.IsCreated() {
		return result.NewUpdated("Timezone is now " + args[0])
	}

	return result.NewError("Timezone update failed : " + res.Message())
}

// ==========================
// systemdtimesyncd processor
// ==========================

type systemdtimesyncd struct {
}

func (e systemdtimesyncd) key() string {
	return "configure-systemd-timesyncd"
}

func (e systemdtimesyncd) describe(args []string) string {
	return "NTP configuration using server " + args[0]
}

func (e systemdtimesyncd) run(args []string) result.Result {
	server := args[0]
	service := "systemd-timesyncd.service"
	updated := false

	// Update config file
	fileupdated, err1 := ini.SetValue("/etc/systemd/timesyncd.conf", "Time", "NTP", server, false, true)
	if err1 != nil {
		return result.NewError(e.describe(args) + " " + err1.Error())
	}
	msg := ""
	if fileupdated {
		msg = "configuration updated, "
		updated = true
	} else {
		msg = "configuration already updated, "
	}

	// (re)start service
	if fileupdated {
		restarted, err2 := systemctl.Restart(service)
		if err2 != nil {
			return result.NewError(e.describe(args) + " " + err2.Error())
		}
		if restarted {
			msg += "service restarted, "
		} else {
			msg += "service started, "
		}
		updated = true
	} else {
		activated, err3 := systemctl.Activate(service)
		if err3 != nil {
			return result.NewError(e.describe(args) + " " + err3.Error())
		}
		if activated {
			msg += "service activated, "
			updated = true
		} else {
			msg += "service already started, "
		}
	}

	// Enable service
	enabled, err4 := systemctl.Enable(service)
	if err4 != nil {
		return result.NewError(e.describe(args) + " " + err4.Error())
	}
	if enabled {
		msg += "service enabled"
		updated = true
	} else {
		msg += "service already enabled"
	}

	if updated {
		return result.NewUpdated(msg)
	}
	return result.NewUnchanged(msg)
}

// ==================================
// configure-systemd-logind processor
// ==================================

type systemdlogind struct {
}

func (e systemdlogind) key() string {
	return "configure-systemd-logind"
}

func (e systemdlogind) describe(args []string) string {
	return "Disables lid switch management with systemd logind"
}

func (e systemdlogind) run(args []string) result.Result {
	path := "/etc/systemd/logind.conf"
	old := "#HandleLidSwitch=suspend"
	new := "HandleLidSwitch=ignore"

	// Note to myself
	// $old may or may not be present before calling this method
	// Thats why this method looks a bit over complicated

	// is update needed
	contains, err := filesystem.StringFileContains(path, new)
	if err != nil {
		return result.NewError(e.describe(args) + " " + err.Error())
	}
	if contains {
		return result.NewUnchanged(e.describe(args))
	}

	// update needed
	res := filesystem.UpdateLineInFile(path, old, new, true)
	if res.IsError() {
		return result.NewError(e.describe(args) + " " + res.Message())
	}

	return result.NewUpdated(e.describe(args))
}

// =========
// Internals
// =========

func install(args []string, describe string) result.Result {

	// No need to install
	if isAllInstalled(args) {
		switch len(args) {
		case 0:
			return result.NewError("NO Packages provided on the configuration file. Please double check it.")
		case 1:
			return result.NewUnchanged("Package " + args[0] + " already installed")
		default:
			return result.NewUnchanged("Packages " + strings.Join(args, ", ") + " already installed")
		}
	}

	// installation needed
	fullArgs := append([]string{"install", "--yes"}, args...)
	cmd := exec.Command("/usr/bin/apt", fullArgs...)
	return misc.RunCmd(cmd, describe)
}

func isAllInstalled(packageNames []string) bool {
	for _, name := range packageNames {
		if ok, _ := isInstalled(name); !ok {
			return false
		}
	}
	return true
}

func isInstalled(packageName string) (bool, error) {
	out, err := exec.Command("/usr/bin/dpkg-query", "-W", "-f=${db:Status-Abbrev}", packageName).Output()
	if err != nil {
		status := getStatusCode(err)
		switch status {
		case 0:
			return true, nil
		case 1:
			return false, nil
		default:
			return false, err
		}
	} else {
		output := string(out)
		if output == "ii" || output == "ii " {
			return true, nil
		} else {
			return false, nil
		}
	}
}

// TODO when upgrading to golang 1.12, use exitError.ExitCode() instead
// if exitError, ok := err.(*exec.ExitError); ok {
//   code := exitError.ExitCode()
// }
// return -1 if unknown
func getStatusCode(err error) int {
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	} else {
		return -1
	}
	return -1
}
