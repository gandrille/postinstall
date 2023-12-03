package main

import (
	"fmt"
	"os"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/postinstall/systemmode"
	"github.com/gandrille/postinstall/usermode"

	"github.com/gandrille/postinstall/backupmode"
)

// =========================================
// IMPORTANT
// TODO EDIT version number before releasing
// =========================================
func version() string {
	return "v23.10"
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		result.PrintRed("Missing parameters: please provide one of the following commands")
		usage()
		os.Exit(1)
	}

	switch args[0] {
	// General infos
	case "help":
		usage()
	case "version":
		fmt.Println(version())

	// System install
	case "system-install-info":
		systemmode.Describe()
	case "system-install":
		if env.Username() != "root" {
			result.PrintRed("You must be root to update the system installation!")
		} else {
			systemmode.Run()
		}

	// User install
	case "user-install-info":
		usermode.Describe()
	case "user-install":
		usermode.Run()

	// Backup
	case "user-backup-info":
		backupmode.Describe()
	case "user-backup":
		backupmode.Backup(args[1:])
	case "user-restore":
		if len(args) >= 2 {
			backupmode.Restore(args[1:])
		} else {
			usage()
		}

	default:
		usage()
	}
}

// usage prints an helper message
func usage() {
	fmt.Println("General infos")
	result.Describe("help", "                prints this help")
	result.Describe("version", "             prints version number ("+version()+")")
	fmt.Println("")
	fmt.Println("System install eases the installation of important packages")
	result.Describe("system-install-info", " describes what the installer does")
	result.Describe("system-install", "      runs the installer")
	fmt.Println("")
	fmt.Println("User install configures user desktop with nice defaults (according to me!)")
	result.Describe("user-install-info", "   describes what the installer does")
	result.Describe("user-install", "        runs the installer")
	fmt.Println("")
	fmt.Println("Backup and restore user configuration")
	result.Describe("user-backup-info", "    describes what the backup does")
	result.Describe("user-backup [file]", "  saves the user defined config to a file")
	result.Describe("user-restore file", "   restores a user defined config from a file")
	fmt.Println("")
	fmt.Println("The source code is available at https://github.com/gandrille/postinstall")
}
