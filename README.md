# Postinstall

Helper tool to configure an [(X)Ubuntu](https://xubuntu.org/) operating system right after it has been installed.


## Usage

Syntax : `postinstall <command>`

Here is the output produced with `postinstall help` 

```
help                prints this help
system-install-info prints what the system installation does
system-install      eases the installation of important packages
user-install-info   prints what the user installation does
user-install        runs the user installation
user-backup-info    prints what user backup does
user-backup [file]  saves the user defined config to a file
user-restore file   restores a user defined config from a file
```


## Build

```
go get github.com/gandrille/postinstall/...
src/github.com/gandrille/postinstall/update-assets
go install src/github.com/gandrille/postinstall/postinstall.go 
```


## Changelog

**v1.0** Designed for XUbuntu 18.04

**[v0.9](../../releases/tag/v0.9)** Pre-release


## Known Bugs

Since installing `xfce4-power-manager-plugin` fails, the `xfce-plugins` statement also fails.


## TODO

Choose a license complient with third parties.

