package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mercul3s/placechicken/placer"
)

// Mux holds the configuration information for an router.
type Mux struct {
	Place  placer.Place
	Router *mux.Router
}

// NewMux returns a new mux router with all routes defined
func NewMux(p placer.Place) Mux {
	r := mux.NewRouter()
	m := Mux{
		Place:  p,
		Router: r,
	}
	m.Router.HandleFunc("/", m.index).Methods("GET")
	m.Router.HandleFunc("/{width}/{height}", m.resizeHandler).Methods("GET")
	m.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../static/"))))
	m.Router.NotFoundHandler = http.HandlerFunc(eggHandler)

	return m
}

func (m Mux) index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func (m Mux) resizeHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	width, err := strconv.Atoi(v["width"])
	if err != nil {
		handleClientError(w, err)
		return
	}
	height, err := strconv.Atoi(v["height"])
	if err != nil {
		handleClientError(w, err)
		return
	}
	image, err := m.Place.GetImage(width, height)
	if err != nil {
		handleServerError(w, err)
		return
	}
	fmt.Fprintf(w, image)
}

func eggHandler(w http.ResponseWriter, r *http.Request) {
	chicken, err := ioutil.ReadFile("../static/chicken")
	if err != nil {
		fmt.Println("error reading file", err)
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, string(chicken))
}

func handleServerError(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusBadRequest)
	msg := fmt.Sprintf("unable to process request: %s", e)
	fmt.Fprintf(w, msg)
}

func handleClientError(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusInternalServerError)
	msg := fmt.Sprintf("unable to process request: %s", e)
	fmt.Fprintf(w, msg)
}
