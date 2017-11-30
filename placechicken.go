package main

import (
	"fmt"
	"net/http"

	"github.com/mercul3s/placechicken/router"
)

func main() {
	router := router.NewRouter()
	http.Handle("/", router)
	err := http.ListenAndServe(":8888", router)
	if err != nil {
		fmt.Println(err)
	}
}
