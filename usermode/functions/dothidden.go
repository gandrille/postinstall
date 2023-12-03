package functions

import (
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/postinstall/assets"
)

// DotHiddenFunction structure
type DotHiddenFunction struct {
}

// Infos function
func (f DotHiddenFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "~/.hidden file",
		ShortDescription: "Updates the ~/.hidden file",
		LongDescription:  "Creates a default ~/.hidden file IF the file does NOT already exist",
	}
}

// Run function
func (f DotHiddenFunction) Run() result.Result {
	f1 := func() result.Result { return assets.WriteAsset("resources/hidden", "~/.hidden", false) }

	return execute(f.Infos().Title, f1)
}
