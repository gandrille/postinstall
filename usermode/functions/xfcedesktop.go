package functions

import (
	_ "embed"
	"strconv"
	"strings"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
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
* Enables tile_on_move
* Power management
* Apparence / Style (colors)
* Apparence / Icon set
* Xfce Terminal configuration
* <print screen> keyboard shortcut
* <Shift>+<print screen> keyboard shortcut
* Touchpad tapping`,
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

	// Window manager / tile_on_move
	f4 := func() result.Result {
		return env.SetXfconfProperty("xfwm4", "/general/tile_on_move", "true").StandardizeMessage("Window manager / tile_on_move", "enabled")
	}

	// Power-manager Display Power Management Signaling (DPMS)
	f5 := func() result.Result {
		return env.SetXfconfProperty("xfce4-power-manager", "/xfce4-power-manager/dpms-enabled", "false").StandardizeMessage("Power-manager DPMS", "disabled")
	}

	// Power-manager blank-on-battery
	f6 := func() result.Result {
		return env.SetXfconfProperty("xfce4-power-manager", "/xfce4-power-manager/blank-on-battery", "15").StandardizeMessage("Power-manager blank-on-battery", "15min")
	}

	// Power-manager blank-on-ac
	f7 := func() result.Result {
		return env.SetXfconfProperty("xfce4-power-manager", "/xfce4-power-manager/blank-on-ac", "15").StandardizeMessage("Power-manager blank-on-ac", "15min")
	}

	// Apparence / Style (colors)
	f8 := func() result.Result {
		return env.SetXfconfProperty("xsettings", "/Net/ThemeName", "Arc").StandardizeMessage("Apparence / Style", "Arc")
	}

	// Apparence / Icon set
	f9 := func() result.Result {
		return env.SetXfconfProperty("xsettings", "/Net/IconThemeName", "bloom-classic").StandardizeMessage("Apparence / Icon set", "Bloom-classic")
	}

	// Xfce Terminal scrolling and Tango color theme
	f10 := func() result.Result {
		return filesystem.WriteBinaryFile("~/.config/xfce4/terminal/terminalrc", assetTerminalrc, true)
	}

	// 'Print' command
	f11 := func() result.Result {
		return env.SetXfconfProperty("xfce4-keyboard-shortcuts", "/commands/custom/Print", "bash -lc \"/usr/bin/xfce4-screenshooter --fullscreen --save shot-$(date +%s).png\"").StandardizeMessage("<print screen> keyboard shortcut", "full screenshot")
	}

	// '<Shift>Print' command
	f12 := func() result.Result {
		return env.SetXfconfProperty("xfce4-keyboard-shortcuts", "/commands/custom/<Shift>Print", "xfce4-screenshooter --region --clipboard").StandardizeMessage("<Shift><print screen> keyboard shortcut", "region screenshot")
	}

	// Touchpad tapping
	f13 := func() result.Result { return setTouchpadTapping("1" /* enable */) }

	return execute(f.Infos().Title, f1, f2, f3, f4, f5, f6, f7, f8, f9, f10, f11, f12, f13)
}

func setBackgroundImage(image string) result.Result {
	channel := "xfce4-desktop"
	propertySuffix := "/last-image"
	description := "Workspace image background"
	return setMultipleProperties(channel, propertySuffix, description, image)
}

func setBackgroundStyle(style string) result.Result {
	channel := "xfce4-desktop"
	propertySuffix := "/image-style"
	description := "Workspace image style"
	return setMultipleProperties(channel, propertySuffix, description, style)
}

func setTouchpadTapping(value string) result.Result {
	channel := "pointers"
	propertySuffix := "/libinput_Tapping_Enabled"
	description := "Touchpad tapping"
	return setMultipleProperties(channel, propertySuffix, description, value)
}

// TODO move this to the library
func setMultipleProperties(channel, propertySuffix, description, newValue string) result.Result {

	properties, err := env.ListXfconfProperties(channel)
	if err != nil {
		return result.NewError("Can't list XfconfProperties from " + channel)
	}

	msg := ""
	nb := 0
	nbSkipped := 0
	nbUpdated := 0
	for _, prop := range properties {
		if strings.HasSuffix(prop, propertySuffix) {
			if res := env.SetXfconfProperty(channel, prop, newValue); !res.IsSuccess() {
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
				msg += res.StandardizeMessage("Xfconf prop "+prop, newValue).Message()
			}
		}
	}

	return computeResult(nb, nbSkipped, nbUpdated, description, newValue, msg)
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

//go:embed assets/terminalrc
var assetTerminalrc []byte
