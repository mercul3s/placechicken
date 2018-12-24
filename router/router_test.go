package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	r := NewRouter()

	tt := []struct {
		name           string
		route          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "expect GET to '/' return the index page",
			route:          "/",
			expectedStatus: 200,
			expectedBody:   "hello world",
		},
		{
			name:           "expect GET to '/{height}/{width} to return a random resized image'",
			route:          "/300/500",
			expectedStatus: 200,
			expectedBody:   "image.jpg",
		},
		{
			name:           "expect GET to '/static' to return a list of files",
			route:          "/static/",
			expectedStatus: 200,
			expectedBody:   "image.jpg",
		},
		{
			name:           "expect route not found to return 404 page",
			route:          "/bogus",
			expectedStatus: 404,
			expectedBody:   "nil",
		},
	}
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", test.route, nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(t, test.expectedStatus, rr.Code, test.name)
		})
	}
}

func TestStaticServer(t *testing.T) {
	assert := assert.New(t)
	r := NewRouter()
	mockServer := httptest.NewServer(r)

	response, err := http.Get(mockServer.URL + "/static/chicken")
	if err != nil {
		assert.FailNow("error serving assets", err)
	}
	assert.Equal(http.StatusOK, response.StatusCode)
}
