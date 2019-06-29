package functions

import (
	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

// Blueman is a GTK+ Bluetooth Manager
// https://github.com/blueman-project/blueman

// BluemanFunction structure
type BluemanFunction struct {
}

// Infos function
func (f BluemanFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "Blueman download folder",
		ShortDescription: "Sets blueman shared folder",
		LongDescription:  "Sets blueman shared folder to the current value of $XDG_DOWNLOAD_DIR (to avoid a popup with an error message at startup).",
	}
}

// Run function
func (f BluemanFunction) Run() result.Result {

	dir, err := env.ReadXdgDir("DOWNLOAD")
	if err != nil {
		return result.NewError(err.Error())
	}

	key := "/org/blueman/transfer/shared-path"
	value := "'" + dir + "'"

	f1 := func() result.Result { return filesystem.CreateFolderIfNeeded(dir) }
	f2 := func() result.Result {
		return env.WriteDconfKey(key, value).StandardizeMessage("Blueman shared folder", value)
	}

	return execute(f.Infos().Title, f1, f2)
}
