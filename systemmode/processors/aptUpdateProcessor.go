package processors

import (
	"os/exec"

	"github.com/gandrille/go-commons/misc"
	"github.com/gandrille/go-commons/result"
)

type aptUpdate struct {
}

func (e aptUpdate) Key() string {
	return "apt-update"
}

func (e aptUpdate) Describe(args []string) string {
	return "apt update: updates packets database"
}

func (e aptUpdate) Run(args []string) result.Result {
	cmd := exec.Command("/usr/bin/apt", "update", "--yes")
	return misc.RunCmd(cmd, e.Describe(args))
}
