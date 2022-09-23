package functions

import (
	"os/exec"
	"strings"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

// DotHiddenFunction structure
type PythonFunction struct {
}

const pipPath = "/usr/bin/pip"

// Infos function
func (f PythonFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "Python packages installation",
		ShortDescription: "Install a set of python packages",
		LongDescription:  "Install the following python packages: slixmpp, pyperclip",
	}
}

// Run function
func (f PythonFunction) Run() result.Result {
	f1 := func() result.Result { return pipInstall("slixmpp") }
	f2 := func() result.Result { return pipInstall("pyperclip") }

	return execute(f.Infos().Title, f1, f2)
}

func pipInstall(name string) result.Result {
	// Check if pip executable exists
	exists, err := filesystem.RegularFileExists(pipPath)
	if err != nil {
		return result.NewError("Can't check if " + pipPath + " exists")
	}
	if !exists {
		return result.NewError(pipPath + " does NOT exist")
	}

	// check if package is already installed
	cmd1 := exec.Command("bash", "-lc", pipPath+" list | grep -E \"^"+name+" \" | wc -l")
	out1, err1 := cmd1.Output()
	if err1 != nil {
		return result.NewError("Can't check if the package " + name + " is already installed")
	}
	if strings.TrimSuffix(string(out1), "\n") == "1" {
		return result.NewUnchanged("Package " + name + " is already installed")
	}

	// install package
	cmd2 := exec.Command("pip", "install", name)
	_, err2 := cmd2.Output()
	if err2 != nil {
		return result.NewError("Error when installing python package " + name)
	}
	return result.NewUpdated("Package " + name + " successfully installed")
}
