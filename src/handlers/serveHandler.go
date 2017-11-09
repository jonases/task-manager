package handlers

import (
	"log"
	"net/http"
)

var (
	contentType = "application/json"
)

// ServeSubscribe will handle something
func ServeSubscribe(res http.ResponseWriter, req *http.Request) {
	log.Println("Invoking ServeSubscribe handler")

	if req.Method != http.MethodGet {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusForbidden)
		return
	}

	res.Header().Set("Content-Type", contentType)

}
