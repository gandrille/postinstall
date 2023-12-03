package backupmode

import (
	"archive/zip"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/misc"
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/go-commons/strpair"
	"github.com/gandrille/go-commons/zipfile"
)

// A source is "something" which can be backup and restore.
// It also has a describe function.
type source interface {
	Backup(writer *zip.Writer, kvstore *[]strpair.StrPair) result.Result
	Restore(zip zipfile.ZipFile, store []strpair.StrPair) result.Result
	Describe()
}

func getSources() []source {
	var sources []source

	// standard files
	sources = append(sources, FileSource{"~/.netrc", "netrc"})
	sources = append(sources, FileSource{"~/.sendxmpprc", "sendxmpprc"})
	sources = append(sources, FileSource{"~/.bashrc-perso", "bashrc-perso"})
	sources = append(sources, FileSource{"~/.gitconfig", "gitconfig"})
	sources = append(sources, FileSource{"~/.m2/settings.xml", "settings.xml"})
	sources = append(sources, FileSource{"~/.config/onedrive/config", "onedrive-config.txt"})

	// standard dirs
	sources = append(sources, DirSource{"~/.ssh", "ssh"})
	sources = append(sources, DirSource{"~/.aws", "aws"})
	sources = append(sources, DirSource{"~/.config/imapsync", "imapsync"})
	sources = append(sources, DirSource{"~/.config/filezilla", "filezilla"})
	sources = append(sources, DirSource{"~/.config/GIMP/2.10/scripts/", "gimp-scripts-2.10"})

	// custom dirs
	// DEPRECATED sources = append(sources, DirSource{"~/.config/auth", "auth"})
	// DEPRECATED sources = append(sources, DirSource{"~/.config/ics", "ics"})
	sources = append(sources, DirSource{"~/.config/card", "card"})
	sources = append(sources, DirSource{"~/.config/quick-links", "quick-links"})
	sources = append(sources, DirSource{"~/.config/webdav", "webdav"})
	sources = append(sources, DirSource{"~/.config/listen-xmpp", "listen-xmpp"})

	// homemade backup/restore procedure
	sources = append(sources, UnisonSource{})
	sources = append(sources, CrontabSource{})
	sources = append(sources, DpkgSource{})
	sources = append(sources, XfceSource{})
	sources = append(sources, XfcePropertySource{})

	// homemade backup only (for info purpose)
	sources = append(sources, FstabSource{})
	sources = append(sources, LsblkSource{})

	return sources
}

func getValue(store []strpair.StrPair, key string) (string, error) {
	for _, pair := range store {
		if pair.Str1() == key {
			return pair.Str2(), nil
		}
	}
	return "", errors.New("Pair with key " + key + " not found")
}

/* ========== */
/* FileSource */
/* ========== */

type FileSource struct {
	filePath string
	zipPath  string
}

func (src FileSource) Backup(writer *zip.Writer, kvstore *[]strpair.StrPair) result.Result {
	return ProcessFile(src.filePath, src.zipPath, writer)
}

func (src FileSource) Restore(zip zipfile.ZipFile, store []strpair.StrPair) result.Result {
	if zip.HasFile(src.zipPath) {
		return zip.GetFile(src.zipPath).Write(src.filePath)
	}
	return result.NewError(src.zipPath + " not found in zip file (skipped)")
}

func (src FileSource) Describe() {
	result.Describe(src.filePath, "file will be stored into "+src.zipPath)
}

/* ========= */
/* DirSource */
/* ========= */

type DirSource struct {
	dirPath string
	zipPath string
}

func (src DirSource) Backup(writer *zip.Writer, kvstore *[]strpair.StrPair) result.Result {
	return ProcessDir(src.dirPath, src.zipPath, writer)
}

func (src DirSource) Restore(zip zipfile.ZipFile, store []strpair.StrPair) result.Result {
	files := zip.FilesStartingWith(src.zipPath)
	return copyFilesToFolder(files, src.zipPath, src.dirPath)
}

func (src DirSource) Describe() {
	result.Describe(src.dirPath, "folder will be stored into "+src.zipPath)
}

/* ============ */
/* UnisonSource */
/* ============ */

type UnisonSource struct {
}

type matchUnison struct {
}

// TODO use a pure function instead ?
func (f matchUnison) Match(element srcdst) bool {
	return strings.HasSuffix(element.src, ".prf")
}

func (src UnisonSource) Backup(writer *zip.Writer, kvstore *[]strpair.StrPair) result.Result {
	return ProcessDirFiltered("~/.unison", "unison", matchUnison{}, writer)
}

