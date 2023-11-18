package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	go func() {
		fmt.Fprintf(w, "Hello From Backend Server")
		w.Header().Set("Content-Type", "application/json")
	}()
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
