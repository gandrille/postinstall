package processors

import (
	"os/exec"

	"github.com/gandrille/go-commons/misc"
	"github.com/gandrille/go-commons/result"
)

type aptAutoremove struct {
}

func (e aptAutoremove) Key() string {
	return "apt-autoremove"
}

func (e aptAutoremove) Describe(args []string) string {
	return "apt autoremove: removes packages not necessary anymore"
}

func (e aptAutoremove) Run(args []string) result.Result {
	cmd := exec.Command("/usr/bin/apt", "autoremove", "--yes")
	return misc.RunCmd(cmd, e.Describe(args))
}
