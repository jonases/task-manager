package main

import (
	"encoding/gob"
	"handlers"
	"log"
	"models"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"utils"

	"db"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/josephspurrier/csrfbanana"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	// initializes the cookie session store
	initSessionStore()

	// creates a client to be used to connect to the Cloudant database
	db.CloudantInit()

	// find full path of the current executable including the file name
	ex, e := os.Executable()
	if e != nil {
		log.Println(e)
		os.Exit(1)
	}

	// returns the path, excluding the file name
	models.Path = filepath.Dir(ex) + "/src/"

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

	// Set up the CSRF protection mechanism
	csrfProtection := csrfbanana.New(http.DefaultServeMux, utils.Store, utils.Name)
	// Generate a new token after each success/failure (also prevents double submits)
	csrfProtection.ClearAfterUsage(true)
	// Set a specific handler when an nvalid token is received
	csrfProtection.FailureHandler(http.HandlerFunc(handlers.InvalidToken))

	srv := &http.Server{
		Handler:      csrfProtection,
		Addr:         ":8001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// start the HTTP server
	log.Fatal(srv.ListenAndServe())

}

func initSessionStore() {
	// cookie store settings
	session := &utils.Session{
		// authenticate the cookie value using HMAC
		HashKey: []byte("F564O4sK16j8eEybQt2ht6DLehxuV4iHioUBsUwSpDU=vUjHATXHn8T89lX3Cg1"),
		// encryption key to encrypt the cookie using AES-256
		BlockKey: []byte("oQGCK9HFaQYAAJrmukcKclXN8WCL+yTs"),
		Name:     models.CookieName,
		Options: sessions.Options{
			Path:     "/",
			Domain:   "",
			MaxAge:   86400,
			Secure:   false,
			HttpOnly: true,
		},
	}

	// allow serializing maps in securecookie
	// http://golang.org/pkg/encoding/gob/#Register
	// need to register structure that are used in the cookie sessions, so that we can attach them to each session
	gob.Register(utils.Flash{})
	gob.Register(csrfbanana.StringMap{})

	// Configures the session cookie store
	utils.Configure(*session)
}
