package main

import (
	"log"
	"net/http"
	"workshop/routes"
)

func main() {

	r := routes.SetupRouter()

	// Start listen to incoming requests.
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Unable to start:", err)
	}
}
