package placer

import (
	"fmt"
	"image"
	"strings"

	"github.com/disintegration/imaging"
)

// Place holds configuration info for the image resizing function.
type Place struct {
	Dir              Directory
	OriginalFilePath string
	ResizedFilePath  string
}

// Image ...
type Image struct {
	Name string
}

// Directory provides a function for listing files in a local or
// remote directory.
type Directory interface {
	List(string) ([]Image, error)
	RandImg(string) (Image, error)
}

// Config returns a Place configuration with file settings
func Config(dir Directory, oPath string, rPath string) Place {
	return Place{
		Dir:              dir,
		OriginalFilePath: oPath,
		ResizedFilePath:  rPath,
	}
}

// GetImage takes a width and height and returns an image sized to the
// dimensions specified.
func (p *Place) GetImage(w int, h int) (image.Image, error) {
	// get a random image from the images dir
	srcImg, err := p.Dir.RandImg(p.OriginalFilePath)
	if err != nil {
		return nil, err
	}

	src, err := imaging.Open(p.OriginalFilePath + srcImg.Name)
	if err != nil {
		return nil, err
	}

	name := p.newFileName(srcImg.Name, w, h)
	resized := imaging.Resize(src, w, h, imaging.Lanczos)
	err = imaging.Save(resized, name)
	if err != nil {
		fmt.Printf("save error: %s\n", err.Error())
		return nil, err
	}
	return resized, nil
}

func (p *Place) newFileName(name string, w int, h int) string {
	idx := strings.Index(name, ".jpg")
	if idx > -1 {
		return fmt.Sprintf("%s%s-%dx%d%s", p.ResizedFilePath, name[:idx], w, h, name[idx:])
	}
	return name
}
