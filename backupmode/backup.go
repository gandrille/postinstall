package backupmode

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/go-commons/strpair"
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
	if err != nil {
		return results.Add(result.NewError("Error while creating zip file " + zipFilePath))
	}
	defer f.Close()

	if _, err := f.Write(buf.Bytes()); err != nil {
		return results.Add(result.NewError("Error while writing zip file"))
	}

	return results.Add(result.NewUpdated("zip file written to " + zipFilePath))
}

func backupAll(writer *zip.Writer) result.Set {
	var results result.Set
	sources := getSources()

	kv := make([]strpair.StrPair, 0)
	for _, src := range sources {
		r := src.Backup(writer, &kv)
		results.Add(r)
	}
	results.Add(writeStore(writer, kv))

	return results
}

func writeStore(writer *zip.Writer, store []strpair.StrPair) result.Result {
	content := "# Key Value Store content\n"
	for _, pair := range store {
		content += pair.Str1() + "=" + pair.Str2() + "\n"
	}
	return ProcessBytes("Key/Value store", "kvstore.txt", []byte(content), writer)
}

func restoreAll(zip zipfile.ZipFile) result.Set {
	var results result.Set
	sources := getSources()

	kv := make([]strpair.StrPair, 0)
	results.Add(readStore(zip, &kv))

	for _, src := range sources {
		results.Add(src.Restore(zip, kv))
	}

	return results
}

func readStore(zip zipfile.ZipFile, store *[]strpair.StrPair) result.Result {
	file := zip.GetFile("kvstore.txt")
	if file == nil {
		return result.NewInfo("No Key/Value store available")
	}

	content, err := file.StringContent()
	if err != nil {
		return result.NewError("Can't read Key/Value store: " + err.Error())
	}

	for _, line := range strings.Split(content, "\n") {
		if len(line) != 0 && !strings.HasPrefix(line, "#") {
			idx := strings.Index(line, "=")
			if idx == -1 {
				return result.NewError("'=' separator not found in line " + line)
			}
			key := line[:idx]
			val := line[idx+1:]
			*store = append(*store, strpair.New(key, val))
		}
	}

	return result.NewInfo("Key/Value store loaded with " + strconv.Itoa(len(*store)) + " pair(s)")
}
