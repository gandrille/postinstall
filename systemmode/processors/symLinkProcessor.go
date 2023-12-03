package processors

import (
	"strings"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

type symlink struct {
}

func (e symlink) Key() string {
	return "sym-link"
}

func (e symlink) Describe(args []string) string {
	return "Symbolic link: " + strings.Join(args, " --> ")
}

func (e symlink) Run(args []string) result.Result {
	return filesystem.UpdateOrCreateSymlink(args[0], args[1])
}
