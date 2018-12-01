package handlers

import (
	"net/http"

	"github.com/jonases/task-manager/shared"
)

// NotFound handler returns the pre-defined 404 page together with 404 HTTP status
func NotFound(res http.ResponseWriter, req *http.Request) {
	//log.Println("Invoking NotFound")

	// renders NotFound page
	shared.RenderPage(res, req)

}
