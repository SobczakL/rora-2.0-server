package main

import (
	"log"
	"net/http"

	"rora-server/routes"
)

func main() {
	router := http.NewServeMux()

	routes.RegisterAPIRoutes(router)

	port := ":8080"
	log.Printf("Server running on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
