package fsx

import (
	"fmt"
	"io"
	"os"
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

func CopyFile(src, dst string, BUFFERSIZE int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("File %s already exists", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}
