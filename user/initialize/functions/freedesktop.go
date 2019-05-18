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
	return execute(f.Infos().Title, xdgDirsConfig, desktopFilesConfig)
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

func desktopFilesConfig() result.Result {

	srcDir := "/usr/share/applications/"
	dstDir := "~/.local/share/applications/"

	allskipped := true

	// Assert folder exists
	if res := filesystem.CreateFolderIfNeeded(dstDir); !res.IsSuccess() {
		return result.NewError(res.Message())
	} else {
		allskipped = allskipped && res.IsUnchanged()
	}

	msg := ""
	// Mail
	dstMail := dstDir + "exo-mail-reader.desktop"
	if res := filesystem.CopyFileWithUpdate(srcDir+"exo-mail-reader.desktop", dstMail, "Icon=", "Icon=internet-mail", true); !res.IsSuccess() {
		return result.NewError(res.Message())
	} else {
		newMsg := "Mail icon: file " + dstMail + " " + strings.ToLower(res.Status().String())
		msg += newMsg
		allskipped = allskipped && res.IsUnchanged()
	}

	// Browser
	if res := filesystem.CopyFileWithUpdate(srcDir+"exo-web-browser.desktop", dstDir+"exo-web-browser.desktop", "Icon=", "Icon=firefox", true); !res.IsSuccess() {
		return result.NewError(res.Message())
	} else {
		msg += "\n"
		newMsg := "Browser icon: file " + dstMail + " " + strings.ToLower(res.Status().String())
		msg += newMsg
		allskipped = allskipped && res.IsUnchanged()
	}

	if allskipped {
		return result.NewUnchanged(msg)
	} else {
		return result.NewUpdated(msg)
	}
}
