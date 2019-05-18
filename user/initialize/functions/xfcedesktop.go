package functions

import (
	"strings"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/result"
)

// XfceDesktopFunction structure
type XfceDesktopFunction struct {
}

// Infos function
func (f XfceDesktopFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "Xfce Desktop",
		ShortDescription: "Setup the desktop configuration",
		LongDescription: `Updates:
* Workspace image background
* Workspace image background style
* Workspace count
* Window manager theme
* Apparence / Style (colors)
* Apparence / Icon set
* <print screen> keyboard shortcut`,
	}
}

// Run function
func (f XfceDesktopFunction) Run() result.Result {

	// Workspace image background
	f1 := func() result.Result { return setBackgroundImage() }

	// Workspace count
	f2 := func() result.Result { return env.SetXfconfProperty("xfwm4", "/general/workspace_count", "3") }

	// Window manager theme
	f3 := func() result.Result { return env.SetXfconfProperty("xfwm4", "/general/theme", "Arc") }

	// Apparence / Style (colors)
	f4 := func() result.Result { return env.SetXfconfProperty("xsettings", "/Net/ThemeName", "Arc") }

	// Apparence / Icon set
	f5 := func() result.Result { return env.SetXfconfProperty("xsettings", "/Net/IconThemeName", "deepin") }

	// Apparence / Icon set
	f6 := func() result.Result {
		return env.SetXfconfProperty("xfce4-keyboard-shortcuts", "/commands/custom/Print", "xfce4-screenshooter")
	}

	return execute(f.Infos().Title, f1, f2, f3, f4, f5, f6)
}

func setBackgroundImage() result.Result {
	channel := "xfce4-desktop"
	image := "/usr/share/xfce4/backdrops/xubuntu-yakkety.png"

	properties, err := env.ListXfconfProperties(channel)
	if err != nil {
		return result.Failure("Can't list XfconfProperties from " + channel)
	}

	msg := ""
	for _, prop := range properties {
		if strings.HasSuffix(prop, "/last-image") {
			if res := env.SetXfconfProperty(channel, prop, image); !res.IsSuccess() {
				return res
			} else {
				if msg != "" {
					msg += "\n"
				}
				msg += res.Message()
			}
		}

		if strings.HasSuffix(prop, "/image-style") {
			if res := env.SetXfconfProperty(channel, prop, "5" /*5=zoomé*/); !res.IsSuccess() {
				return res
			} else {
				if msg != "" {
					msg += "\n"
				}
				msg += res.Message()
			}
		}
	}

	return result.Success(msg)
}
