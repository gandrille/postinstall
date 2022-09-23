package functions

import (
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/postinstall/assets"
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
		LongDescription: `* ~/.config/zim/style.conf theme file is updated with a more modern theme
* simple-web-template is deployed for HTTP rendering`,
	}
}

// Run function
func (f ZimFunction) Run() result.Result {

	// style.conf
	f1 := func() result.Result {
		return assets.WriteAsset("resources/zim-style.conf", "~/.config/zim/style.conf", true)
	}

	// SimpleWeb template
	f2 := func() result.Result {
		return assets.CopyAssetDirectory("resources/zim-simple-web-template", "~/.local/share/zim/templates/html", true)
	}

	return execute(f.Infos().Title, f1, f2)
}
