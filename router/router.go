package router

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"github.com/mercul3s/placechicken/placer"
)

// Mux holds the configuration information for an router.
type Mux struct {
	Place       placer.Place
	Router      *mux.Router
	staticDir   string
	templateDir string
}

// PageData stores information for output in a template.
type PageData struct {
	Image string
}

// NewMux returns a new mux router with all routes defined
func NewMux(place placer.Place, sDir string, tDir string) Mux {
	r := mux.NewRouter()
	m := Mux{
		Place:       place,
		Router:      r,
		staticDir:   sDir,
		templateDir: tDir,
	}
	m.Router.HandleFunc("/", m.index).Methods("GET")
	m.Router.HandleFunc("/{width}/{height}", m.resizeHandler).Methods("GET")
	m.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(sDir))))
	m.Router.NotFoundHandler = http.HandlerFunc(m.eggHandler)

	return m
}

func (m Mux) index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(m.templateDir + "index.html"))
	data := PageData{Image: "/605/0"}
	err := t.Execute(w, data)
	if err != nil {
		msg := fmt.Sprintf("unable to process request: %s", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func (m Mux) resizeHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	width, wErr := strconv.Atoi(v["width"])
	height, hErr := strconv.Atoi(v["height"])
	if wErr != nil && hErr != nil {
		msg := fmt.Sprintf("unable to process request: must provide width or height")
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	image, err := m.Place.GetImage(width, height)
	if err != nil {
		msg := fmt.Sprintf("unable to process request: %s", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	imaging.Encode(w, image, imaging.JPEG)
}

func (m Mux) eggHandler(w http.ResponseWriter, r *http.Request) {
	chicken, err := ioutil.ReadFile(m.templateDir + "chicken")
	if err != nil {
		msg := fmt.Sprintf("error reading file: %s", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, string(chicken))
}
