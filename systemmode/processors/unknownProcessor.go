package processors

import (
	"strings"

	"github.com/gandrille/go-commons/result"
)

type unknown struct {
	str string
}

func (e unknown) Key() string {
	return e.str
}

func (e unknown) Describe(args []string) string {
	return "unknown command " + e.str + " " + strings.Join(args, " ")
}

func (e unknown) Run(args []string) result.Result {
	return result.NewError("NOT executed: " + e.Describe(args))
}
