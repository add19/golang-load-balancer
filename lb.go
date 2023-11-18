package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type RoundRobinScheduler struct {
	index   int
	servers []string
	client  *http.Client
}

var lbObject = RoundRobinScheduler{
	index:   0,
	servers: []string{"http://127.0.0.1:8080/", "http://127.0.0.1:8081/", "http://127.0.0.1:8082/"},
	client:  &http.Client{},
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var wg sync.WaitGroup
	var mu sync.Mutex
	var responses []string
	serverEndpoint := lbObject.servers[lbObject.index%len(lbObject.servers)]
	lbObject.index++

	for range lbObject.servers {
		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()
			request, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			response, err := lbObject.client.Do(request)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer response.Body.Close()

			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			responses = append(responses, fmt.Sprintf("Received request from %s \nGET / HTTP/1.1\nHost: %s\nUser-Agent: %s\nAccept: */*\n%s", r.RemoteAddr, r.Host, r.UserAgent(), responseBody))
		}(serverEndpoint)
	}

	wg.Wait()

	// Writing collected responses to the ResponseWriter outside of goroutines
	for _, resp := range responses {
		fmt.Fprintf(w, "%s\n", resp)
	}
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var wg sync.WaitGroup
// 	serverEndpoint := lbObject.servers[lbObject.index%len(lbObject.servers)]
// 	lbObject.index++

// 	for range lbObject.servers {
// 		wg.Add(1)
// 		go func(serverEndpoint string) {
// 			defer wg.Done()
// 			request, err := http.NewRequest("GET", serverEndpoint, nil)

// 			if err != nil {
// 				fmt.Println(err)
// 				return
// 			}

// 			response, error := lbObject.client.Do(request)

// 			if error != nil {
// 				fmt.Println(error)
// 			}
// 			responseBody, error := io.ReadAll(response.Body)

// 			if error != nil {
// 				fmt.Println(error)
// 			}
// 			defer response.Body.Close()
// 			fmt.Fprintf(w, "Received request from %s \nGET / HTTP/1.1\nHost: %s\nUser-Agent: %s\nAccept: */*\n%s", r.RemoteAddr, r.Host, r.UserAgent(), responseBody)
// 		}(serverEndpoint)
// 	}
// 	wg.Wait()
// }

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
