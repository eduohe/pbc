package pbc

import (
	"os"
	"path/filepath"
)

/*
  Finds target files for packaging into passbook pass
*/
func findTargets(root string) ([]string, error) {
	var targets []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if isTarget(path) && !info.IsDir() {
			targets = append(targets, path)
		}

		if info.IsDir() && !isTarget(path) && path != "." {
			return filepath.SkipDir
		}
		return nil
	})

	return targets, err
}

func isTarget(file string) bool {
	name := ""
	parts := filepath.SplitList(file)
	if len(parts) > 0 {
		name = parts[len(parts)-1]
	}

	if name == "signature" {
		return false
	}

	if match, err := filepath.Match("*.pkpass", name); err == nil && match {
		return false
	}

	return true
}
