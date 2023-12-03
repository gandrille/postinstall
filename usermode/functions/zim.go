package functions

import (
	_ "embed"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

// Zim is probably the best desktop wiki ever created!
// http://zim-wiki.org/

// ZimFunction structure
type ZimFunction struct {
}

// Infos function
func (f ZimFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "Zim wiki",
		ShortDescription: "Better configuration for Zim",
		LongDescription:  "* ~/.config/zim/style.conf theme file is updated with a more modern theme",
	}
}

// Run function
func (f ZimFunction) Run() result.Result {

	f1 := func() result.Result {
		return filesystem.WriteBinaryFile("~/.config/zim/style.conf", assetZimStyle, true)
	}

	return execute(f.Infos().Title, f1)
}

//go:embed assets/zim-style.conf
var assetZimStyle []byte
