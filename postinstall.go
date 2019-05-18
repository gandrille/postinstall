package main

import (
	"os"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/postinstall/system"
	"github.com/gandrille/postinstall/user/backup"
	"github.com/gandrille/postinstall/user/initialize"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		result.PrintRed("Missing parameters")
		usage()
		os.Exit(1)
	}

	switch args[0] {

	case "help":
		usage()

	case "system-install-info":
		system.Describe()
	case "system-install":
		if env.Username() != "root" {
			result.PrintRed("You must be root to update the system installation!")
		} else {
			system.Run()
		}

	case "user-install-info":
		initialize.Describe()
	case "user-install":
		initialize.Run()

	case "user-backup-info":
		backup.Describe()
	case "user-backup":
		backup.Backup(args[1:])
	case "user-restore":
		if len(args) >= 2 {
			backup.Restore(args[1:])
		} else {
			usage()
		}

	default:
		usage()
	}
}

// usage prints an helper message
func usage() {
	result.Describe("help", "               prints this help")
	result.Describe("system-install-info", "prints what the system installation does")
	result.Describe("system-install", "     eases the installation of important packages")
	result.Describe("user-install-info", "  prints what the user installation does")
	result.Describe("user-install", "       runs the user installation")
	result.Describe("user-backup-info", "   prints what user backup does")
	result.Describe("user-backup [file]", " saves the user defined config to a file")
	result.Describe("user-restore file", "  restores a user defined config from a file")
}
