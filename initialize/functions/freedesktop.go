package functions

import (
	"strings"

	"github.com/gandrille/go-commons/env"
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
		LongDescription: `Updates the ~/.config/user-dirs.dirs file with XDG directories
  paths (TEMPLATES, PUBLICSHARE,...) set with $HOME directory`,
	}
}

// Run function
func (f FreedesktopFunction) Run() result.Result {
	return execute(f.Infos().Title, xdgDirsConfig)
}

// ======================================
// ~/.config/user-dirs.dirs configuration
// ======================================

func xdgDirsConfig() result.Result {

	// path updates
	var inError []string
	allskipped := true

	dirs := []strpair.StrPair{
		strpair.New("DOWNLOAD", "$HOME"),
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
