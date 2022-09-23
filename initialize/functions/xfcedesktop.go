package functions

import (
	"strconv"
	"strings"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/ini"
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/go-commons/strpair"
)

// XfceDesktopFunction structure
type XfceDesktopFunction struct {
}

// Infos function
func (f XfceDesktopFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "Xfce Desktop",
		ShortDescription: "Updates Xfce desktop configuration",
		LongDescription: `* Workspace image background
* Workspace image background style
* Window manager theme
* Apparence / Style (colors)
* Apparence / Icon set
* <print screen> keyboard shortcut
* Thunar: changes the behavior of backspace shortcut to move to parent folder
* Xfce Terminal: enable infinite history and sets the color palette`,
	}
}

// Run function
func (f XfceDesktopFunction) Run() result.Result {

	// Workspace image background
	f1 := func() result.Result { return setBackgroundImage("/usr/share/xfce4/backdrops/xubuntu-yakkety.png") }

	// Workspace image background
	f2 := func() result.Result { return setBackgroundStyle("5" /* zoom */) }

	// Window manager theme
	f3 := func() result.Result {
		return env.SetXfconfProperty("xfwm4", "/general/theme", "Arc").StandardizeMessage("Window manager theme", "Arc")
	}

	// Apparence / Style (colors)
	f4 := func() result.Result {
		return env.SetXfconfProperty("xsettings", "/Net/ThemeName", "Arc").StandardizeMessage("Apparence / Style (colors)", "Arc")
	}

	// Apparence / Icon set
	f5 := func() result.Result {
		return env.SetXfconfProperty("xsettings", "/Net/IconThemeName", "bloom-classic").StandardizeMessage("Apparence / Icon set", "Bloom-classic")
	}

	// 'Print' command
	f6 := func() result.Result {
		return env.SetXfconfProperty("xfce4-keyboard-shortcuts", "/commands/custom/Print", "bash -lc \"/usr/bin/xfce4-screenshooter --fullscreen --save shot-$(date +%s).png\"").StandardizeMessage("<print screen> keyboard shortcut", "full screenshot")
	}

	// '<Shift>Print' command
	f7 := func() result.Result {
		return env.SetXfconfProperty("xfce4-keyboard-shortcuts", "/commands/custom/<Shift>Print", "xfce4-screenshooter --region --clipboard").StandardizeMessage("<Shift><print screen> keyboard shortcut", "region screenshot")
	}

	// Xfce Terminal scrolling and Tango color theme
	f8 := configureXfceTerminal

	return execute(f.Infos().Title, f1, f2, f3, f4, f5, f6, f7, f8)
}

func setBackgroundImage(image string) result.Result {
	channel := "xfce4-desktop"

	properties, err := env.ListXfconfProperties(channel)
	if err != nil {
		return result.NewError("Can't list XfconfProperties from " + channel)
	}

	msg := ""
	nb := 0
	nbSkipped := 0
	nbUpdated := 0
	for _, prop := range properties {
		if strings.HasSuffix(prop, "/last-image") {
			if res := env.SetXfconfProperty(channel, prop, image); !res.IsSuccess() {
				return res
			} else {
				if msg != "" {
					msg += "\n"
				}
				nb++
				if res.IsUnchanged() {
					nbSkipped++
				} else if res.IsUpdated() {
					nbUpdated++
				}
				msg += res.StandardizeMessage("Xfconf prop "+prop, image).Message()
			}

		}
	}

	return computeResult(nb, nbSkipped, nbUpdated, "Workspace image background", image, msg)
}

func setBackgroundStyle(style string) result.Result {
	channel := "xfce4-desktop"

	properties, err := env.ListXfconfProperties(channel)
	if err != nil {
		return result.NewError("Can't list XfconfProperties from " + channel)
	}

	msg := ""
	nb := 0
	nbSkipped := 0
	nbUpdated := 0
	for _, prop := range properties {
		if strings.HasSuffix(prop, "/image-style") {
			if res := env.SetXfconfProperty(channel, prop, style); !res.IsSuccess() {
				return res
			} else {
				if msg != "" {
					msg += "\n"
				}
				nb++
				if res.IsUnchanged() {
					nbSkipped++
				} else if res.IsUpdated() {
					nbUpdated++
				}
				msg += res.StandardizeMessage("Xfconf prop "+prop, "5").Message()
			}
		}
	}

	return computeResult(nb, nbSkipped, nbUpdated, "Workspace image style", style, msg)
}

func computeResult(nb, nbSkipped, nbUpdated int, name, value, defaultMessage string) result.Result {
	if nbSkipped == nb {
		return result.NewUnchanged(name + " already set in all " + strconv.Itoa(nbSkipped) + " props with value " + value)
	} else if nbUpdated == nb {
		return result.NewUpdated(name + " updated in " + strconv.Itoa(nbUpdated) + " props with value " + value)
	} else if nbUpdated+nbSkipped == nb {
		return result.NewUpdated(name + " updated in " + strconv.Itoa(nbUpdated) + " props with value " + value + " (" + strconv.Itoa(nbSkipped) + " props already set)")
	} else {
		return result.NewUpdated(defaultMessage)
	}
}

func configureXfceTerminal() result.Result {
	configFile := "~/.config/xfce4/terminal/terminalrc"
	section := "Configuration"

	props := []strpair.StrPair{
		strpair.New("FontName", "DejaVu Sans Mono 9"),
		strpair.New("ColorForeground", "#eeeeee"),
		strpair.New("ColorBackground", "#000000"),
		strpair.New("ColorPalette", "#000000;#cc0000;#4e9a06;#c4a000;#3465a4;#75507b;#06989a;#d3d7cf;#555753;#ef2929;#8ae234;#fce94f;#739fcf;#ad7fa8;#34e2e2;#eeeeec"),
		strpair.New("ColorSelection", "#163b59"),
		strpair.New("ColorSelectionUseDefault", "FALSE"),
		strpair.New("ColorCursor", "#0f4999"),
		strpair.New("ColorBold", "#ffffff"),
		strpair.New("TabActivityColor", "#0f4999"),
		strpair.New("ScrollingLines", "30000"),
	}

	allskipped := true
	for _, prop := range props {
		if updated, err := ini.SetValue(configFile, section, prop.Str1(), prop.Str2(), false, false); err != nil {
			return result.NewError("Error while updating " + configFile + ": " + err.Error())
		} else {
			allskipped = allskipped && !updated
		}
	}

	if allskipped {
		return result.NewUnchanged("Xfce terminal config is already up to date")
	} else {
		return result.NewUpdated("Xfce terminal config updated")
	}
}
