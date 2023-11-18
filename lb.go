package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type RoundRobinScheduler struct {
	index        int
	servers      []string
	serverStatus map[string]bool
	client       *http.Client
}

var lbObject = RoundRobinScheduler{
	index:   0,
	servers: []string{"http://127.0.0.1:8080/", "http://127.0.0.1:8081/", "http://127.0.0.1:8082/"},
	serverStatus: map[string]bool{
		"http://127.0.0.1:8080/": true,
		"http://127.0.0.1:8081/": true,
		"http://127.0.0.1:8082/": true,
	},
	client: &http.Client{},
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var wg sync.WaitGroup
	var mu sync.Mutex
	var responses []string
	serverEndpoint := lbObject.servers[lbObject.index%len(lbObject.servers)]
	lbObject.index++

	if _, ok := lbObject.serverStatus[serverEndpoint]; !ok {
		serverEndpoint = lbObject.servers[lbObject.index%len(lbObject.servers)]
		lbObject.index++
	}

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

	wg.Wait()

	for _, resp := range responses {
		fmt.Fprintf(w, "%s\n", resp)
	}
}

func healthCheck() {
	for {
		time.Sleep(10 * time.Second)
		fmt.Println("Checking servers")
		for _, server := range lbObject.servers {
			response, err := http.Get(server)
			if err != nil {
				delete(lbObject.serverStatus, server)
				continue
			}

			if response.StatusCode != http.StatusOK {
				delete(lbObject.serverStatus, server)
			} else {
				lbObject.serverStatus[server] = true
			}
			response.Body.Close()
		}
	}
}

func main() {
	go healthCheck()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
