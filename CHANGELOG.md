Changes for postinstall
=======================

## [v24.04](../../releases/tag/v24.04)
Designed for XUbuntu 24.04 LTS () - minor update
* System install mode
   - Minor updates
* User install mode
   - Minor updates
* Backup/Restore user config mode
   - Minor updates

## [v23.10](../../releases/tag/v23.10)
Designed for XUbuntu 23.10 (Mantic Minotaur) - next LTS preparation
* System install mode
   - New processors: `snap-install`, `flatpak-add-remote`, `flatpak-install`, `drawio-install`, `vscode-install`, `configure-unattendedUpgrade`
   - Improved `configure-imagemagick` processor
   - Profiles reordered and updated, `snap-and-flatpak` profile added
* User install mode
   - SDKman! updated with Java 17 and 21 (Temurin flavor)
   - XFCE desktop updated (power management, XFCE terminal,...)
   - Zim Web template removed
   - Firefox configuration removed
* Backup/Restore user config mode
   - Minor updates
* Source code
   - Now requires go version >= 1.21 (uses `slices` package)
   - Packages layout updated
   - `//go:embed` is now used instead of the brave old legacy `go-bindata`

## [v1.4](../../releases/tag/v1.4)
Designed for XUbuntu 22.04 LTS (Jammy Jellyfish) - minor update
* System install mode
   - New processor : `configure-imagemagick`
   - Profiles updated
* User install mode
   - Web and mail exo configuration removed
   - SDKman! updated with Java 8 and 11
   - Python package installer
   - `/commands/custom/<Shift>Print` shortcut
* Backup/Restore user config mode
   - Minor updates
* Source code
   - `go-bindata` now has to be installed using `apt`

## [v1.3](../../releases/tag/v1.3)
Designed for XUbuntu 20.04 LTS (Focal Fossa) - minor update
* System install mode
   - Profiles updated
* User install mode
   - Zim configuration updated
   - Firefox configuration updated
   - `~/.hidden` file initialization
* Backup/Restore user config mode
   - Dpkg packages list is now dumped inside the zip file for info purpose
   - Minor updates

## [v1.2](../../releases/tag/v1.2)
Designed for XUbuntu 20.04 LTS (Focal Fossa)
* System install mode
   - Removed processors : `configure-timezone`, `configure-systemd-logind`
   - Profiles updated
* User install mode
   - Default `.m2/settings.xml` cleaned
   - Removed screensaver initialization

## [v1.1](../../releases/tag/v1.1)
First pre-release BEFORE 20.04 LTS (Focal Fossa)
* System install mode
   - New processors: `configure-timezone`, `configure-systemd-timesyncd`, `configure-systemd-logind`
   - Renamed processor: `fuse-conf` into `configure-fuse`

* User install mode
   - zim simple web template added
   - Default `.m2/settings.xml` updated
   - `SdkManFunction`, `SSHFunction`, `FirefoxFunction`, `ScreensaverAndLockFunction` added

## [v1.0](../../releases/tag/v1.0)
Designed for XUbuntu 18.04 LTS (Bionic Beaver)

## [v0.9](../../releases/tag/v0.9)
Pre-release
