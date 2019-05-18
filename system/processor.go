package system

import (
	"os/exec"
	"strings"

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
	list = append(list, fuse{})
	list = append(list, xfceplugins{})
	list = append(list, debconf{})
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
	return result.New(false, "NOT executed: "+e.describe(args))
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

// ==============
// fuse processor
// ==============

type fuse struct {
}

func (e fuse) key() string {
	return "fuse-conf"
}

func (e fuse) describe(args []string) string {
	return "fuse config: sets the user_allow_other option"
}

func (e fuse) run(args []string) result.Result {
	_, err := filesystem.CreateOrAppendIfNotInFile("/etc/fuse", "user_allow_other")
	if err != nil {
		return result.New(false, e.describe(args)+" "+err.Error())
	}
	return result.New(true, e.describe(args))
}

// ======================
// xfce-plugins processor
// ======================

type xfceplugins struct {
}

func (e xfceplugins) key() string {
	return "xfce-plugins"
}

func (e xfceplugins) describe(args []string) string {
	return "xfce-plugins: installs all the Xfce plugins (quite a lot)"
}

func (e xfceplugins) run(args []string) result.Result {
	out, err := exec.Command("/usr/bin/apt-cache", "search", "--names-only", "xfce4-").Output()
	if err != nil {
		return result.New(false, e.key()+" : Can NOT get list of Xfce plugins")
	}

	// lines have the following format :
	// xfce4-power-manager - Gestion de l'énergie...
	// 1. keeps only the package name
	// 2. asset the package name contains "plugin"
	var plugins []string
	for _, line := range strings.Split(string(out), "\n") {
		idx := strings.Index(line, " ")
		if idx != -1 {
			name := line[:idx-1]
			if strings.Index(name, "plugin") != -1 {
				plugins = append(plugins, name)
			}
		}
	}

	if len(plugins) == 0 {
		return result.New(true, e.key()+" NO Xfce plugin found")
	}

	return install(plugins, e.key()+" "+strings.Join(plugins, ","))
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
	str := strings.Join(args, " ")
	cmd := exec.Command("/usr/bin/debconf-set-selections", "-v")
	return misc.RunCmdStdIn("Debconf "+str, str+"\n", cmd)
}

// =========
// Internals
// =========

func install(args []string, describe string) result.Result {
	fullArgs := append([]string{"install", "--yes"}, args...)
	cmd := exec.Command("/usr/bin/apt", fullArgs...)
	return misc.RunCmd(cmd, describe)
}
