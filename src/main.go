package main

import (
	"db"
	"encoding/gob"
	"handlers"
	"log"
	"models"
	"net/http"
	"time"
	"utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	// setup session settings
	session := &utils.Session{
		SecretKey: "K1fuGi8xsXihxiDbxiL4z2A80",
		Name:      models.CookieName,
		Options: sessions.Options{
			Path:     "/",
			Domain:   "",
			MaxAge:   86400,
			Secure:   false,
			HttpOnly: true,
		},
	}
	// need to register structures to be used in the sessions
	gob.Register(utils.Flash{})
	gob.Register([]models.MsgData{})

	// Configures the session cookie store
	utils.Configure(*session)

	// creates a client to be used to connect to the Cloudant database
	db.CloudantInit()

	// set up the routes for the HTTP handle
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.ServeContent).Methods("GET")
	router.HandleFunc("/{pageAlias}", handlers.ServeContent).Methods("GET")
	router.HandleFunc("/contact", handlers.SendMessage).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFound)

	http.HandleFunc("/js/", handlers.ServeResource)
	http.HandleFunc("/css/", handlers.ServeResource)

	http.Handle("/", router)

	srv := &http.Server{
		Handler:      nil,
		Addr:         ":8001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// start the HTTP server
	log.Fatal(srv.ListenAndServe())

}
