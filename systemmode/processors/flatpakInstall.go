package processors

import (
	"os/exec"
	"slices"
	"strings"

	"github.com/gandrille/go-commons/misc"
	"github.com/gandrille/go-commons/result"
)

type flatpakInstall struct {
}

func (e flatpakInstall) Key() string {
	return "flatpak-install"
}

func (e flatpakInstall) Describe(args []string) string {
	return "Installs a flatpak package: " + strings.Join(args, ",")
}

func (e flatpakInstall) Run(args []string) result.Result {
	cmd1 := exec.Command("/usr/bin/flatpak", "--columns=app", "list")
	if remotes, err := runCmdAndGetFirstColumn(cmd1, true); err != nil {
		return result.NewError("Can't retrieve flatpak applications list " + err.Error())
	} else if slices.Contains(remotes, args[1]) {
		return result.NewUnchanged("flatpak application " + args[1] + " is already installed")
	}

	// snap installation
	fullArgs := append([]string{"install", "--noninteractive", "--assumeyes"}, args...)
	cmd2 := exec.Command("/usr/bin/flatpak", fullArgs...)
	return misc.RunCmd(cmd2, e.Describe(args))
}
