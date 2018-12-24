package resizer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageResizer(t *testing.T) {
	tt := []struct {
		name                  string
		width                 int
		height                int
		expectedImageReturned string
	}{
		{
			name:                  "image resized to 300x500",
			width:                 500,
			height:                300,
			expectedImageReturned: "image300x500.jpg",
		},
	}
	for _, table := range tt {
		fmt.Println(table.expectedImageReturned)
	}
}

func TestGetRandomImageValidPath(t *testing.T) {
	rImage, err := randImg("../static/images/test/")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, rImage)
}

func TestGetRandomImageWithInvalidPath(t *testing.T) {
	rImage, err := randImg("./bogus")
	assert.Nil(t, rImage)
	assert.NotNil(t, err)
}
