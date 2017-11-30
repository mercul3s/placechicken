package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	assert := assert.New(t)
	r := NewRouter()
	mockServer := httptest.NewServer(r)

	response, err := http.Get(mockServer.URL + "/egg")
	if err != nil {
		assert.FailNow("Error requesting url from server", err)
	}

	chickenAsset, err := ioutil.ReadFile("../static/chicken")
	if err != nil {
		assert.FailNow("error reading file")
	}
	assert.Equal(http.StatusOK, response.StatusCode)
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		assert.FailNow("error reading response body", err)
	}
	assert.Equal(string(chickenAsset), string(b))
	fmt.Println(string(b))
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
