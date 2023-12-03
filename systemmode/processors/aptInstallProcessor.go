package processors

import (
	"os/exec"
	"strings"
	"syscall"

	"github.com/gandrille/go-commons/misc"
	"github.com/gandrille/go-commons/result"
)

type aptInstall struct {
}

func (e aptInstall) Key() string {
	return "apt-install"
}

func (e aptInstall) Describe(args []string) string {
	return "package installation: " + strings.Join(args, ",")
}

func (e aptInstall) Run(args []string) result.Result {
	return install(args, e.Describe(args))
}

// =========
// Internals
// =========

func install(args []string, describe string) result.Result {

	// No need to install
	if isAllInstalled(args) {
		switch len(args) {
		case 0:
			return result.NewError("NO Packages provided on the configuration file. Please double check it.")
		case 1:
			return result.NewUnchanged("Package " + args[0] + " already installed")
		default:
			return result.NewUnchanged("Packages " + strings.Join(args, ", ") + " already installed")
		}
	}

	// installation needed
	fullArgs := append([]string{"install", "--yes"}, args...)
	cmd := exec.Command("/usr/bin/apt", fullArgs...)
	return misc.RunCmd(cmd, describe)
}

func isAllInstalled(packageNames []string) bool {
	for _, name := range packageNames {
		if ok, _ := isInstalled(name); !ok {
			return false
		}
	}
	return true
}

func isInstalled(packageName string) (bool, error) {
	out, err := exec.Command("/usr/bin/dpkg-query", "-W", "-f=${db:Status-Abbrev}", packageName).Output()
	if err != nil {
		status := getStatusCode(err)
		switch status {
		case 0:
			return true, nil
		case 1:
			return false, nil
		default:
			return false, err
		}
	} else {
		output := string(out)
		// TODO next line is probably a bug...
		if output == "ii" || output == "ii " {
			return true, nil
		} else {
			return false, nil
		}
	}
}

// TODO when upgrading to golang 1.12, use exitError.ExitCode() instead
//
//	if exitError, ok := err.(*exec.ExitError); ok {
//	  code := exitError.ExitCode()
//	}
//
// return -1 if unknown
func getStatusCode(err error) int {
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	} else {
		return -1
	}
	return -1
}
