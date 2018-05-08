package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/jstorer/gannett/api"
)


func main() {
	fmt.Println("...Supermarket Server Starting...")
	api.Initialize(false)
	log.Fatal(http.ListenAndServe(":8080", api.Handlers()))
}
