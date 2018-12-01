package handlers

import (
	"log"
	"net/http"
)

// InvalidToken replies with "bad request" because token received was invalid
func InvalidToken(res http.ResponseWriter, req *http.Request) {

	log.Println("invalid token provided for the page, bad request.")

	http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
