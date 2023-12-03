package functions

import (
	_ "embed"
	"strings"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
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
		return updateBashrc()
	}

	// .bashrc-common
	f2 := func() result.Result { return filesystem.WriteBinaryFile("~/.bashrc-common", assetBashrcCommon, true) }

	// .bashrc-perso
	f3 := func() result.Result { return filesystem.WriteBinaryFile("~/.bashrc-perso", assetBashrcPerso, false) }

	return execute(f.Infos().Title, f1, f2, f3)
}

//go:embed assets/bashrc-common
var assetBashrcCommon []byte

//go:embed assets/bashrc-perso
var assetBashrcPerso []byte

//go:embed assets/bashrc-integration
var assetBashrcIntegration string

func updateBashrc() result.Result {
	sdkmanMarker := "#THIS MUST BE AT THE END OF THE FILE FOR SDKMAN TO WORK!!!"
	filePath := strings.Replace("~/.bashrc", "~", filesystem.HomeDir(), 1)

	// Checks if the file exists
	exists, err := filesystem.RegularFileExists(filePath)
	if err != nil {
		return result.NewError("Don't know if file " + filePath + " exists")
	}

	// The file exists
	if exists {

		// Checks if the content is already present
		contains1, err1 := filesystem.StringFileContains(filePath, assetBashrcIntegration)
		if err1 != nil {
			return result.NewError(err1.Error())
		}
		if contains1 {
			return result.NewUnchanged(filePath + " already updated with bash integration")
		}

		// Checks if SDKMAN! marker is present
		contains2, err2 := filesystem.StringFileContains(filePath, sdkmanMarker)
		if err2 != nil {
			return result.NewError(err2.Error())
		}
		if contains2 {
			res := filesystem.UpdateLineInFile(filePath, sdkmanMarker, assetBashrcIntegration+"\n\n"+sdkmanMarker, false)
			if res.IsSuccess() {
				return result.NewUpdated(filePath + " updated with bash integration")
			}
			return res
		}

		// Append content
		appended, err3 := filesystem.CreateOrAppendIfNotInFile(filePath, assetBashrcIntegration)
		if err3 != nil {
			return result.NewError(err3.Error())
		}

		if !appended {
			return result.NewError(filePath + "was not modified whereas it should have been")
		}

		return result.NewUpdated(filePath + " updated with bash integration")
	}

	// The file does NOT exist
	appended, err := filesystem.CreateOrAppendIfNotInFile(filePath, assetBashrcIntegration)
	if err != nil {
		return result.NewError(err.Error())
	}
	if appended {
		return result.NewCreated(filePath + " initialized with bash integration")
	}
	return result.NewError(filePath + "was not created whereas it should have been")
}