func (src UnisonSource) Restore(zip zipfile.ZipFile, store []strpair.StrPair) result.Result {
	// we have backup only the necessary... so we can restore everything
	files := zip.FilesStartingWith("unison")
	return copyFilesToFolder(files, "unison", "~/.unison")
}

func (src UnisonSource) Describe() {
	result.Describe("unison", "~/.unison/*.prf files will be stored into unison")
}

/* =========== */
/* FstabSource */
/* =========== */

type DpkgSource struct {
}

func (src DpkgSource) Backup(writer *zip.Writer, kvstore *[]strpair.StrPair) result.Result {
	if out, err := exec.Command("/usr/bin/dpkg", "--list").Output(); err != nil {
		return result.NewError("dpkg info not available: " + err.Error())
	} else {
		return ProcessBytes("dpkg", "system/dpkg-list.txt", out, writer)
	}
}

func (src DpkgSource) Restore(zip zipfile.ZipFile, store []strpair.StrPair) result.Result {
	return result.NewInfo("dpkg list of installed packages will NOT be restored (It's for info only!)")
}

func (src DpkgSource) Describe() {
	result.Describe("dpkg", "Installed packages list will be stored into system/dpkg-list.txt")
}

/* ============= */
/* CrontabSource */
/* ============= */

type CrontabSource struct {
}

func (src CrontabSource) Backup(writer *zip.Writer, kvstore *[]strpair.StrPair) result.Result {
	ct := getCurrentCrontab()

	if ct == "" {
		return result.NewUnchanged("NO user crontab defined")
	}
	return ProcessBytes("crontab", "crontab", []byte(ct), writer)
}

func (src CrontabSource) Restore(zip zipfile.ZipFile, store []strpair.StrPair) result.Result {
	zipCT, err := getCrontabInZip(zip)
	if err != nil {
		return result.NewError(err.Error())
	}

	curCT := getCurrentCrontab()
	if zipCT == "" {
		if curCT != "" {
			return result.NewUnchanged("crontab not available in zip file (current crontab left UNCHANGED)")
		}
		return result.NewUnchanged("crontab not available in zip file")
	}
	return restoreCrontab(curCT, zipCT)
}

func getCurrentCrontab() string {
	if out, err := exec.Command("/usr/bin/crontab", "-l").Output(); err != nil {
		return ""
	} else {
		return string(out)
	}
}

func getCrontabInZip(zip zipfile.ZipFile) (string, error) {
	file := zip.GetFile("crontab")
	if file != nil {
		return file.StringContent()
	}
	return "", nil
}

func restoreCrontab(old, new string) result.Result {
	if old == new {
		return result.NewUnchanged("crontab is already up to date")
	}

	cmd := exec.Command("/usr/bin/crontab", "-")
	res := misc.RunCmdStdIn("crontab", new, cmd)
	if res.IsSuccess() {
		return result.NewUpdated("User crontab restored")
	} else {
		return result.NewError("User crontab revert failed")
	}
}

func (src CrontabSource) Describe() {
	result.Describe("crontab", "user crontab will be stored into crontab")
}

/* =========== */
/* FstabSource */
/* =========== */

type FstabSource struct {
}

func (src FstabSource) Backup(writer *zip.Writer, kvstore *[]strpair.StrPair) result.Result {
	return ProcessFile("/etc/fstab", "system/fstab", writer)
}

func (src FstabSource) Restore(zip zipfile.ZipFile, store []strpair.StrPair) result.Result {
	return result.NewInfo("fstab will NOT be restored (You are not root!)")
}

func (src FstabSource) Describe() {
	result.Describe("fstab", "system fstab will be stored into fstab")
}

/* ===== */
/* lsblk */
/* ===== */

type LsblkSource struct {
}

func (src LsblkSource) Backup(writer *zip.Writer, kvstore *[]strpair.StrPair) result.Result {
	if out, err := exec.Command("/bin/lsblk", "-o", "NAME,FSTYPE,SIZE,LABEL,UUID,MOUNTPOINT").Output(); err != nil {
		return result.NewError("lsblk info not available: " + err.Error())
	} else {
		return ProcessBytes("lsblk", "system/lsblk", out, writer)
	}
}

func (src LsblkSource) Restore(zip zipfile.ZipFile, store []strpair.StrPair) result.Result {
	return result.NewInfo("lsblk will NOT be restored (it's only for information purpose)")
}

func (src LsblkSource) Describe() {
	result.Describe("lsblk", "list of all block mount points (partition table)")
}

/* ========== */
/* XfceSource */
/* ========== */

