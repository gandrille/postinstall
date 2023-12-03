package systemmode

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed profils
var efs embed.FS

// FunctionInfos type
type Profile struct {
	Path     string
	File     string
	Name     string
	Contents string
}

func GetProfils() (profils []Profile, err error) {

	if err := fs.WalkDir(&efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if content, err := efs.ReadFile(path); err != nil {
			return err
		} else {
			tokens := strings.Split(path, "/")
			file := tokens[len(tokens)-1]
			name := strings.TrimSuffix(file, ".txt")
			name = strings.ReplaceAll(name, "-", " ")

			profile := Profile{
				Path:     path,
				File:     file,
				Name:     name,
				Contents: string(content),
			}
			profils = append(profils, profile)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return profils, nil
}
