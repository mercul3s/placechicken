package testHelpers

import (
	"os"

	"github.com/stretchr/testify/mock"
)

// Dir is a mock interface for directory functions.
type Dir struct {
	mock.Mock
}

// List is a mock directory listing method.
func (t *Dir) List(p string) ([]os.FileInfo, error) {
	args := t.Called(p)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}
