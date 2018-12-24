package router

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mercul3s/placechicken_go/resizer"
)

// NewRouter returns a new mux router with all routes defined
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/{height}/{width}", resizeHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.NotFoundHandler = http.HandlerFunc(eggHandler)
	return r
}

func resizeHandler(w http.ResponseWriter, r *http.Request) {
	resizer.ImageResize(w, r)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")

}

func eggHandler(w http.ResponseWriter, r *http.Request) {

	chicken, err := ioutil.ReadFile("../static/chicken")
	if err != nil {
		fmt.Println("error reading file", err)
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, string(chicken))

}
