package processors

import (
	"os/exec"

	"github.com/gandrille/go-commons/misc"
	"github.com/gandrille/go-commons/result"
)

type aptUpgrade struct {
}

func (e aptUpgrade) Key() string {
	return "apt-upgrade"
}

func (e aptUpgrade) Describe(args []string) string {
	return "apt upgrade: upgrades all the installed packets"
}

func (e aptUpgrade) Run(args []string) result.Result {
	cmd := exec.Command("/usr/bin/apt", "upgrade", "--yes")
	return misc.RunCmd(cmd, e.Describe(args))
}
