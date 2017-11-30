package router

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mercul3s/placechicken/resizer"
)

// NewRouter returns a new mux router with all routes defined
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/egg", eggHandler).Methods("GET")
	r.HandleFunc("/{height}/{width}", resizeHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	return r
}

func resizeHandler(w http.ResponseWriter, r *http.Request) {
	resizer.ImageResize(w, r)
}

func eggHandler(w http.ResponseWriter, r *http.Request) {
	chicken, err := ioutil.ReadFile("../static/chicken")
	if err != nil {
		fmt.Println("error reading file", err)
	}

	fmt.Fprintf(w, string(chicken))
}
