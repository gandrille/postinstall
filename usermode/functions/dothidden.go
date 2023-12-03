package functions

import (
	_ "embed"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
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
	f1 := func() result.Result { return filesystem.WriteBinaryFile("~/.hidden", assetDotHidden, false) }
	return execute(f.Infos().Title, f1)
}

//go:embed assets/hidden
var assetDotHidden []byte
