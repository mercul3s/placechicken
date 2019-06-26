package placer

import (
	"io/ioutil"
	"os"
)

type dir struct{}

func (r *dir) List(p string) ([]os.FileInfo, error) {
	fileList, err := ioutil.ReadDir(p)
	return fileList, err
}
