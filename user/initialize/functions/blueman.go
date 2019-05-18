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
		ShortDescription: "Sets the blueman shared folder",
		LongDescription:  "Sets the blueman shared folder to an arbitrary existing folder to avoid a popup with an error message at startup.",
	}
}

// Run function
func (f BluemanFunction) Run() result.Result {
	dir := filesystem.HomeDir() + "/.config/blueman"
	key := "/org/blueman/transfer/shared-path"
	value := "'" + dir + "'"

	f1 := func() result.Result { return filesystem.CreateFolderIfNeeded(dir) }
	f2 := func() result.Result { return env.WriteDconfKey(key, value) }

	return execute(f.Infos().Title, f1, f2)
}
