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

// func handle(conn net.Conn) {
// 	remoteHost, _ := net.LookupAddr(strings.Split(conn.RemoteAddr().String(), ":")[0])
// 	var b bytes.Buffer
// 	fmt.Fprintf(&b, "Received request from %s \nGET / HTTP/1.1\nHost: %s\nUser-Agent: %s\nAccept: */*\n\nReplied with a hello message", conn.RemoteAddr(), remoteHost, conn.LocalAddr())

// 	_, err := conn.Write(b.Bytes())
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// }

// func main() {
// 	listener, err := net.Listen("tcp", "localhost:8080")
// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 		return
// 	}
// 	defer listener.Close()

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Println("Error: ", err)
// 			continue
// 		}

// 		go handle(conn)
// 	}
// }
