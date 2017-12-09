package handlers

import "net/http"

// InvalidToken replies with "bad request" because token received was invalid
func InvalidToken(res http.ResponseWriter, req *http.Request) {
	http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	return
}
