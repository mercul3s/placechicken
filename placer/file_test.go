package placer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirList(t *testing.T) {
	path := "../static/images/test/"
	d := Dir{}
	iList, err := d.List(path)
	assert.Nil(t, err)
	assert.Equal(t, len(iList), 1)
}

func TestDirGetRandomImage(t *testing.T) {
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
			expectedError:  errors.New("open ./bogus: no such file or directory"),
		},
	}

	for _, table := range tt {
		d := Dir{}
		rImage, err := d.RandImg(table.path)
		assert.Equal(t, table.expectedResult, rImage.Name)
		if err != nil {
			assert.Equal(t, table.expectedError.Error(), err.Error())
		}
	}
}
