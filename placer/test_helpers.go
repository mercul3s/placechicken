package placer

import (
	"github.com/stretchr/testify/mock"
)

// MockDir is a mock interface for file directory functions.
type MockDir struct {
	mock.Mock
}

// List is a mock directory listing method.
func (t *MockDir) list(p string) ([]Image, error) {
	args := t.Called(p)
	return args.Get(0).([]Image), args.Error(1)
}

// RandImg is a mock random get method
func (t *MockDir) RandImg(p string) (Image, error) {
	args := t.Called(p)
	return args.Get(0).(Image), args.Error(1)
}
