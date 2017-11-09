package main

import (
	"fmt"
	// "fmt"
	// "io"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	dec, err := w.Write([]byte("This is an example server.\n"))
	fmt.Printf("Dec: %d, err: %s", dec, err)
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func root(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	dec, err := w.Write([]byte("This is the root.\n\n\n\n\t\t\tWelcome!"))
	fmt.Printf("Dec: %d, err: %s", dec, err)

}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
