package placer

import (
	"io/ioutil"
	"math/rand"
	"strings"
)

// Dir struct exists as a placeholder to allow abstracting os directory methods.
type Dir struct{}

// List returns only image contents of a directory.
func (d *Dir) List(p string) ([]Image, error) {
	fileList, err := ioutil.ReadDir(p)
	i := []Image{}
	for _, file := range fileList {
		if strings.Contains(file.Name(), "original") {
			i = append(i, Image{Name: file.Name()})
		}
	}
	return i, err
}

// RandImg gets the contents of a directory, filters it for only images, and
// then returns a random image.
func (d *Dir) RandImg(p string) (Image, error) {
	files, err := d.List(p)
	i := Image{}
	if err != nil {
		return i, err
	}

	if len(files) > 0 {
		randIdx := rand.Intn(len(files))
		randFile := files[randIdx]
		return randFile, nil
	}
	return i, nil
}
