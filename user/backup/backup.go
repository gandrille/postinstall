package backup

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/go-commons/zipfile"
)

func Backup(args []string) {
	zipFilePath := getFilePath(args)
	results := backups(zipFilePath)
	results.Print()
}

func Restore(args []string) {
	zipContent := zipfile.Open(args[0])
	if !zipContent.IsValid() {
		result.PrintError(zipContent.Err.Error())
	} else {
		results := restoreAll(zipContent)
		results.Print()
	}
}

func getFilePath(args []string) string {
	if len(args) == 0 {
		now := time.Now()
		zipFileName := fmt.Sprintf("%d-%02d-%02d-%s-backup.zip", now.Year(), now.Month(), now.Day(), env.Hostname())
		return filesystem.HomeDir() + "/" + zipFileName
	} else {
		return args[0]
	}
}

func Describe() {
	fmt.Println("Here is how the information will be stored into the archive:")
	for _, element := range getSources() {
		element.Describe()
	}
}

func backups(zipFilePath string) result.Set {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	results := backupAll(w)

	if err := w.Close(); err != nil {
		return results.Add(result.NewError("Error while closing zip byte content"))
	}

	f, err := os.Create(zipFilePath)
	defer f.Close()
	if err != nil {
		return results.Add(result.NewError("Error while creating zip file " + zipFilePath))
	}

	if _, err := f.Write(buf.Bytes()); err != nil {
		return results.Add(result.NewError("Error while writing zip file"))
	}

	return results.Add(result.NewUpdated("zip file written to " + zipFilePath))
}

func backupAll(writer *zip.Writer) result.Set {
	var results result.Set
	sources := getSources()

	for _, src := range sources {
		r := src.Backup(writer)
		results.Add(r)
	}

	return results
}

func restoreAll(zip zipfile.ZipFile) result.Set {
	var results result.Set
	sources := getSources()

	for _, src := range sources {
		results.Add(src.Restore(zip))
	}

	return results
}