type XfceSource struct {
}

func (src XfceSource) Backup(writer *zip.Writer, kvstore *[]strpair.StrPair) result.Result {
	allskipped := true

	if res := ProcessDir("~/.config/xfce4/panel", "xfce4/panel", writer); !res.IsSuccess() {
		return res
	} else {
		allskipped = allskipped && res.IsUnchanged()
	}

	if res := ProcessFile("~/.config/xfce4/xfconf/xfce-perchannel-xml/xfce4-panel.xml", "xfce4/xfce4-panel.xml", writer); !res.IsSuccess() {
		return res
	} else {
		allskipped = allskipped && res.IsUnchanged()
	}

	if allskipped {
		return result.NewUnchanged("Xfce configuration is NOT available")
	} else {
		return result.NewUpdated("Xfce configuration added")
	}
}

func (src XfceSource) Restore(zip zipfile.ZipFile, store []strpair.StrPair) result.Result {
	allskipped := true

	files := zip.FilesStartingWith("xfce4/panel")
	if res := copyFilesToFolder(files, "xfce4/panel", "~/.config/xfce4/panel"); !res.IsSuccess() {
		return res
	} else {
		allskipped = allskipped && res.IsUnchanged()
	}

	if zip.HasFile("xfce4/xfce4-panel.xml") {
		if res := zip.GetFile("xfce4/xfce4-panel.xml").Write("~/.config/xfce4/xfconf/xfce-perchannel-xml/xfce4-panel.xml"); !res.IsSuccess() {
			return result.NewError("xfce4/xfce4-panel.xml update failed")
		} else {
			allskipped = allskipped && res.IsUnchanged()
		}
	}

	if allskipped {
		return result.NewUnchanged("Xfce config already up to date")
	}
	return result.NewUpdated("Xfce config restored")
}

func (src XfceSource) Describe() {
	result.Describe("Xfce panel", "configuration files will be stored into xfce4")
}

/* ================== */
/* XfcePropertySource */
/* ================== */

type XfcePropertySource struct {
}

func (src XfcePropertySource) Backup(writer *zip.Writer, kvstore *[]strpair.StrPair) result.Result {

	val, err := env.ReadXfconfProperty("xfwm4", "/general/workspace_count")
	if err != nil {
		return result.NewError("Can't read workspace count")
	}
	*kvstore = append(*kvstore, strpair.New("workspace_count", val))

	return result.NewUpdated("Xfce Key/Value store properties registered")
}

func (src XfcePropertySource) Restore(zip zipfile.ZipFile, store []strpair.StrPair) result.Result {

	val, err := getValue(store, "workspace_count")
	if err != nil {
		return result.NewUnchanged("No update made")
	}

	return env.SetXfconfProperty("xfwm4", "/general/workspace_count", val)
}

func (src XfcePropertySource) Describe() {
	result.Describe("Xfce properties", "Configuration properties will be added to Key/Value store")
}

/* ================= */
/* COMMON / INTERNAL */
/* ================= */

func copyFilesToFolder(files []zipfile.ZipElement, zippath, destinationPath string) result.Result {
	if len(files) == 0 {
		return result.NewUnchanged("No file to restore to " + destinationPath)
	}

	var inError []string
	skipped := 0
	for _, file := range files {
		if file.IsValid() {
			if !strings.HasPrefix(file.Name(), zippath) {
				inError = append(inError, file.Name())
			} else {
				dst := strings.TrimSuffix(destinationPath, "/") + "/" + strings.TrimPrefix(strings.TrimPrefix(file.Name(), zippath), "/")
				res := file.Write(dst)
				if !res.IsSuccess() {
					inError = append(inError, file.Name())
				} else {
					if res.IsUnchanged() {
						skipped++
					}
				}
			}
		} else {
			inError = append(inError, file.Name())
		}
	}

	if len(inError) == 0 {
		if len(files) == skipped {
			return result.NewUnchanged(fileString(len(files)) + " with already the expected content into " + destinationPath)
		} else if skipped == 0 {
			return result.NewUpdated(fileString(len(files)) + " written to " + destinationPath)
		} else {
			return result.NewUpdated(fileString(len(files)-skipped) + " written to " + destinationPath + " (" + fileString(skipped) + " with already the expected content)")
		}
	}

	return result.NewError(fileString(len(inError)) + " in error : " + strings.Join(inError, ","))
}

func fileString(nb int) string {
	if nb == 0 {
		return "no file"
	}
	if nb == 1 {
		return "1 file"
	}
	return strconv.Itoa(nb) + " files"
}
