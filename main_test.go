package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jonases/cybersecuryproject/db"
	"github.com/jonases/cybersecuryproject/handlers"
	"github.com/jonases/cybersecuryproject/models"
)

var (
	token  string
	cookie *http.Cookie
)

func TestInitSetup(t *testing.T) {
	// initializes cookie store
	initSessionStore()
	// creates a client to be used to connect to the Cloudant database
	db.CloudantInit()
	// append "_test" to the database names
	models.UsersDB += "_test"
	models.MessagesDB += "_test"
	// set the base path
	models.Path = ""
	// set up the logging output standard
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

func TestGetContactPage(t *testing.T) {
	Get(t, "/contact", handlers.ServeContent, http.StatusOK, "")
}

func TestPostContactPage(t *testing.T) {
	postBody := strings.NewReader("name=Test Name&email=myemail@email.com&message=This is a test msg")

	Post(t, "/contact", postBody, handlers.SendMessage, http.StatusOK, "")
}

func TestGetRegisterPage(t *testing.T) {
	Get(t, "/register", handlers.ServeContent, http.StatusOK, "")
}

func TestPostRegisterPage(t *testing.T) {
	postBody := strings.NewReader("first_name=First Name Test&last_name=Last Name Test&email=myemail@email.com&password=ThisIsMyPassword99&password_verify=ThisIsMyPassword99")

	Post(t, "/register", postBody, handlers.Register, http.StatusOK, "")
}

func TestGetLoginPage(t *testing.T) {
	Get(t, "/login", handlers.ServeContent, http.StatusOK, "")
}

func TestPostLoginPage(t *testing.T) {
	postBody := strings.NewReader("email=myemail@email.com&password=ThisIsMyPassword99")

	Post(t, "/login", postBody, handlers.Login, http.StatusFound, "")
}

func TestGetMessages(t *testing.T) {
	Get(t, "/messages", handlers.ServeContent, http.StatusOK, "")
}

func TestLogout(t *testing.T) {
	Get(t, "/logout", handlers.ServeContent, http.StatusFound, "")
}

func TestCleanUpDB(t *testing.T) {
	db.DeleteDocument(db.MsgsDoc.ID, db.MsgsDoc.Rev)
	db.CreateDBConnection(models.UsersDB)
	db.DeleteDocument(db.UsersDoc.ID, db.UsersDoc.Rev)
}

// Get defines the GET requests, expecting the desired HTTP code and body
func Get(t *testing.T, url string, hFunc http.HandlerFunc, expectedStatus int, expectedBody string) {

	// Create a request to pass to the handler
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	// set cookie header if it exists
	if cookie != nil {
		req.AddCookie(cookie)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(hFunc)

	// the handler satisfy http.Handler, so we can call ServeHTTP method
	// directly and pass in Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what's expected
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v, want %v",
			status, expectedStatus)
	}

	if expectedBody != "" {
		// Check the response body is what we expect.
		expected := string(expectedBody)
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v, want %v",
				rr.Body.String(), expected)
		}
	}
}

// Post defines the POST requests, expecting the desired HTTP code and body
func Post(t *testing.T, url string, body *strings.Reader, hFunc http.HandlerFunc, expectedStatus int, expectedBody string) {

	// Create a request to pass to the handler
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		t.Fatal(err)
	}

	// set the correct content type header
	req.Header.Set("Content-Type", " application/x-www-form-urlencoded")

	// populate the form values
	err = req.ParseForm()
	if err != nil {
		t.Error(err)
	}
	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(hFunc)

	// the handler satisfy http.Handler, so we can call ServeHTTP method
	// directly and pass in the Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// save the cookie returned after logging in
	if url == "/login" {
		result := rr.Result()
		cookies := result.Cookies()
		cookie = cookies[0]
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v, want %v",
			status, expectedStatus)
	}

	// Check the response body is what we expect.
	if expectedBody != "" {
		expected := string(expectedBody)
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v, want %v",
				rr.Body.String(), expected)
		}
	}
}
