package main

import (
	"log"
	"net/http"

	listdebates "github.com/raphink/debate/functions/list-debates"
)

func main() {
	http.HandleFunc("/", listdebates.ListDebatesHandler)

	port := "8080"
	log.Printf("list-debates server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
