package firefox

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

// IMPORTANT
// This file provides functions that works... most of the time...
// Use at your own risks!
// TODO update it according to Firefox specs

const filename = "search.json.mozlz4"
const searchEngine = "Qwant"
const appName = "Firefox"

// Disclaimer from https://hg.mozilla.org/mozilla-central/file/tip/toolkit/components/search/SearchEngine.jsm#l228
const disclaimer = "By modifying this file, I agree that I am doing so only within $appName itself, using official, user-driven search engine selection processes, and in a way which does not circumvent user consent. I acknowledge that any attempt to change this file from outside of $appName is a malicious act, and will be responded to accordingly."

// SetSearchEngine sets the search engine inside the dedicated file
func SetSearchEngine() result.Result {

	// Get file path
	jsonFilePath, err0 := getJSONFilePath()
	if err0 != nil {
		return result.NewError("Error while computing file path: " + err0.Error())
	}

	// Get file content
	content, err1 := getJSONFileContent(jsonFilePath)
	if err1 != nil {
		return result.NewError("Error while reading " + filename + ": " + err1.Error())
	}

	// GetInfos (this implementation is a bit weak...)
	curStart, curEnd, curErr := findValue(content, "current")
	if curErr != nil {
		return result.NewError("Error while parsing json: " + curErr.Error())
	}
	curVal := content[curStart:curEnd]
	hashStart, hashEnd, hashErr := findValue(content, "hash")
	if hashErr != nil {
		return result.NewError("Error while parsing json: " + hashErr.Error())
	}
	hashVal := content[hashStart:hashEnd]

	// Hash computation
	hash, err2 := computeHash()
	if err2 != nil {
		return result.NewError("Error while computing hash: " + err2.Error())
	}

	// Here is the logic!
	if searchEngine == curVal {
		if hashVal == hash {
			return result.NewUnchanged("Search engine is already " + searchEngine)
		}
		return result.NewUnchanged("Search engine is already " + searchEngine + " BUT HASH IS MALFORMED (LEFT UNCHANGED)")
	}

	// Update is needed!
	if curStart > hashStart {
		return result.NewError("Json file content is not as expected... I prefer not to update")
	}
	newContent := content[:curStart] + searchEngine + content[curEnd:hashStart] + hash + content[hashEnd:]

	// Compress data
	r := strings.NewReader(newContent)
	w := &bytes.Buffer{}
	err := Compress(r, w, len(newContent))
	if err != nil {
		return result.NewError("Failed to compress new data: " + err.Error())
	}

	if res := filesystem.WriteBinaryFile(jsonFilePath, w.Bytes(), true); res.IsFailure() {
		return result.NewError("Failed to update search engine: " + res.Message())
	}
	return result.NewUpdated("Search engine is now " + searchEngine)
}

func getJSONFilePath() (string, error) {
	// Profile folder
	path, err1 := GetProfileFolder()
	if err1 != nil {
		return "", err1
	}
	return path + "/" + filename, nil
}

func getJSONFileContent(filePath string) (string, error) {

	// Raw file content
	raw, err2 := filesystem.ReadFileAsBinary(filePath)
	if err2 != nil {
		return "", err2
	}

	// Json content
	rt, err3 := NewDecompressReader(bytes.NewReader(raw))
	if err3 != nil {
		return "", errors.New("Failed to decompress data: " + err3.Error())
	}
	out, err4 := ioutil.ReadAll(rt)
	if err4 != nil {
		return "", errors.New("Failed to ReadAll data: " + err4.Error())
	}

	return string(out), nil
}

func computeHash() (string, error) {
	key, err1 := GetProfileKey()
	if err1 != nil {
		return "", err1
	}
	salt := key + searchEngine + strings.ReplaceAll(disclaimer, "$appName", appName)
	sum256 := sha256.Sum256([]byte(salt))
	encoded64 := base64.StdEncoding.EncodeToString(sum256[:])
	return encoded64, nil
}

func findValue(content, hash string) (int, int, error) {
	full := "\"" + hash + "\""

	// find key
	idx := strings.Index(content, full)
	if idx == -1 {
		return 0, 0, errors.New("key " + full + " not found")
	}
	offsetStart := idx + len(full)
	cut := content[offsetStart:]

	// Value start
	idx = strings.Index(cut, "\"")
	if idx == -1 {
		return 0, 0, errors.New("Value start not found for key " + full)
	}
	offsetStart += idx + 1
	cut = content[offsetStart:]

	// Value end
	idx = strings.Index(cut, "\"")
	if idx == -1 {
		return offsetStart, 0, errors.New("Value end not found for key " + full)
	}
	offsetEnd := offsetStart + idx

	return offsetStart, offsetEnd, nil
}
