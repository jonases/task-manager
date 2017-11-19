package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var (
	stacktraceBuffer []byte
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	//w.Write([]byte("This is an example server.\n"))
	dec, err := w.Write([]byte("This is an example server.\n"))
	fmt.Printf("Dec: %d, err: %s", dec, err)
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func root(w http.ResponseWriter, req *http.Request) {
	//template.

	w.Header().Set("Content-Type", "text/plain")
	dec, err := w.Write([]byte("This is the root.\n\n\n\n\t\t\tWelcome!"))
	fmt.Printf("Dec: %d, err: %s", dec, err)

}

func main() {

	// preallocate array for stack trace
	if stacktraceBuffer == nil {
		stacktraceBuffer = make([]byte, 4096)
	}

	t := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	// Write startup info to console (stdout)
	log.Println("*** Starting web app...", t)

	// handle unix signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs

		fmt.Printf("SIGNAL Received: '%s'. Stopping service.\n", sig)
		log.Printf("SIGNAL Received: '%s'. Stopping service.\n", sig)
		time.Sleep(5 * time.Second)
		os.Exit(-3)
	}()

	// intercept panics: print error and stacktrace
	defer func() {
		if err := recover(); err != nil {
			count := runtime.Stack(stacktraceBuffer, true)

			fmt.Printf("PANIC: '%s'. Stopping service.\n", err)
			fmt.Printf("STACKTRACE (%d bytes): %s\n", count, stacktraceBuffer[:count])

			log.Printf("PANIC: '%s'. Stopping service.\n", err)
			log.Printf("STACKTRACE (%d bytes): %s\n", count, stacktraceBuffer[:count])

			time.Sleep(5 * time.Second)
			os.Exit(-2)
		}
	}()

	http.HandleFunc("/", root)
	http.HandleFunc("/hello", HelloServer)

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
		if err != nil {
			log.Println("web app aborted with", err)
			os.Exit(1)
		} else {
			fmt.Println("web app stopped listening.")
		}
	}()

	err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	if err != nil {
		log.Println("web app aborted with", err)
		os.Exit(1)
	} else {
		fmt.Println("web app stopped listening.")
	}

	log.Println("web app stopped")
	time.Sleep(5 * time.Second)
	os.Exit(0)
}
