package processors

import (
	"github.com/gandrille/go-commons/filesystem"
	"github.com/gandrille/go-commons/result"
)

type imagemagick struct {
}

func (e imagemagick) Key() string {
	return "configure-imagemagick"
}

func (e imagemagick) Describe(args []string) string {
	return "imagemagick config: allows full pdf file manipulation"
}

func (e imagemagick) Run(args []string) result.Result {
	imgMagickFile := "/etc/ImageMagick-6/policy.xml"

	res1 := filesystem.UpdateLineInFile(imgMagickFile, "  <policy domain=\"coder\" rights=\"none\" pattern=\"PDF\" />", "  <policy domain=\"coder\" rights=\"read|write\" pattern=\"PDF\" />", false)
	res2 := filesystem.UpdateLineInFile(imgMagickFile, "  <policy domain=\"resource\" name=\"disk\" value=\"1GiB\" />", "  <policy domain=\"resource\" name=\"disk\" value=\"8GiB\" />", false)

	if res1.IsError() {
		return result.NewError(imgMagickFile + " update failed: " + res1.Message())
	}
	if res2.IsError() {
		return result.NewError(imgMagickFile + " update failed: " + res2.Message())
	}
	if res1.IsUpdated() || res2.IsUpdated() {
		return result.NewUpdated(imgMagickFile + " successfully updated")
	}
	if res1.IsUnchanged() && res2.IsUnchanged() {
		return result.NewUnchanged(imgMagickFile + " left unchanged")
	}

	return result.NewError("Unexpected status: res1=" + res1.Status().String() + " and res2=" + res2.Status().String())
}
