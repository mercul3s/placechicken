package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mercul3s/placechicken/placer"
	"github.com/mercul3s/placechicken/router"
)

var logger = log.New(os.Stdout, "placechicken:", log.Lshortfile)

func main() {
	static := os.Getenv("STATIC")
	resized := os.Getenv("RESIZED")
	p := placer.Config(static, resized)
	logger.Printf("placechicken started with directories static: %s and resized: %s", static, resized)
	m := router.NewMux(p, "static/", "templates/")
	err := http.ListenAndServe(":8888", m.Router)
	if err != nil {
		logger.Print(err)
	}
}
