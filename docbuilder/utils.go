package docbuilder

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadFile(fpath string) (content []byte, err error) {
	content, err = ioutil.ReadFile(fpath)
	if err != nil {
		return
	}
	return
}

func CheckFile(fpath string) (exist bool, err error) {
	if _, err = os.Stat(fpath); os.IsNotExist(err) {
		return
	}

	exist = true
	return
}

func WriteFile(fpath string, content []byte, force bool) (err error) {
	if ok, _ := CheckFile(fpath); ok && !force {
		err = fmt.Errorf("file exist")
		return
	}

	os.MkdirAll(filepath.Dir(fpath), 0755)
	err = ioutil.WriteFile(fpath, content, 0644)
	if err != nil {
		return
	}
	return
}
