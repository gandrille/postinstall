package processors

import (
	"os/exec"
	"strings"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

type installDrawIO struct {
}

func (e installDrawIO) Key() string {
	return "drawio-install"
}

func (e installDrawIO) Describe(args []string) string {
	return "Installs draw.io from a deb file which is automatically downloaded"
}

func (e installDrawIO) Run(args []string) result.Result {
	if exists, _ := filesystem.RegularFileExists("/usr/bin/drawio"); exists {
		return result.NewUnchanged("Draw.io is already installed")
	}

	// Finding latest release
	out1, err1 := exec.Command("bash", "-lc", "curl -si https://github.com/jgraph/drawio-desktop/releases/latest | grep -E '^location' | sed 's|^.*tag/v||'").Output()
	if err1 != nil {
		return result.NewError("Error while finding latest Draw.io release " + err1.Error())
	}
	version := strings.Trim(string(out1), " \n\r")

	// Downloading deb file
	url := "https://github.com/jgraph/drawio-desktop/releases/download/v" + version + "/drawio-amd64-" + version + ".deb"
	if err := exec.Command("/usr/bin/wget", "-O", "drawio.deb", url).Run(); err != nil {
		return result.NewError("Error while downloading Draw.io from " + url + " " + err.Error())
	}

	// Installing deb file
	if err := exec.Command("/usr/bin/apt", "install", "--yes", "./drawio.deb").Run(); err != nil {
		return result.NewError("Error while installing Draw.io after dl from " + url + " " + err.Error())
	}

	// Removing deb file (nothing to do if it fails)
	exec.Command("/bin/rm", "./drawio.deb").Run()

	return result.NewUpdated("Draw.io v" + version + " installed")
}
