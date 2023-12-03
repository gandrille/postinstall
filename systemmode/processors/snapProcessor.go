package processors

import (
	"os/exec"
	"slices"
	"strings"

	"github.com/gandrille/go-commons/misc"
	"github.com/gandrille/go-commons/result"
)

type snapInstall struct {
}

func (e snapInstall) Key() string {
	return "snap-install"
}

func (e snapInstall) Describe(args []string) string {
	return "snap installation: " + strings.Join(args, ",")
}

func (e snapInstall) Run(args []string) result.Result {
	cmd1 := exec.Command("/usr/bin/snap", "list")
	if snaps, err := runCmdAndGetFirstColumn(cmd1, false); err != nil {
		return result.NewError("Can't retrieve snap list " + err.Error())
	} else if slices.Contains(snaps, args[0]) {
		return result.NewUnchanged("snap " + args[0] + " already installed")
	}

	// snap installation
	fullArgs := append([]string{"install"}, args...)
	cmd2 := exec.Command("/usr/bin/snap", fullArgs...)
	return misc.RunCmd(cmd2, e.Describe(args))
}
