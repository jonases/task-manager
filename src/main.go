package main

import (
	"handlers"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.ServeContent).Methods("GET")
	router.HandleFunc("/{pageAlias}", handlers.ServeContent).Methods("GET")

	http.HandleFunc("/js/", handlers.ServeResource)
	http.HandleFunc("/css/", handlers.ServeResource)
	http.Handle("/", router)

	srv := &http.Server{
		Handler:      nil,
		Addr:         ":8001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
