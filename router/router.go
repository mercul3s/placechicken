package router

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mercul3s/placechicken_go/placer"
)

var p = placer.Config()

// NewRouter returns a new mux router with all routes defined
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/{height}/{width}", resizeHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../static/"))))
	r.NotFoundHandler = http.HandlerFunc(eggHandler)
	return r
}

func resizeHandler(w http.ResponseWriter, r *http.Request) {
	image, err := p.GetImage(500, 300)
	if err != nil {
		fmt.Println("error resizing file")
	}
	fmt.Fprintf(w, image)
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
