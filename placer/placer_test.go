package placer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testDir struct {
	mock.Mock
}

func (t *testDir) List(p string) ([]os.FileInfo, error) {
	args := t.Called(p)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

func TestImageResizer(t *testing.T) {
	err := os.Mkdir("/tmp/placechicken", 0700)
	if err != nil {
		t.Fatal(err)
	}
	tt := []struct {
		name           string
		fileName       string
		path           string
		dirList        []os.FileInfo
		width          int
		height         int
		expectedResult string
		expectedErr    error
	}{
		{
			name:           "image resized to 300x500",
			fileName:       "original-test-image.jpg",
			path:           "../static/images/test/",
			width:          500,
			height:         300,
			expectedResult: "/tmp/placechicken/original-test-image-500X300.jpg",
			expectedErr:    nil,
		},
	}
	for _, table := range tt {
		td := &testDir{}

		place := Place{
			Dir:              td,
			OriginalFilePath: "../static/images/test/",
			ResizedFilePath:  "/tmp/placechicken/",
		}
		// get test image
		fileInfo, err := os.Stat(table.path + table.fileName)
		fileList := []os.FileInfo{fileInfo}
		if err != nil {
			t.Fatal(err)
		}
		td.On("List", "../static/images/test/").Return(fileList, table.expectedErr)
		resized, err := place.GetImage(table.width, table.height)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, table.expectedResult, resized)
		// check that the new file exists
		assert.FileExists(t, resized)
	}
	err = os.RemoveAll("/tmp/placechicken")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetRandomImage(t *testing.T) {
	tt := []struct {
		name           string
		path           string
		expectedResult string
		expectedError  error
	}{
		{
			name:           "returns only image in directory",
			path:           "../static/images/test/",
			expectedResult: "original-test-image.jpg",
			expectedError:  nil,
		},
		{
			name:           "returns error for non existent directory",
			path:           "./bogus",
			expectedResult: "",
			expectedError:  nil,
		},
	}

	for _, table := range tt {
		d := &dir{}
		place := Place{
			Dir:              d,
			OriginalFilePath: table.path,
		}
		rImage, err := place.randImg()
		if rImage != nil {
			assert.Equal(t, table.expectedResult, rImage.Name())
			assert.Nil(t, err)
		} else {
			assert.Error(t, err)
		}
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
			expected: "/testpath/test-image-400X800.jpg",
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
