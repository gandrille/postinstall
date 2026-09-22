package functions

import (
	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/result"
)

// MimeFunction structure
type MimeFunction struct {
}

// Infos function
func (f MimeFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "Mimetype default app",
		ShortDescription: "Updates the application that is used by default to open some file",
		LongDescription:  "vlc will be used for mpeg files",
	}
}

// Run function
func (f MimeFunction) Run() result.Result {
	f1 := func() result.Result { return env.WriteMimeDefautApp("audio/mpeg", "vlc.desktop") }
	return execute(f.Infos().Title, f1)
}
