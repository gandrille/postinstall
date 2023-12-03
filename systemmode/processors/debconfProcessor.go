package processors

import (
	"strings"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/result"
)

type debconf struct {
}

func (e debconf) Key() string {
	return "deb-conf"
}

func (e debconf) Describe(args []string) string {
	return "debconf: setting " + strings.Join(args, " ")
}

func (e debconf) Run(args []string) result.Result {
	packageName := args[0]
	keyName := args[1]
	typeName := args[2]
	value := args[3]
	return env.WriteDebconfKey(packageName, keyName, typeName, value)
}
