package backup

import (
	"archive/zip"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gandrille/go-commons/misc"
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/go-commons/zipfile"
)

// A source is "something" which can be backup and restore.
// It also has a describe function.
type source interface {
	Backup(writer *zip.Writer) result.Result
	Restore(zip zipfile.ZipFile) result.Result
	Describe()
}

func getSources() []source {
	var sources []source

	sources = append(sources, FileSource{"~/.netrc", "netrc"})
	sources = append(sources, FileSource{"~/.gitconfig", "gitconfig"})
	sources = append(sources, FileSource{"~/.m2/settings.xml", "settings.xml"})

	sources = append(sources, DirSource{"~/.ssh", "ssh"})
	sources = append(sources, DirSource{"~/.aws", "aws"})
	sources = append(sources, DirSource{"~/.config/auth", "auth"})
	sources = append(sources, DirSource{"~/.config/ics", "ics"})
	sources = append(sources, DirSource{"~/.config/quick-links", "quick-links"})
	sources = append(sources, DirSource{"~/.config/webdav", "webdav"})

	sources = append(sources, UnisonSource{})
	sources = append(sources, CrontabSource{})
	sources = append(sources, FstabSource{})
	sources = append(sources, LsblkSource{})

	sources = append(sources, XfceSource{})

	return sources
}

/* ========== */
/* FileSource */
/* ========== */

type FileSource struct {
	filePath string
	zipPath  string
}

func (src FileSource) Backup(writer *zip.Writer) result.Result {
	return ProcessFile(src.filePath, src.zipPath, writer)
}

func (src FileSource) Restore(zip zipfile.ZipFile) result.Result {
	if zip.HasFile(src.zipPath) {
		return zip.GetFile(src.zipPath).Write(src.filePath)
	}
	return result.Success(src.zipPath + " not found in zip file (skipped)")
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

func (src DirSource) Backup(writer *zip.Writer) result.Result {
	return ProcessDir(src.dirPath, src.zipPath, writer)
}

func (src DirSource) Restore(zip zipfile.ZipFile) result.Result {
	files := zip.FilesStartingWith(src.zipPath)
	return copyFilesToFolder(files, src.dirPath)
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

func (src UnisonSource) Backup(writer *zip.Writer) result.Result {
	return ProcessDirFiltered("~/.unison", "unison", matchUnison{}, writer)
}

func (src UnisonSource) Restore(zip zipfile.ZipFile) result.Result {
	// we have backup only the necessary... so we can restore everything
	files := zip.FilesStartingWith("unison")
	return copyFilesToFolder(files, "~/.unison")
}

func (src UnisonSource) Describe() {
	result.Describe("unison", "~/.unison/*.prf files will be stored into unison")
}

/* ============= */
/* CrontabSource */
/* ============= */

type CrontabSource struct {
}

func (src CrontabSource) Backup(writer *zip.Writer) result.Result {
	ct := getCurrentCrontab()

	if ct == "" {
		return result.Success("Crontab is NOT available")
	}
	return ProcessBytes("crontab", "crontab", []byte(ct), writer)
}

func (src CrontabSource) Restore(zip zipfile.ZipFile) result.Result {
	zipCT, err := getCrontabInZip(zip)
	if err != nil {
		return result.Failure(err.Error())
	}
	if zipCT == "" {
		return result.Success("crontab not found in zip file (skipped)")
	}
	curCT := getCurrentCrontab()
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
		return result.Success("Crontab already up to date")
	}

	cmd := exec.Command("/usr/bin/crontab", "-")
	return misc.RunCmdStdIn("crontab", new, cmd)
}

func (src CrontabSource) Describe() {
	result.Describe("crontab", "User crontab will be stored into crontab")
}

/* =========== */
/* FstabSource */
/* =========== */

type FstabSource struct {
}

func (src FstabSource) Backup(writer *zip.Writer) result.Result {
	return ProcessFile("/etc/fstab", "system/fstab", writer)
}

func (src FstabSource) Restore(zip zipfile.ZipFile) result.Result {
	return result.Success("fstab will NOT be restored (You are not root!)")
}

func (src FstabSource) Describe() {
	result.Describe("fstab", "System fstab will be stored into fstab")
}

/* ===== */
/* lsblk */
/* ===== */

type LsblkSource struct {
}

func (src LsblkSource) Backup(writer *zip.Writer) result.Result {
	if out, err := exec.Command("/bin/lsblk", "-o", "NAME,FSTYPE,SIZE,LABEL,UUID,MOUNTPOINT").Output(); err != nil {
		return result.Failure("lsblk info not available: " + err.Error())
	} else {
		return ProcessBytes("lsblk", "system/lsblk", out, writer)
	}
}

func (src LsblkSource) Restore(zip zipfile.ZipFile) result.Result {
	return result.Success("lsblk will NOT be restored (it's only for information purpose)")
}

func (src LsblkSource) Describe() {
	result.Describe("lsblk", "List of all block mount points (partition table)")
}

/* ========== */
/* XfceSource */
/* ========== */

type XfceSource struct {
}

func (src XfceSource) Backup(writer *zip.Writer) result.Result {
	if result := ProcessDir("~/.config/xfce4/panel", "xfce4/panel", writer); !result.IsSuccess() {
		return result
	}
	if result := ProcessFile("~/.config/xfce4/xfconf/xfce-perchannel-xml/xfce4-panel.xml", "xfce4/xfce4-panel.xml", writer); !result.IsSuccess() {
		return result
	}
	return result.Success("Xfce configuration added")
}

func (src XfceSource) Restore(zip zipfile.ZipFile) result.Result {

	files := zip.FilesStartingWith("xfce4/panel")
	if result := copyFilesToFolder(files, "~/.config/xfce4/panel"); !result.IsSuccess() {
		return result
	}

	if zip.HasFile("xfce4/xfce4-panel.xml") {
		if res := zip.GetFile("xfce4/xfce4-panel.xml").Write("~/.config/xfce4/xfconf/xfce-perchannel-xml/xfce4-panel.xml"); !res.IsSuccess() {
			return result.Failure("xfce4/xfce4-panel.xml update failed")
		}
	}

	return result.Success("Xfce config restored")
}

func (src XfceSource) Describe() {
	result.Describe("Xfce panel", "configuration files will be stored into xfce4")
}

/* ================= */
/* COMMON / INTERNAL */
/* ================= */

func copyFilesToFolder(files []zipfile.ZipElement, destinationPath string) result.Result {
	if len(files) == 0 {
		return result.New(true, "No files to restore to "+destinationPath)
	}

	var inError []string
	for _, file := range files {
		if file.IsValid() {
			dst := strings.TrimSuffix(destinationPath, "/") + file.Name()[strings.Index(file.Name(), "/"):]
			if result := file.Write(dst); !result.IsSuccess() {
				inError = append(inError, file.Name())
			}
		} else {
			inError = append(inError, file.Name())
		}
	}

	if len(inError) == 0 {
		return result.Success(strconv.Itoa(len(files)) + " files written to " + destinationPath)
	}

	return result.Failure(strconv.Itoa(len(inError)) + " files in error : " + strings.Join(inError, ","))
}
