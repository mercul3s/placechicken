package main

import (
	"fmt"
	"net/http"

	"github.com/mercul3s/placechicken/placer"
	"github.com/mercul3s/placechicken/router"
)

func main() {
	p := placer.Config("./static/images/", "./static/images/resized/")
	m := router.NewMux(p, "static/", "templates/")
	err := http.ListenAndServe(":8888", m.Router)
	if err != nil {
		fmt.Println(err)
	}
}
