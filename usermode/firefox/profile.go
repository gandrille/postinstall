package firefox

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

// IMPORTANT
// This file provides functions that works... most of the time...
// Use at your own risks!
// TODO update it according to Firefox specs

// CreateProfile creates a profile if it is not already done
func CreateProfile() result.Result {
	if exists, _ := filesystem.Exists("~/.mozilla/firefox"); exists {
		return result.NewUnchanged("Firefox profile already exists")
	}

	if err := exec.Command("firefox", "-silent").Run(); err != nil {
		return result.NewError("Firefox profile creation failed: " + err.Error() + " - Please start Firefox once")
	}

	return result.NewCreated("Firefox profile created")
}

// GetProfileKey returns the defaul profile forlder name
func GetProfileKey() (string, error) {
	str, err := filesystem.ReadFileAsString("~/.mozilla/firefox/profiles.ini")
	if err != nil {
		return "", errors.New("Can't read profiles.ini file")
	}

	for _, line := range strings.Split(str, "\n") {
		if strings.HasPrefix(line, "Default=") {
			return strings.TrimPrefix(line, "Default="), nil
		}
	}

	return "", errors.New("Profile key not found in profiles.ini file")
}

// GetProfileFolder returns the full default profile folder path
func GetProfileFolder() (string, error) {
	key, err := GetProfileKey()
	if err != nil {
		return "", err
	}

	path := "~/.mozilla/firefox/" + key
	exists, err := filesystem.Exists(path)
	if exists {
		return path, nil
	}

	return "", errors.New("Error while checking if folder " + path + " exists")
}
