package functions

import (
	"strings"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/go-commons/strpair"
)

// FreedesktopFunction structure
type FreedesktopFunction struct {
}

// Infos function
func (f FreedesktopFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "Freedesktop",
		ShortDescription: "Updates configuration according to Freedesktop specification",
		LongDescription: `* Updates the ~/.config/user-dirs.dirs file with XDG directories
  paths (TEMPLATES, PUBLICSHARE,...) set with $HOME directory
* Updates the email client desktop files for using a custom icon
* Updates the browser client desktop files for using a custom icon`,
	}
}

// Run function
func (f FreedesktopFunction) Run() result.Result {
	return execute(f.Infos().Title, xdgDirsConfig, assertFolderExists, mailConfig, webConfig)
}

// ======================================
// ~/.config/user-dirs.dirs configuration
// ======================================

func xdgDirsConfig() result.Result {

	// path updates
	var inError []string
	allskipped := true

	dirs := []strpair.StrPair{
		strpair.New("TEMPLATES", "$HOME"),
		strpair.New("PUBLICSHARE", "$HOME"),
		strpair.New("DOCUMENTS", "$HOME"),
		strpair.New("MUSIC", "$HOME"),
		strpair.New("PICTURES", "$HOME"),
		strpair.New("VIDEOS", "$HOME")}

	for _, dir := range dirs {
		changed, err := env.UpdateXdgDir(dir.Str1(), dir.Str2())
		if err != nil {
			inError = append(inError, dir.Str1())
			allskipped = false
		} else {
			allskipped = allskipped && !changed
		}
	}

	if len(inError) != 0 {
		return result.NewError("At least one directory config failed: " + strings.Join(inError, ", "))
	}

	if allskipped {
		return result.NewUnchanged("~/.config/user-dirs.dirs no modification needed")
	}

	return result.NewUpdated("~/.config/user-dirs.dirs updated")
}

// ===========================
// desktop files configuration
// ===========================

const srcDir = "/usr/share/applications/"
const dstDir = "~/.local/share/applications/"

func assertFolderExists() result.Result {
	return filesystem.CreateFolderIfNeeded(dstDir)
}

func mailConfig() result.Result {
	name := "Mail"
	dst := "exo-mail-reader.desktop"
	icon := "internet-mail"
	return desktopFilesConfig(name, dst, icon)
}

func webConfig() result.Result {
	name := "Web"
	dst := "exo-web-browser.desktop"
	icon := "firefox"
	return desktopFilesConfig(name, dst, icon)
}

func desktopFilesConfig(name, ficName, icon string) result.Result {
	srcPath := srcDir + ficName
	dstPath := dstDir + ficName
	if res := filesystem.CopyFileWithUpdate(srcPath, dstPath, "Icon=", "Icon="+icon, true); !res.IsSuccess() {
		return result.NewError(res.Message())
	} else {
		msg := name + " icon: file " + dstPath
		return result.New(res.Status(), msg)
	}
}
