package functions

import (
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/postinstall/assets"
)

// BashrcFunction structure
type BashrcFunction struct {
}

// Infos function
func (f BashrcFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            ".bashrc ecosystem",
		ShortDescription: "Updates ~/.bashrc, ~/.bashrc-common, ~/.bashrc-perso",
		LongDescription: `* ~/.bashrc is updated in order to include the following files
* ~/.bashrc-common content is reverted to a default one
* ~/.bashrc-perso is initialized with a default content IF the file does NOT exist`,
	}
}

// Run function
func (f BashrcFunction) Run() result.Result {

	// .bashrc
	f1 := func() result.Result {
		return assets.CreateOrAppendAssetContentIfNotInFile("resources/bashrc", "~/.bashrc")
	}

	// .bashrc-common
	f2 := func() result.Result { return assets.WriteAsset("resources/bashrc-common", "~/.bashrc-common", true) }

	// .bashrc-perso
	f3 := func() result.Result { return assets.WriteAsset("resources/bashrc-perso", "~/.bashrc-perso", false) }

	return execute(f.Infos().Title, f1, f2, f3)
}
