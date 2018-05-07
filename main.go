package main

import (
	"fmt"
	"github.com/jstorer/gannett/api"
	"log"
	"net/http"
)

func main() {
	fmt.Println("...Supermarket Server Starting...")
	api.Initialize()
	log.Fatal(http.ListenAndServe(":8080", api.Handlers()))
}
