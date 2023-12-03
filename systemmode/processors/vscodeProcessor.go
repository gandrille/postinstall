package processors

import (
	"os/exec"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

type installVsCode struct {
}

func (e installVsCode) Key() string {
	return "vscode-install"
}

func (e installVsCode) Describe(args []string) string {
	return "Installs vscode from a deb file which is automatically downloaded"
}

func (e installVsCode) Run(args []string) result.Result {
	if exists, _ := filesystem.RegularFileExists("/usr/bin/code"); exists {
		return result.NewUnchanged("Vscode is already installed")
	}

	// Downloading deb file
	url := "https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-x64"
	if err := exec.Command("/usr/bin/wget", "-O", "vscode.deb", url).Run(); err != nil {
		return result.NewError("Error while downloading Vscode from " + url + " " + err.Error())
	}

	// Installing deb file
	if err := exec.Command("/usr/bin/apt", "install", "--yes", "./vscode.deb").Run(); err != nil {
		return result.NewError("Error while installing Vscode after dl from " + url + " " + err.Error())
	}

	// Removing deb file (nothing to do if it fails)
	exec.Command("/bin/rm", "./vscode.deb").Run()

	return result.NewUpdated("Vscode installed")
}
