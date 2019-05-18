package assets

import (
	"log"
	"os"
	"strings"

	"github.com/gandrille/go-commons/result"

	"github.com/gandrille/go-commons/filesystem"
)

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go resources/ resources/profils resources/zim-simple-web-template resources/zim-simple-web-template/SimpleWeb

// AssetAsString gets an asset and convert it into a string
// If an error occurs (ie not found), exit(1)
// (but of course you will always ask for an existing asset, wont you?)
// TODO remove the Exit(1) ?
func AssetAsString(assetName string) string {
	bytes, err := Asset(assetName)
	if err != nil {
		log.Fatal("FATAL ERROR: Can't find '" + assetName + "'")
		os.Exit(1)
	}
	return string(bytes)
}

// CreateOrAppendAssetContentIfNotInFile appends the content of an asset to a file on the filesystem, if the file does NOT already contain the asset content.
// The destination file is created if it does not already exists.
func CreateOrAppendAssetContentIfNotInFile(assetName, filePath string) result.Result {
	fileContent := AssetAsString(assetName)
	switch updated, err := filesystem.CreateOrAppendIfNotInFile(filePath, fileContent); {
	case err != nil:
		return result.Failure(filePath + " NOT updated : " + err.Error())
	case updated:
		return result.Success(filePath + " updated with specific content")
	case !updated:
		return result.Success(filePath + " already includes expected content")
	}
	return result.Failure("Just for Linter not to complain")
}

// WriteAsset writes an asset to the filesystem.
func WriteAsset(assetName, filePath string, overwrite bool) result.Result {

	// Get new content
	bytes, err := Asset(assetName)
	if err != nil {
		return result.Failure("Can't find asset inside executable")
	}

	return filesystem.WriteBinaryFile(filePath, bytes, overwrite)
}

// CopyAssetDirectory writes an asset directory to the filesystem.
func CopyAssetDirectory(assetPrefix, dirPath string, overwrite bool) result.Result {

	// Ensure assetPrefix ends with '/'
	if !strings.HasSuffix(assetPrefix, "/") {
		assetPrefix = assetPrefix + "/"
	}

	// Ensure dirPath ends with '/'
	if !strings.HasSuffix(dirPath, "/") {
		dirPath = dirPath + "/"
	}

	msg := ""
	for _, assetName := range AssetNames() {
		if strings.HasPrefix(assetName, assetPrefix) {
			filePath := dirPath + strings.TrimPrefix(assetName, assetPrefix)
			if res := WriteAsset(assetName, filePath, overwrite); !res.IsSuccess() {
				return result.Failure("Can't copy " + assetName + " to " + filePath + ": " + res.Message())
			} else {
				if msg != "" {
					msg += "\n"
				}
				msg += res.Message()
			}
		}
	}

	return result.Success(msg)
}
