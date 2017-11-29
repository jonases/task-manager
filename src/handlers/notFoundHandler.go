package handlers

import (
	"net/http"
	"utils"
)

// NotFound handler returns the pre-defined 404 page together with 404 HTTP status
func NotFound(res http.ResponseWriter, req *http.Request) {
	//log.Println("Invoking NotFound")

	// renders NotFound page
	utils.RenderPage(res, req)

}
