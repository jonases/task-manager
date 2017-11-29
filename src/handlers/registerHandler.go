package handlers

import (
	"db"
	"html"
	"log"
	"models"
	"net/http"
	"strings"
	"utils"
)

// Register handles POST requests made from the Register page
func Register(res http.ResponseWriter, req *http.Request) {

	// creates connection to "users" document
	db.CreateDBConnection(models.UsersDB)
	// query by email address
	result := db.Query("email", html.EscapeString(req.FormValue("email")))

	// deny registering if email is found in the db
	if len(result) > 0 {
		log.Println("Email " + req.FormValue("email") + " is already registed")
		session := utils.NewSession(req)
		session.AddFlash(utils.Flash{Message: "Email " + req.FormValue("email") + " is already registed", Class: utils.FlashError})
		err := session.Save(req, res)
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		utils.RenderPage(res, req, "first_name", "last_name", "email")
		return
	}

	// validate required fields
	if result := utils.Validate(req, []string{"first_name", "last_name", "email", "password", "password_verify"}); len(result) > 0 {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Missing fields:" + strings.Join(result, ",")))
		return
	}

	// check whether passwords match
	if html.EscapeString(req.FormValue("password")) != html.EscapeString(req.FormValue("password_verify")) {
		log.Println("Passwords don't match. Email:", req.FormValue("email"))
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Passwords dont match"))
		return
	}

	var user models.UserDB
	var err error

	user.Email = html.EscapeString(req.FormValue("email"))
	// bcrypt hash password
	user.Password, err = utils.HashString(req.FormValue("password"))
	if err != nil {
		log.Println(err)
	}
	user.FirstName = req.FormValue("first_name")
	user.LastName = req.FormValue("last_name")
	user.AccountActive = true

	// creates user in "users" document
	err = db.CreateDocument(user)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	session := utils.NewSession(req)
	session.AddFlash(utils.Flash{Message: "User " + req.FormValue("email") + " successfully registered", Class: utils.FlashSuccess})
	err = session.Save(req, res)
	if err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// renders contact page, after the request
	utils.RenderPage(res, req)

}
