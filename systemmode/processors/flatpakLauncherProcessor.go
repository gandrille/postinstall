package processors

import (
	"os"
	"strings"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/postinstall/assets"
)

type flatpakLauncher struct {
}

func (e flatpakLauncher) Key() string {
	return "flatpak-launcher"
}

func (e flatpakLauncher) Describe(args []string) string {
	return "Creates a launcher for a flatpak application: " + strings.Join(args, ",")
}

func (e flatpakLauncher) Run(args []string) result.Result {
	filePath := args[0]
	flatpakId := args[1]
	expectedContents := strings.Replace(assets.AssetAsString("resources/flatpak-launcher"), "FLATPAK_ID", flatpakId, -1)

	// Updating file
	res := filesystem.WriteStringFile(filePath, expectedContents, true)
	if res.IsError() {
		return res
	}
	updated := !res.IsUnchanged()

	// Check if perms are fine
	if fileInfo, err := os.Lstat(filePath); err != nil {
		return result.NewError("Can't retreive " + filePath + " file permissions")
	} else {
		if fileInfo.Mode().Perm().String() != "-rwxr-xr-x" {
			updated = true
			if err := os.Chmod(filePath, 0755); err != nil {
				return result.NewError("Can't update " + filePath + " file permissions")
			}
		}
	}

	if updated {
		return result.NewUpdated("flatpak launcher " + filePath + " updated")
	} else {
		return result.NewUnchanged("flatpak launcher " + filePath + " already well defined")
	}
}
