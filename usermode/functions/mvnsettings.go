package functions

import (
	_ "embed"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

// MvnSettingsFunction structure
type MvnSettingsFunction struct {
}

// Infos function
func (f MvnSettingsFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "Maven ~/.m2/settings.xml file",
		ShortDescription: "Updates the ~/.m2/settings.xml file",
		LongDescription:  "Creates a default ~/.m2/settings.xml file IF the file does NOT already exist",
	}
}

// Run function
func (f MvnSettingsFunction) Run() result.Result {
	f1 := func() result.Result { return filesystem.WriteBinaryFile("~/.m2/settings.xml", assetMvnSettings, false) }

	return execute(f.Infos().Title, f1)
}

//go:embed assets/settings.xml
var assetMvnSettings []byte
