package processors

import (
	"github.com/gandrille/go-commons/result"
)

// Function interface
type processor interface {
	Key() string
	Describe(args []string) string
	Run(args []string) result.Result
}

func GetProcessor(key string) processor {
	for _, p := range getProcessors() {
		if p.Key() == key {
			return p
		}
	}
	return unknown{key}
}

func getProcessors() []processor {
	var list []processor
	list = append(list, aptAutoremove{})
	list = append(list, aptInstall{})
	list = append(list, aptUpdate{})
	list = append(list, aptUpgrade{})
	list = append(list, debconf{})
	list = append(list, installDrawIO{})
	list = append(list, installVsCode{})
	list = append(list, fuse{})
	list = append(list, imagemagick{})
	list = append(list, inifile{})
	list = append(list, snapInstall{})
	list = append(list, flatpakAddRemote{})
	list = append(list, flatpakInstall{})
	list = append(list, flatpakLauncher{})
	list = append(list, symlink{})
	list = append(list, systemdtimesyncd{})
	list = append(list, unattendedUpgrade{})

	return list
}
