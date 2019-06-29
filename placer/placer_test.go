package placer

import (
	"fmt"
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageResizer(t *testing.T) {
	if _, err := os.Stat("/tmp/placechicken"); os.IsNotExist(err) {
		err := os.Mkdir("/tmp/placechicken", 0700)
		if err != nil {
			t.Fatal(err)
		}
	}
	tt := []struct {
		name           string
		fileName       string
		path           string
		dirList        []os.FileInfo
		width          int
		height         int
		expectedResult image.Rectangle
		expectedErr    error
	}{
		{
			name:     "image resized to 300x500",
			fileName: "original-test-image.jpg",
			path:     "../static/images/test/",
			width:    500,
			height:   300,
			expectedResult: image.Rectangle{
				Max: image.Point{X: 500, Y: 300},
			},
			expectedErr: nil,
		},
	}
	for _, table := range tt {
		td := MockDir{}

		place := Place{
			Dir:              &td,
			OriginalFilePath: "../static/images/test/",
			ResizedFilePath:  "/tmp/placechicken/",
		}
		// get test image
		fileInfo, err := os.Stat(table.path + table.fileName)
		if err != nil {
			assert.FailNowf(t, "could not get test image", err.Error())
		}
		fileList := []Image{}
		file := Image{Name: fileInfo.Name()}
		fileList = append(fileList, file)
		td.On("List", "../static/images/test/").Return(fileList, table.expectedErr)
		td.On("RandImg", "../static/images/test/").Return(file, table.expectedErr)
		image, err := place.GetImage(table.width, table.height)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(image.Bounds().Max)
	}
	err := os.RemoveAll("/tmp/placechicken")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewFileName(t *testing.T) {
	tt := []struct {
		name     string
		path     string
		width    int
		height   int
		expected string
	}{
		{
			name:     "test-image.jpg",
			path:     "/testpath/",
			width:    400,
			height:   800,
			expected: "/testpath/test-image-400x800.jpg",
		},
		{
			name:     "test-image",
			path:     "",
			width:    400,
			height:   800,
			expected: "test-image",
		},
	}

	for _, table := range tt {
		p := Place{
			OriginalFilePath: table.path,
			ResizedFilePath:  table.path,
		}
		name := p.newFileName(table.name, table.width, table.height)
		assert.Equal(t, table.expected, name)
	}
}
