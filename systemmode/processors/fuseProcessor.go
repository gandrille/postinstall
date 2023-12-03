package processors

import (
	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

type fuse struct {
}

func (e fuse) Key() string {
	return "configure-fuse"
}

func (e fuse) Describe(args []string) string {
	return "fuse config: sets the user_allow_other option"
}

func (e fuse) Run(args []string) result.Result {
	apended, err := filesystem.CreateOrAppendIfNotInFile("/etc/fuse", "user_allow_other")
	if err != nil {
		return result.NewError(e.Describe(args) + " " + err.Error())
	}
	if apended {
		return result.NewUpdated("File /etc/fuse updated with user_allow_other option")
	} else {
		return result.NewUnchanged("File /etc/fuse already contains user_allow_other option")
	}
}
