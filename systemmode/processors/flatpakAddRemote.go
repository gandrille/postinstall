package processors

import (
	"os/exec"
	"slices"
	"strings"

	"github.com/gandrille/go-commons/misc"
	"github.com/gandrille/go-commons/result"
)

type flatpakAddRemote struct {
}

func (e flatpakAddRemote) Key() string {
	return "flatpak-add-remote"
}

func (e flatpakAddRemote) Describe(args []string) string {
	return "Registers a flatpak remote: " + strings.Join(args, ",")
}

func (e flatpakAddRemote) Run(args []string) result.Result {
	cmd1 := exec.Command("/usr/bin/flatpak", "remotes")
	if remotes, err := runCmdAndGetFirstColumn(cmd1, true); err != nil {
		return result.NewError("Can't retrieve flatpak remotes list " + err.Error())
	} else if slices.Contains(remotes, args[0]) {
		return result.NewUnchanged("flatpak remote " + args[0] + " already installed")
	}

	// snap installation
	fullArgs := append([]string{"remote-add", "--if-not-exists"}, args...)
	cmd2 := exec.Command("/usr/bin/flatpak", fullArgs...)
	return misc.RunCmd(cmd2, e.Describe(args))
}
