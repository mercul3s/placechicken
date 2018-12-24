package resizer

import (
	"fmt"
	"image"
	"io/ioutil"
	"math/rand"

	"github.com/disintegration/imaging"
)

// GetImage takes a width and height and returns an image sized to the
// dimensions specified.
func GetImage(h int, w int) (string, error) {
	// get a random image from the images dir
	srcImg, err := randImg("./static/images/")
	if err != nil {
		return "error", err
	}
	resized := imaging.Resize(srcImg, w, h, imaging.Lanczos)
	fmt.Println(resized)
	return "hello", err
	//	err := imaging.Save(resized, )
}

func randImg(path string) (image.Image, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	randIdx := rand.Intn(len(files))
	randFile := files[randIdx]
	src, err := imaging.Open(path + randFile.Name())
	return src, err
}
