package functions

import (
	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

// XfceAppFunction structure
type XfceAppFunction struct {
}

// Infos function
func (f XfceAppFunction) Infos() FunctionInfos {
	return FunctionInfos{
		Title:            "Xfce Applications",
		ShortDescription: "Configures Thunar and Xfce Terminal",
		LongDescription: `* Thunar: changes the behavior of backspace shortcut to move to parent folder
* Xfce Terminal: enable infinite history`,
	}
}

// Run function
func (f XfceAppFunction) Run() result.Result {

	// Thunar
	f1 := func() result.Result {
		updated1, err1 := filesystem.CreateOrAppendIfNotInFile("~/.config/Thunar/accels.scm", "(gtk_accel_path \"<Actions>/ThunarWindow/open-parent\" \"BackSpace\")")
		if err1 != nil {
			return result.Failure("Can't configure Thunar")
		}
		if updated1 {
			return result.Success("Thunar backspace shortcut updated to move to parent folder")
		} else {
			return result.Success("Thunar backspace shortcut already set")
		}
	}

	// Xfce Terminal config file initialization
	f2 := func() result.Result {
		return filesystem.WriteStringFile("~/.config/xfce4/terminal/terminalrc", "[Configuration]", false)
	}

	// Xfce Terminal Unlimitted scrolling
	f3 := func() result.Result {
		return filesystem.UpdateLineInFile("~/.config/xfce4/terminal/terminalrc", "ScrollingUnlimited=", "ScrollingUnlimited=TRUE", true)
	}

	// Xfce Terminal Tango color theme
	f4 := func() result.Result {
		return filesystem.UpdateLineInFile("~/.config/xfce4/terminal/terminalrc", "ColorPalette=", "ColorPalette=#000000;#cc0000;#4e9a06;#c4a000;#3465a4;#75507b;#06989a;#d3d7cf;#555753;#ef2929;#8ae234;#fce94f;#739fcf;#ad7fa8;#34e2e2;#eeeeec", true)
	}

	return execute(f.Infos().Title, f1, f2, f3, f4)
}
