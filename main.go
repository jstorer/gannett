package main

import (
	"log"
	"net/http"
	"github.com/jstorer/gannett/api"
	"fmt"
)

func main() {
	fmt.Println("...Supermarket Server Starting...")
	api.Initialize()
	log.Fatal(http.ListenAndServe(":8080", api.Handlers()))

}
