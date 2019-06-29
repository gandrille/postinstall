package functions

import (
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/postinstall/assets"
	"github.com/gandrille/postinstall/initialize/firefox"
)

// FirefoxFunction structure
type FirefoxFunction struct {
}

const engineName = "Qwant"

// Infos function
func (f FirefoxFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "Firefox configuration",
		ShortDescription: "Sets nice defaults for Firefox",
		LongDescription: `* Creates a default profile if needed
* Creates a user.js file for pushing nice default settings
* Changes default search engine for using ` + engineName + " (requires that Firefox has already been started once)",
	}
}

// Run function
func (f FirefoxFunction) Run() result.Result {
	return execute(f.Infos().Title, firefox.CreateProfile, setUserConfig, setSearchEngine)
}

func setUserConfig() result.Result {
	path, err := firefox.GetProfileFolder()
	if err != nil {
		return result.NewError(err.Error())
	}
	return assets.WriteAsset("resources/firefox-user.js", path+"/user.js", true)
}

func setSearchEngine() result.Result {
	return firefox.SetSearchEngine(engineName)
}
