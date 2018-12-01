package shared

import (
	"net/http"
	"net/url"
)

var (
	// FlashError is a bootstrap class
	FlashError = "alert-danger"
	// FlashSuccess is a bootstrap class
	FlashSuccess = "alert-success"
	// FlashNotice is a bootstrap class
	FlashNotice = "alert-info"
	// FlashWarning is a bootstrap class
	FlashWarning = "alert-warning"

	// Todos is sent back to the client
	Todos []Todo
)

// View attributes
type View struct {
	Vars map[string]interface{}
	//request *http.Request
	Title   string
	Section string
}

// Flash Message
type Flash struct {
	Message string
	Class   string
}

// NewView returns a new view
func NewView(req *http.Request) *View {
	v := &View{}
	v.Vars = make(map[string]interface{})

	// set auth level to annonymous by default
	v.Vars["AuthLevel"] = "anon"

	// Get session
	sess := NewSession(req)

	// used to populate the todo table in "/todo" endpoint
	if sess.Values["todos"] != nil {
		v.Vars["todos"] = Todos
	}

	// Set the AuthLevel to auth if the user is logged in
	if sess.Values["email"] != nil {
		v.Vars["AuthLevel"] = "auth"
	}

	return v
}

// Validate checks if required fields have been set, and return its values
func Validate(req *http.Request, required []string) (result []string) {
	for _, v := range required {
		if req.FormValue(v) == "" {
			result = append(result, v)
		}
	}

	return
}

// Repopulate updates the dst map so the form fields can be refilled
func Repopulate(list []string, src url.Values, dst map[string]interface{}) {
	for _, v := range list {
		dst[v] = src.Get(v)
	}
}
