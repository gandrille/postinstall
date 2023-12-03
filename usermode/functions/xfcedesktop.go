package functions

import (
	"strconv"
	"strings"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/postinstall/assets"
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
	return assets.WriteAsset("resources/terminalrc", "~/.config/xfce4/terminal/terminalrc", true)

	/*
		strpair.New("FontName", "DejaVu Sans Mono 9"),
		strpair.New("MiscDefaultGeometry", "120x32"),
		strpair.New("ColorPalette", "rgb(0,0,0);rgb(213,72,72);rgb(184,214,140);rgb(196,160,0);rgb(52,101,164);rgb(117,80,123);rgb(6,152,154);rgb(211,215,207);rgb(85,87,83);rgb(210,61,61);rgb(160,207,93);rgb(252,233,79);rgb(115,159,207);rgb(173,127,168);rgb(52,226,226);rgb(238,238,236)"),
		strpair.New("ColorBackground", "#090909090909"),
		strpair.New("ColorForeground", "#eeeeeeeeeeee"),
		strpair.New("ColorSelection", "#16163b3b5959"),
		strpair.New("ColorSelectionUseDefault", "FALSE"),
		strpair.New("ColorBold", "#ffffffffffff"),
		strpair.New("TabActivityColor", "#0f0e49499999"),
		strpair.New("ScrollingUnlimited", "TRUE"),
	*/
}
