package system

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/postinstall/assets"
)

func describeAll() {
	for _, name := range getProfilsNames() {
		result.PrintInfo(name)
		println(getProfileContent(name))
	}
}

func askLines() []string {
	printWelcomeMessage()
	profils := askProfiles()
	printSelectedProfiles(profils)
	confirmInstall()
	return collectLines(profils)
}

const assetsPrefix = "resources/profils/"
const assetsSuffix = ".txt"

func askProfiles() []string {
	var profils []string

	for _, name := range getProfilsNames() {
		if includeProfile(name) {
			profils = append(profils, name)
		}
	}

	return profils
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

func getProfilsNames() []string {
	var profils []string

	for _, name := range assets.AssetNames() {
		if strings.HasPrefix(name, assetsPrefix) {
			profil := strings.TrimSuffix(strings.TrimPrefix(name, assetsPrefix), assetsSuffix)
			profils = append(profils, profil)
		}
	}
	sort.Strings(profils)

	return profils
}

func includeProfile(name string) bool {
	var input string
	for {
		result.PrintInfo("Do you want to install profil " + name + "  [ynl]")
		fmt.Scanln(&input)

		if input == "y" {
			return true
		}
		if input == "n" {
			return false
		}
		if input == "l" {
			fmt.Println(getProfileContent(name))
		}
	}
}

func getProfileContent(name string) string {
	assetname := assetsPrefix + name + assetsSuffix
	bytes, err := assets.Asset(assetname)
	if err != nil {
		log.Fatal("Can't read asset " + assetname)
	}

	// Filtering and merging consecutive empty lines into a single one
	lines := strings.Split(string(bytes), "\n")
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

func printSelectedProfiles(profils []string) {
	println("")
	println("You have selected the following profils:")
	for _, p := range profils {
		fmt.Printf(" * %s\n", p)
	}
}

func confirmInstall() {
	println("")
	println("In order to run the installation process, press enter key.")
	println("To cancel, press Ctrl+C keys.")
	var input string
	fmt.Scanln(&input)
}

func collectLines(profils []string) []string {
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
