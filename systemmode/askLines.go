package systemmode

import (
	"fmt"
	"strings"

	"github.com/gandrille/go-commons/result"
)

func describeAll() {
	if availableProfils, err := GetProfils(); err != nil {
		result.PrintError(err.Error())
	} else {
		for _, profile := range availableProfils {
			result.PrintInfo(profile.Name)
			println(getProfileContent(profile))
		}
	}
}

func askLines() ([]string, error) {
	printWelcomeMessage()
	profils, err := askProfiles()
	if err != nil {
		return []string{}, err
	}
	printSelectedProfiles(profils)
	confirmInstall()
	return collectLines(profils), nil
}

func printWelcomeMessage() {
	println(`First, you need to select the profils you want to install.
After each profil, please type:
 - y --> YES, install profil
 - n --> NO, do NOT install
 - l --> LIST profil content, before the question is asked again
At the end of the profil selection, you will be asked to confirm your selection in order to trigger the installation process. 
`)
}

func askProfiles() ([]Profile, error) {
	availableProfils, err := GetProfils()
	if err != nil {
		return []Profile{}, err
	}

	selectedProfils := []Profile{}
	for _, availableProfile := range availableProfils {
		if includeProfile(availableProfile) {
			selectedProfils = append(selectedProfils, availableProfile)
		}
	}

	return selectedProfils, nil
}

func includeProfile(profile Profile) bool {
	var input string
	for {
		result.PrintInfo("Do you want to install profil " + profile.Name + "  [ynl]")
		fmt.Scanln(&input)

		if input == "y" {
			return true
		}
		if input == "n" {
			return false
		}
		if input == "l" {
			fmt.Println(getProfileContent(profile))
		}
	}
}

func getProfileContent(profile Profile) string {
	lines := strings.Split(profile.Contents, "\n")
	var filtered []string
	lastLineEmpty := true
	for _, line := range lines {
		if !strings.HasPrefix(line, "##") {
			if !lastLineEmpty || strings.Trim(line, " ") != "" {
				filtered = append(filtered, line)
			}
			lastLineEmpty = (line == "")
		}
	}

	return strings.Join(filtered, "\n")
}

func printSelectedProfiles(profils []Profile) {
	println("")
	println("You have selected the following profils:")
	for _, p := range profils {
		fmt.Printf(" * %s\n", p.Name)
	}
}

func confirmInstall() {
	println("")
	println("In order to run the installation process, press enter key.")
	println("To cancel, press Ctrl+C keys.")
	var input string
	fmt.Scanln(&input)
}

func collectLines(profils []Profile) []string {
	var lines []string

	for _, profil := range profils {
		profileContent := strings.Split(getProfileContent(profil), "\n")
		for _, line := range profileContent {
			line := strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				lines = append(lines, line)
			}
		}
	}
	return lines
}
