package functions

import (
	"os"
	"os/exec"
	"strconv"

	"strings"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/misc"
	"github.com/gandrille/go-commons/result"
)

// SdkManFunction structure
type SdkManFunction struct {
}

// Infos function
func (f SdkManFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "SdkMan! installation",
		ShortDescription: "Installs SdkMan!, Java11, and Java17",
		LongDescription: `* Installs SdkMan!
* Installs latest Java11 release
* Installs latest Java17 release`,
	}
}

// Run function
func (f SdkManFunction) Run() result.Result {
	return execute(f.Infos().Title, installSdkMan, setAutoAnswer, installMaven, installJDK11, installJDK17)
}

func installSdkMan() result.Result {
	if exists, _ := filesystem.RegularFileExists("~/.sdkman/bin/sdkman-init.sh"); exists {
		return result.NewUnchanged("SdkMan! is aready installed")
	}

	c1 := exec.Command("curl", "-s", "https://get.sdkman.io")
	c2 := exec.Command("bash")
	c2.Stdin, _ = c1.StdoutPipe()
	c2.Stdout = os.Stdout
	if err := c2.Start(); err != nil {
		return result.NewError("SdkMan! install: " + err.Error())
	}
	if err := c1.Run(); err != nil {
		return result.NewError("SdkMan! install: " + err.Error())
	}
	if err := c2.Wait(); err != nil {
		return result.NewError("SdkMan! install: " + err.Error())
	}
	return result.NewUpdated("SdkMan! installed")
}

func setAutoAnswer() result.Result {
	return filesystem.UpdateLineInFile("~/.sdkman/etc/config", "sdkman_auto_answer", "sdkman_auto_answer=true", true)
}

func installMaven() result.Result {
	return sdkManCommand("Maven", "~/.sdkman/candidates/maven/", "sdk install maven")
}

func installJDK11() result.Result {
	return installJDK(11, "open")
}

func installJDK17() result.Result {
	return installJDK(17, "oracle")
}

func installJDK(jdk int, flavor string) result.Result {

	// JDK version
	cmd1 := buildSdkCmd("sdk list java | cut -d '|' -f 6 | tr -s ' ' | grep ' " + strconv.Itoa(jdk) + ".' | grep -- '-" + flavor + "' | sed 's/^ //' | sed 's/ $//' | head -n 1")
	out, err := cmd1.Output()
	if err != nil {
		return result.NewError("Don't know JDK" + strconv.Itoa(jdk) + " version")
	}
	ver := string(out)
	ver = strings.TrimSuffix(ver, "\n")
	if ver == "" {
		return result.NewError("Empty JDK" + strconv.Itoa(jdk) + " version")
	}

	return sdkManCommand("Latest JDK"+strconv.Itoa(jdk)+" ("+ver+")", "~/.sdkman/candidates/java/"+ver, "sdk install java "+ver)
}

func sdkManCommand(name, checkPath, cmd string) result.Result {
	exists, err := filesystem.Exists(checkPath)
	if exists {
		return result.NewUnchanged(name + " is already installed")
	}
	if err != nil {
		return result.NewError("Can't know if " + name + " is installed")
	}

	return misc.RunCmd(buildSdkCmd(cmd), name+" installed")
}

func buildSdkCmd(sdkcmd string) *exec.Cmd {
	return exec.Command("bash", "-lc", "source ~/.sdkman/bin/sdkman-init.sh && "+sdkcmd)
}
