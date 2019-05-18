# Postinstall

Helper tool to configure an [(X)Ubuntu](https://xubuntu.org/) operating system right after it has been installed.

→ You can [download the latest pre-built release](https://github.com/gandrille/postinstall/releases/latest), or [build it yourself](#build).


## Usage

Syntax : `postinstall <command>`

Here is the output produced with `postinstall help` 

```
General infos
help                 prints this help
version              prints version number (v1.0)

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

```
cd ${GOPATH:-~/go}
go get github.com/gandrille/postinstall/...
src/github.com/gandrille/postinstall/update-assets
go install src/github.com/gandrille/postinstall/postinstall.go 
```


## Changelog

**[v1.0](../../releases/tag/v1.0)** Designed for XUbuntu 18.04 LTS (Bionic Beaver)

**[v0.9](../../releases/tag/v0.9)** Pre-release


## License

This project is released under the
[Apache 2.0 license](https://www.apache.org/licenses/LICENSE-2.0.html).


**Dependencies**
* [github.com/fatih/color](https://github.com/fatih/color/): [MIT](https://github.com/fatih/color/blob/master/LICENSE.md)
* [github.com/go-bindata/go-bindata](https://github.com/go-bindata/go-bindata/): [CC0](https://github.com/go-bindata/go-bindata/blob/master/LICENSE)
* [github.com/gandrille/go-commons](https://github.com/gandrille/go-commons): [Apache 2.0 license](https://github.com/gandrille/go-commons/blob/master/LICENSE.txt)
