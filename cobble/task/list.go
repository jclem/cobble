package task

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/jclem/cobble/cobble/config"
)

var list []string

func List() ([]string, error) {
	if list != nil {
		return list, nil
	}

	scaffoldsDir, err := config.ScaffoldsDir()
	if err != nil {
		list = []string{}
		return nil, err
	}

	var scaffolds []string

	// Walk through scaffoldsDir, returning all scaffolds. A scaffold is a
	// directory with a "scaffold.yml" file inside of it, such as
	// "pretter/scaffold.yml".

	if err := filepath.WalkDir(scaffoldsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		filename := filepath.Base(path)
		if filename == "scaffold.yml" || filename == "scaffold.yaml" {
			name := strings.TrimPrefix(path, fmt.Sprintf("%s/", scaffoldsDir))
			name = strings.TrimSuffix(name, fmt.Sprintf("/%s", filename))
			scaffolds = append(scaffolds, name)
		}

		return nil
	}); err != nil {
		list = []string{}
		return nil, err
	}

	list = scaffolds
	return scaffolds, nil
}
