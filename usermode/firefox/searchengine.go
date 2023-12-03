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

const filename = "search.json.mozlz4"
const appName = "Firefox"

// Disclaimer from https://hg.mozilla.org/mozilla-central/file/tip/toolkit/components/search/SearchEngine.jsm#l228
const disclaimer = "By modifying this file, I agree that I am doing so only within $appName itself, using official, user-driven search engine selection processes, and in a way which does not circumvent user consent. I acknowledge that any attempt to change this file from outside of $appName is a malicious act, and will be responded to accordingly."

// SetSearchEngine sets the search engine inside the dedicated file
func SetSearchEngine(engineName string) result.Result {

	// Get file path
	filePath, err0 := getSearchEngineFilePath()
	if err0 != nil {
		return result.NewError("Error while computing file path: " + err0.Error())
	}

	// Get file content
	oldContent, err1 := getMozlz4FileContent(filePath)
	if err1 != nil {
		return result.NewError("Error while reading " + filePath + ": " + err1.Error() + " (are you sure Firefox has already been started once?)")
	}

	// Find metaData section
	start, end, oldMetaData, err2 := getMetaData(oldContent)
	if err2 != nil {
		return result.NewError("Error while reading " + filePath + ": " + err2.Error())
	}

	// Get new metaData
	newMetaData, err3 := getNewMetadataContent(oldMetaData, engineName)
	if err3 != nil {
		return result.NewError("Error while computing new Metadata content: " + err3.Error())
	}

	// Compute new content
	newContent := oldContent[:start] + newMetaData + oldContent[end:]

	// Check for updates
	if oldMetaData == newMetaData {
		return result.NewUnchanged("Search engine was already " + engineName)
	}

	// Write content
	if err4 := writeContent(newContent, filePath); err4 != nil {
		return result.NewError(err4.Error())
	}

	// Success!
	return result.NewUpdated("Search engine is now " + engineName)
}

func getSearchEngineFilePath() (string, error) {
	path, err1 := GetProfileFolder()
	if err1 != nil {
		return "", err1
	}
	return path + "/" + filename, nil
}

func getMozlz4FileContent(filePath string) (string, error) {

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

func getMetaData(content string) (int, int, string, error) {
	idx := strings.Index(content, "\"metaData\"")
	if idx == -1 {
		return -1, -1, "", errors.New("Parsing error: No metaData object found")
	}
	start := strings.Index(content[idx:], "{")
	if start == -1 {
		return -1, -1, "", errors.New("Parsing error: No '{' found after metaData")
	}
	start += idx + 1
	end := strings.Index(content[start:], "}")
	if end == -1 {
		return start, -1, "", errors.New("Parsing error: No '}' found after metaData")
	}
	end += start
	metaData := content[start:end]

	return start, end, metaData, nil
}

func getNewMetadataContent(content, engineName string) (string, error) {
	// Hash computation
	hash, err1 := computeHash(appName, engineName)
	if err1 != nil {
		return "", errors.New("Error while computing hash: " + err1.Error())
	}

	// Empty Section : here is the full content
	if len(content) == 0 {
		return "\"current\":\"" + engineName + "\",\"hash\":\"" + hash + "\"", nil
	}

	// Section is not empty

	// GetInfos (this implementation is a bit weak...)
	curStart, curEnd, curErr := findValue(content, "current")
	if curErr != nil {
		return "", errors.New("Error while parsing json: " + curErr.Error())
	}
	curVal := content[curStart:curEnd]
	hashStart, hashEnd, hashErr := findValue(content, "hash")
	if hashErr != nil {
		return "", errors.New("Error while parsing json: " + hashErr.Error())
	}
	hashVal := content[hashStart:hashEnd]

	// Here is the logic!
	if engineName == curVal {
		if hashVal == hash {
			return content, nil
		}
		return content, errors.New("Search engine is already " + engineName + " BUT HASH IS MALFORMED (LEFT UNCHANGED)")
	}

	// Update is needed!
	if curStart > hashStart {
		return "", errors.New("Json file content is not as expected... I prefer not to update")
	}

	newContent := content[:curStart] + engineName + content[curEnd:hashStart] + hash + content[hashEnd:]
	return newContent, nil
}

func computeHash(appName, engineName string) (string, error) {
	key, err1 := GetProfileKey()
	if err1 != nil {
		return "", err1
	}
	salt := key + engineName + strings.ReplaceAll(disclaimer, "$appName", appName)
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

func writeContent(content, filePath string) error {
	r := strings.NewReader(content)
	w := &bytes.Buffer{}
	err := Compress(r, w, len(content))
	if err != nil {
		return errors.New("Failed to compress new data: " + err.Error())
	}

	if res := filesystem.WriteBinaryFile(filePath, w.Bytes(), true); res.IsFailure() {
		return errors.New("Failed to update search engine: " + res.Message())
	}

	return nil
}
