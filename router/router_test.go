package router

import (
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mercul3s/placechicken/placer"
	"github.com/mercul3s/placechicken/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	if _, err := os.Stat("/tmp/placechicken"); os.IsNotExist(err) {
		err := os.Mkdir("/tmp/placechicken", 0700)
		if err != nil {
			t.Fatal(err)
		}
	}
	d := &testHelpers.Dir{}
	p := placer.Place{
		Dir:              d,
		OriginalFilePath: "../static/images/test/",
		ResizedFilePath:  "/tmp/placechicken/",
	}

	r := NewMux(p, "../static/", "../templates/")

	tt := []struct {
		name           string
		route          string
		expectedStatus int
		expectedHeader []string
		expectedBody   string
		expectedError  error
	}{
		{
			name:           "expect GET to '/' return the index page",
			route:          "/",
			expectedStatus: 200,
			expectedBody:   "Place Chicken",
		},
		{
			name:           "expect GET to '/static' to return a list of files",
			route:          "/static/",
			expectedStatus: 200,
			expectedBody:   "images",
		},
		{
			name:           "expect route not found to return 404 page",
			route:          "/bogus",
			expectedStatus: 404,
			expectedBody:   "( a )",
		},
	}
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			// get test image
			fileInfo, err := os.Stat("../static/images/test/original-test-image.jpg")
			fileList := []os.FileInfo{fileInfo}
			if err != nil {
				t.Fatal(err)
			}
			d.On("List", "../static/images/test/").Return(fileList, test.expectedError)
			req := httptest.NewRequest("GET", test.route, nil)
			rr := httptest.NewRecorder()
			r.Router.ServeHTTP(rr, req)
			bodyBytes, err := ioutil.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}
			bodyString := string(bodyBytes)
			assert.Equal(t, rr.Code, test.expectedStatus, test.name)
			assert.Contains(t, bodyString, test.expectedBody)
		})
	}
	err := os.RemoveAll("/tmp/placechicken")
	if err != nil {
		t.Fatal(err)
	}
}

func TestImageRoute(t *testing.T) {
	if _, err := os.Stat("/tmp/placechicken"); os.IsNotExist(err) {
		err := os.Mkdir("/tmp/placechicken", 0700)
		if err != nil {
			t.Fatal(err)
		}
	}

	tt := []struct {
		name           string
		route          string
		expectedStatus int
		expectedHeader []string
		expectedError  error
	}{
		{
			name:           "expect GET to '/{height}/{width} to return a random resized image'",
			route:          "/300/500",
			expectedStatus: 200,
			expectedHeader: []string{"image/jpeg"},
			expectedError:  nil,
		},
		{
			name:           "expect GET to return an error with image processing",
			route:          "/300/500",
			expectedStatus: 500,
			expectedHeader: []string{"text/plain; charset=utf-8"},
			expectedError:  errors.New("this is an error"),
		},
	}
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			// get test image
			fileInfo, err := os.Stat("../static/images/test/original-test-image.jpg")
			fileList := []os.FileInfo{fileInfo}
			if err != nil {
				t.Fatal(err)
			}
			d := &testHelpers.Dir{}
			p := placer.Place{
				Dir:              d,
				OriginalFilePath: "../static/images/test/",
				ResizedFilePath:  "/tmp/placechicken/",
			}

			d.On("List", "../static/images/test/").Return(fileList, test.expectedError)
			r := NewMux(p, "../static/", "../templates/")
			req := httptest.NewRequest("GET", test.route, nil)
			rr := httptest.NewRecorder()
			r.Router.ServeHTTP(rr, req)
			assert.Equal(t, rr.Code, test.expectedStatus, test.name)
			assert.Equal(t, rr.HeaderMap["Content-Type"], test.expectedHeader)
		})
	}
	err := os.RemoveAll("/tmp/placechicken")
	if err != nil {
		t.Fatal(err)
	}
}
