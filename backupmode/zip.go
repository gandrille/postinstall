package backupmode

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

type srcdst struct {
	src string
	dst string
}

type FileMatcher interface {
	Match(element srcdst) bool
}

type MatchAll struct {
}

func (f MatchAll) Match(element srcdst) bool {
	return true
}

func ProcessDir(dirPath, dstName string, writer *zip.Writer) result.Result {
	return ProcessDirFiltered(dirPath, dstName, MatchAll{}, writer)
}

func ProcessDirFiltered(dirPath, dstName string, matcher FileMatcher, writer *zip.Writer) result.Result {
	dirName := strings.Replace(dirPath, filesystem.HomeDir(), "~", 1)

	if exists, err := filesystem.FolderExists(dirPath); err == nil && !exists {
		return result.NewUnchanged(dirName + " does NOT exists")
	}

	files, err := parseDir(dirPath, dstName)
	if err != nil {
		return result.NewError(err.Error())
	}

	var filtered []srcdst
	for _, file := range files {
		if matcher.Match(file) {
			filtered = append(filtered, file)
		}

	}

	return addFiles(dirName, filtered, writer)
}

func parseDir(dirPath, dstName string) ([]srcdst, error) {
	dirPath = strings.Replace(dirPath, "~", filesystem.HomeDir(), 1)
	dirPath = strings.TrimSuffix(dirPath, "/")
	dstName = strings.TrimPrefix(dstName, "/")

	var result []srcdst

	// Check folder existence
	if exists, err := filesystem.FolderExists(dirPath); err != nil {
		return result, errors.New(dirPath + " can't be parsed")
	} else if !exists {
		return result, errors.New(dirPath + " does NOT exists (skipped)")
	}

	// get files list
	paths, err := filesystem.FolderFiles(dirPath)
	if err != nil {
		return result, errors.New("Can't get list of files into " + dirPath)
	}

	// Registers files
	for _, src := range paths {
		dst := dstName + strings.TrimPrefix(src, dirPath)
		result = append(result, srcdst{src, dst})
	}

	return result, nil
}

func addFiles(dirName string, files []srcdst, writer *zip.Writer) result.Result {
	var results []result.Result
	for _, file := range files {
		results = append(results, ProcessFile(file.src, file.dst, writer))
	}

	success := 0
	failures := 0
	for _, result := range results {
		if result.IsSuccess() {
			success++
		} else {
			failures++
		}
	}

	if failures == 0 {
		if success == 0 {
			return result.NewUnchanged("NO files available into " + dirName)
		}
		return result.NewUpdated(strconv.Itoa(success) + " files added from " + dirName)
	} else {
		return result.NewError(strconv.Itoa(failures) + " failures while adding files from " + dirName)
	}
}

func ProcessFile(filePath, dstName string, writer *zip.Writer) result.Result {
	filePath = strings.Replace(filePath, "~", filesystem.HomeDir(), 1)
	fileName := strings.Replace(filePath, filesystem.HomeDir(), "~", 1)

	// Assert file exists
	exists, err1 := filesystem.RegularFileExists(filePath)
	if err1 != nil {
		return result.NewError(err1.Error())
	}
	if !exists {
		return result.NewUnchanged(fileName + " does NOT exists (not added)")
	}

	// Read file content
	content, err2 := ioutil.ReadFile(filePath)
	if err2 != nil {
		return result.NewError(fileName + " is not readable")
	}

	return ProcessBytes(fileName, dstName, content, writer)
}

func ProcessBytes(displayName, dstName string, content []byte, writer *zip.Writer) result.Result {
	// Create entry
	f, err3 := writer.Create(dstName)
	if err3 != nil {
		return result.NewError(err3.Error())
	}

	// Write file content
	_, err4 := f.Write(content)
	if err4 != nil {
		return result.NewError(displayName + " could NOT be written into the archive")
	}

	return result.NewUpdated(displayName + " added")
}
