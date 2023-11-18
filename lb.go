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
	mu           sync.Mutex
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
	mu:     sync.Mutex{},
}

func getNextServerRoundRobin() string {
	lbObject.mu.Lock()
	serverEndpoint := lbObject.servers[lbObject.index%len(lbObject.servers)]
	lbObject.index++
	lbObject.mu.Unlock()
	return serverEndpoint
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var wg sync.WaitGroup
	var mu sync.Mutex
	var responses []string

	serverEndpoint := getNextServerRoundRobin()

	if _, ok := lbObject.serverStatus[serverEndpoint]; !ok {
		serverEndpoint = getNextServerRoundRobin()
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
				lbObject.mu.Lock()
				delete(lbObject.serverStatus, server)
				lbObject.mu.Unlock()
				continue
			}

			if response.StatusCode != http.StatusOK {
				lbObject.mu.Lock()
				delete(lbObject.serverStatus, server)
				lbObject.mu.Unlock()
			} else {
				lbObject.mu.Lock()
				lbObject.serverStatus[server] = true
				lbObject.mu.Unlock()
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
