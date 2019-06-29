package functions

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

// ScreensaverAndLockFunction structure
type ScreensaverAndLockFunction struct {
}

// Infos function
func (f ScreensaverAndLockFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "Screensaver/locker configuration",
		ShortDescription: "Updates the screensaver and the power-management configuration",
		LongDescription: `* Turn off screensaver
* Updates lock backgroung
* Updates lock command for using Gnome
* Updates lid switch management`,
	}
}

// Run function
func (f ScreensaverAndLockFunction) Run() result.Result {
	return execute(f.Infos().Title, screensaver, lockBackground, lockCommand, lid)
}

func screensaver() result.Result {
	return env.CreateOrSetXfconfProperty("xfce4-screensaver", "/saver/enabled", "bool", "false").StandardizeMessage("Screensaver", "<<disabled>>")
}

func lockBackground() result.Result {
	path := "/usr/share/xfce4/backdrops/xubuntu-yakkety.png"
	return setGnomeBackgroundImage(path).StandardizeMessage("Gnome Lock background image", path)
}

func lockCommand() result.Result {
	cmd := "gnome-screensaver-command --lock"
	return env.CreateOrSetXfconfProperty("xfce4-session", "/general/LockCommand", "string", cmd).StandardizeMessage("Lock command", cmd)
}

func lid() result.Result {
	logindFile := "/etc/systemd/logind.conf"
	contains, err := filesystem.StringFileContains(logindFile, "HandleLidSwitch=ignore")
	if err != nil {
		return result.NewError("Lid switch configuration: " + err.Error())
	}
	if !contains {
		return result.NewError("Lid switch configuration: please run sudo sed -i \"s/#HandleLidSwitch=suspend/HandleLidSwitch=ignore/\" " + logindFile)
	}

	res1 := env.SetXfconfProperty("xfce4-power-manager", "/xfce4-power-manager/logind-handle-lid-switch", "true")
	if res1.IsError() {
		result.NewError("Lid switch configuration (first check): " + res1.Message())
	}
	if res1.IsUnchanged() {
		result.NewUnchanged("Lid switch already configured")
	}

	// Power manager restart required to take into account new config
	if err := exec.Command("xfce4-power-manager", "--restart").Run(); err != nil {
		result.NewError("Lid switch configuration - power manager error: " + err.Error())
	}

	// Waiting that the service restarts...
	time.Sleep(3 * time.Second)

	// Checking again xfconf (yes, the value may be reverted... this is a known bug)
	res2 := env.SetXfconfProperty("xfce4-power-manager", "/xfce4-power-manager/logind-handle-lid-switch", "true")
	if res2.IsError() {
		result.NewError("Lid switch configuration (double check): " + res2.Message())
	}

	return result.NewUpdated("Lid switch configuration success. PLEASE REBOOT YOUR COMPUTER!")
}

func setGnomeBackgroundImage(path string) result.Result {

	// assert ~/.cache/wallpaper exists
	if res := filesystem.CreateFolderIfNeeded("~/.cache/wallpaper"); res.IsFailure() {
		fmt.Println("WARNING: " + res.Message())
	}

	// clear ~/.cache/wallpaper content
	if files, err := filesystem.FolderFiles("~/.cache/wallpaper"); err != nil {
		fmt.Println("WARNING: " + err.Error())
	} else {
		for _, file := range files {
			if err := os.Remove(file); err != nil {
				fmt.Println("WARNING: " + err.Error())
			}
		}
	}

	// write key
	return env.WriteGsettingsKey("org.gnome.desktop.background", "picture-uri", "'file://"+path+"'")
}
