package handlers

import (
	"html"
	"log"
	"net/http"

	"github.com/jonases/task-manager/shared"
)

// Login handles POST requets made from the Login page
func Login(res http.ResponseWriter, req *http.Request) {

	email := html.EscapeString(req.FormValue("email"))
	//password := html.EscapeString(req.FormValue("password"))
	session := shared.NewSession(req)
	// create connection to users db
	shared.CreateDBConnection(shared.UsersDB)
	// query by email to find existing user
	result, err := shared.QueryByFieldAndValue("email", email)
	if err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// if no result came back, there was no match for the email meaning user does not exist
	if len(result) == 0 {
		log.Println("User does not exist:", email)
		session.AddFlash(shared.Flash{Message: "User " + req.FormValue("email") + " does not exit", Class: shared.FlashError})
		err := session.Save(req, res)
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		shared.RenderPage(res, req, "email")
		return
	}

	for _, v := range result {
		r := v.(map[string]interface{})
		// check if account is active
		if r["account_active"].(bool) {
			pwdMatch := shared.MatchString(r["password"].(string), req.FormValue("password"))
			if pwdMatch {
				// login successful
				shared.Empty(session)
				session.Values["email"] = email
				session.Values["fname"] = r["fname"].(string)
				session.AddFlash(shared.Flash{Message: "Welcome, " + html.UnescapeString(session.Values["fname"].(string)), Class: shared.FlashSuccess})
				err := session.Save(req, res)
				if err != nil {
					log.Println(err)
					http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				http.Redirect(res, req, "/", http.StatusFound)
				return
			}
			// if incorrect password was used return error message
			log.Println("Wrong password:", email)
			session.AddFlash(shared.Flash{Message: "Wrong password, please try again", Class: shared.FlashError})
			err := session.Save(req, res)
			if err != nil {
				log.Println(err)
				http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			shared.RenderPage(res, req, "email")
			return

		}
		// Account is inactive
		session.AddFlash(shared.Flash{Message: "Account " + email + " is disabled", Class: shared.FlashNotice})
		err := session.Save(req, res)
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		shared.RenderPage(res, req)
		return

	}

	// renders the login page
	shared.RenderPage(res, req)

}
