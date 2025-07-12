package projectpath

import (
	"errors"
	"os"
	"path/filepath"
)

// FindGoModDir walks up the directory tree to find the directory containing go.mod.
func FindGoModDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached root
		}
		dir = parent
	}

	return "", errors.New("go.mod not found")
}

// MustRoot returns go.mod directory or current working directory as fallback.
func MustRoot() string {
	root, err := FindGoModDir()
	if err != nil {
		root, _ = os.Getwd()
	}
	return root
}
