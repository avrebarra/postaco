package fsx

import (
	"path/filepath"
	"strings"
)

func HasExtension(fpath, ext string) bool {
	return strings.HasSuffix(fpath, ext)
}

func RemoveExtension(fpath, ext string) string {
	return strings.TrimSuffix(fpath, ext)
}

func PathAsRelative(fpath, base string) (rel string, err error) {
	absPath, err := filepath.Abs(fpath)
	if err != nil {
		return
	}

	absBase, err := filepath.Abs(base)
	if err != nil {
		return
	}

	rel = filepath.Clean(strings.TrimPrefix(absPath, absBase))

	return
}
