package resizer

import (
	"fmt"
	"net/http"
)

// ImageResize takes an HTTP request and returns an image sized to the
// dimentions requested in the URL.
func ImageResize(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.String())
}
