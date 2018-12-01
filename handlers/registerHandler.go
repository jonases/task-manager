package handlers

import (
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/jonases/task-manager/shared"
)

// Register handles POST requests made from the Register page
func Register(res http.ResponseWriter, req *http.Request) {

	log.Println("Invoking Register page")

	session := shared.NewSession(req)

	// creates connection to "users" document
	shared.CreateDBConnection(shared.UsersDB)

	email := html.EscapeString(req.FormValue("email"))
	password := html.EscapeString(req.FormValue("password"))

	// query by email address
	result, err := shared.QueryByFieldAndValue("email", email)
	if err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// deny registering if email is found in the db
	if len(result) > 0 {
		log.Println("Email " + req.FormValue("email") + " is already registed")
		session.AddFlash(shared.Flash{Message: "Email " + req.FormValue("email") + " is already registed", Class: shared.FlashError})
		err := session.Save(req, res)
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		shared.RenderPage(res, req, "first_name", "last_name", "email")
		return
	}

	// validate required fields
	if result := shared.Validate(req, []string{"first_name", "last_name", "email", "password", "password_verify"}); len(result) > 0 {
		res.WriteHeader(http.StatusBadRequest)
		if _, err := res.Write([]byte("Missing fields:" + strings.Join(result, ","))); err != nil {
			log.Println(err)
		}
		return
	}

	// check whether passwords match
	if password != html.EscapeString(req.FormValue("password_verify")) {
		log.Println("Passwords don't match. Email:", email)
		res.WriteHeader(http.StatusBadRequest)
		if _, err := res.Write([]byte("Passwords dont match")); err != nil {
			log.Println(err)
		}
		return
	}

	var user shared.UserDB

	user.Email = email
	// bcrypt hash password
	user.Password, err = shared.HashString(password)
	if err != nil {
		log.Println(err)
	}
	user.FirstName = html.EscapeString(req.FormValue("first_name"))
	user.LastName = html.EscapeString(req.FormValue("last_name"))
	user.AccountActive = true

	// creates user in "users" document
	err = shared.CreateDocument(user)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	session.AddFlash(shared.Flash{Message: "User " + req.FormValue("email") + " successfully registered", Class: shared.FlashSuccess})
	err = session.Save(req, res)
	if err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// renders contact page, after the request
	shared.RenderPage(res, req)

}
