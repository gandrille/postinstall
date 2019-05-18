package functions

import (
	"os/exec"
	"strings"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
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
* Updates the browser and email client desktop files for using a custom icon`,
	}
}

// Run function
func (f FreedesktopFunction) Run() result.Result {
	return execute(f.Infos().Title, xdgDirsConfig, desktopFilesConfig)
}

// ======================================
// ~/.config/user-dirs.dirs configuration
// ======================================

const exe = "/usr/bin/xdg-user-dirs-update"

func xdgDirsConfig() result.Result {
	hostName := env.Hostname()
	homeDir := filesystem.HomeDir()

	// check if executable exists
	if exists, err := filesystem.FileExists(exe); err != nil || exists == false {
		return result.New(false, "File "+exe+" does NOT exist")
	}

	// get previous content
	fic := homeDir + "/.config/user-dirs.dirs"
	oldContent, err1 := filesystem.ReadFileAsStringOrEmptyIfNotExists(fic)
	if err1 != nil {
		return result.New(false, err1.Error())
	}

	// path updates
	var inError []string

	inError = updateXdgDir(inError, homeDir, hostName, "TEMPLATES", "$HOME")
	inError = updateXdgDir(inError, homeDir, hostName, "PUBLICSHARE", "$HOME")
	inError = updateXdgDir(inError, homeDir, hostName, "DOCUMENTS", "$HOME")
	inError = updateXdgDir(inError, homeDir, hostName, "MUSIC", "$HOME")
	inError = updateXdgDir(inError, homeDir, hostName, "PICTURES", "$HOME")
	inError = updateXdgDir(inError, homeDir, hostName, "VIDEOS", "$HOME")
	inError = updateXdgDir(inError, homeDir, hostName, "DOWNLOAD", "$HOME")

	// if some path updates reported error, notify caller
	if len(inError) != 0 {
		return result.New(false, "At least one directory config failed: "+strings.Join(inError, ", "))
	}

	// get final content
	newContent, err2 := filesystem.ReadFileAsStringOrEmptyIfNotExists(fic)
	if err2 != nil {
		return result.New(false, err2.Error())
	}

	// return
	switch {
	case oldContent == "":
		return result.New(true, "~/.config/user-dirs.dirs created")
	case oldContent == newContent:
		return result.New(true, "~/.config/user-dirs.dirs no modification needed")
	default:
		return result.New(true, "~/.config/user-dirs.dirs updated")
	}
}

func updateXdgDir(inError []string, homeDir, hostName, key, value string) []string {
	cmd := exec.Command(exe, "--set", key, strings.Replace(strings.Replace(value, "$HOME", homeDir, -1), "$HOSTNAME", hostName, -1))
	if cmd.Run() == nil {
		return inError
	}
	return append(inError, key)
}

// ===========================
// desktop files configuration
// ===========================

func desktopFilesConfig() result.Result {

	srcDir := "/usr/share/applications/"
	dstDir := "~/.local/share/applications/"

	// Assert folder exists
	if res := filesystem.CreateFolderIfNeeded(dstDir); !res.IsSuccess() {
		return result.Failure(res.Message())
	}

	msg := ""
	// Mail
	if res := filesystem.CopyFileWithUpdate(srcDir+"exo-mail-reader.desktop", dstDir+"exo-mail-reader.desktop", "Icon=", "Icon=internet-mail", true); !res.IsSuccess() {
		return result.Failure(res.Message())
	} else {
		msg += res.Message()
	}

	// Browser
	if res := filesystem.CopyFileWithUpdate(srcDir+"exo-web-browser.desktop", dstDir+"exo-web-browser.desktop", "Icon=", "Icon=firefox", true); !res.IsSuccess() {
		return result.Failure(res.Message())
	} else {
		msg += "\n" + res.Message()
	}

	return result.Success(msg)
}
