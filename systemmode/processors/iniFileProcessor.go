package processors

import (
	"github.com/gandrille/go-commons/ini"
	"github.com/gandrille/go-commons/result"
)

type inifile struct {
}

func (e inifile) Key() string {
	return "ini-file"
}

func (e inifile) Describe(args []string) string {
	return "INI file: " + args[0] + " [" + args[1] + "] " + args[2] + " = " + args[3]
}

func (e inifile) Run(args []string) result.Result {
	updated, err := ini.SetValue(args[0], args[1], args[2], args[3], false, false)
	if err != nil {
		return result.NewError(e.Describe(args) + " " + err.Error())
	}
	if updated {
		return result.NewUpdated(e.Describe(args))
	} else {
		return result.NewUnchanged(e.Describe(args))
	}
}
