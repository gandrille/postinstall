package functions

import (
	"fmt"
	"os"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

// GnomeDesktopFunction structure
type GnomeDesktopFunction struct {
}

// Infos function
func (f GnomeDesktopFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "Gnome Desktop",
		ShortDescription: "Updates Gnome desktop configuration",
		LongDescription:  `Updates Gnome background image (used by Gnome screensaver/locker)`,
	}
}

// Run function
func (f GnomeDesktopFunction) Run() result.Result {

	// screen locker
	f1 := func() result.Result {
		path := "/usr/share/xfce4/backdrops/xubuntu-yakkety.png"
		return setGnomeBackgroundImage(path).StandardizeMessage("Gnome background image", path)
	}

	return execute(f.Infos().Title, f1)
}

func setGnomeBackgroundImage(path string) result.Result {

	// assert ~/.cache/wallpaper exists
	if res := filesystem.CreateFolderIfNeeded("~/.cache/wallpaper"); res.IsFailure() {
		fmt.Println("WARNING: " + res.Message())
	}

	// clear ~/.cache/wallpaper content
	if files, err := filesystem.FloderFiles("~/.cache/wallpaper"); err != nil {
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
