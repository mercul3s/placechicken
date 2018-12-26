package placer

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	"github.com/disintegration/imaging"
)

// Place holds configuration info for the image resizing function.
type Place struct {
	Dir              Directory
	OriginalFilePath string
	ResizedFilePath  string
}

// Directory provides a function for listing files in a directory.
type Directory interface {
	List(string) ([]os.FileInfo, error)
}

type dir struct{}

// Config returns a configured Place
func Config() Place {
	return Place{
		Dir:              &dir{},
		OriginalFilePath: "../static/images/",
		ResizedFilePath:  "../static/images/resized/",
	}
}

func (r *dir) List(p string) ([]os.FileInfo, error) {
	fileList, err := ioutil.ReadDir(p)
	return fileList, err
}

// GetImage takes a width and height and returns an image sized to the
// dimensions specified.
func (p *Place) GetImage(w int, h int) (string, error) {
	// get a random image from the images dir
	srcImg, err := p.randImg()
	if err != nil {
		return "error", err
	}

	src, err := imaging.Open(p.OriginalFilePath + srcImg.Name())
	if err != nil {
		return "error", err
	}

	name := p.newFileName(srcImg.Name(), w, h)
	resized := imaging.Resize(src, w, h, imaging.Lanczos)
	err = imaging.Save(resized, name)
	return name, nil
}

func (p *Place) randImg() (os.FileInfo, error) {
	files, err := p.Dir.List(p.OriginalFilePath)
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		randIdx := rand.Intn(len(files))
		randFile := files[randIdx]
		return randFile, nil
	}
	return nil, nil
}

func (p *Place) newFileName(name string, w int, h int) string {
	idx := strings.Index(name, ".jpg")
	if idx > -1 {
		return fmt.Sprintf("%s%s-%dX%d%s", p.ResizedFilePath, name[:idx], w, h, name[idx:])
	}
	return name
}
