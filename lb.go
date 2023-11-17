package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	go func() {
		w.Header().Set("Content-Type", "application/json")

		apiUrl := "http://127.0.0.1:8080/"
		request, err := http.NewRequest("GET", apiUrl, nil)

		if err != nil {
			fmt.Println(err)
			return
		}

		client := &http.Client{}
		response, error := client.Do(request)

		if error != nil {
			fmt.Println(error)
		}
		responseBody, error := io.ReadAll(response.Body)

		if error != nil {
			fmt.Println(error)
		}
		fmt.Fprintf(w, "Received request from %s \nGET / HTTP/1.1\nHost: %s\nUser-Agent: %s\nAccept: */* \n%s", r.RemoteAddr, r.Host, r.UserAgent(), responseBody)
	}()
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
