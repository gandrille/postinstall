# Postinstall

Helper tool to configure an [(X)Ubuntu](https://xubuntu.org/) operating system right after it has been installed.

→ You can [download the latest pre-built release](https://github.com/gandrille/postinstall/releases/latest), or [build it yourself](#build).


## Usage

Syntax : `postinstall <command>`

Here is the output produced with `postinstall help` 

```
General infos
help                 prints this help
version              prints version number (v24.04)

System install eases the installation of important packages
system-install-info  describes what the installer does
system-install       runs the installer

User install configures user desktop with nice defaults (according to me!)
user-install-info    describes what the installer does
user-install         runs the installer

Backup and restore user configuration
user-backup-info     describes what the backup does
user-backup [file]   saves the user defined config to a file
user-restore file    restores a user defined config from a file

The source code is available at https://github.com/gandrille/postinstall
```


## Build

A go version >= 1.21 is required. You can check your go version using `go version` command.

```
git clone git@github.com:gandrille/postinstall.git
cd postinstall
./update-assets
go install postinstall.go
${GOPATH:-~/go}/bin/postinstall help
```

## License

This project is released under the
[Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0.html).
