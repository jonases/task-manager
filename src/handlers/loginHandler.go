package handlers

import (
	"db"
	"html"
	"log"
	"models"
	"net/http"
	"utils"
)

// Login handles POST requets made from the Login page
func Login(res http.ResponseWriter, req *http.Request) {
	//	log.Println("Password:", req.FormValue("password"))

	email := html.EscapeString(req.FormValue("email"))
	//password := html.EscapeString(req.FormValue("password"))
	session := utils.NewSession(req)
	// create connection to users db
	db.CreateDBConnection(models.UsersDB)
	// query by email to find existing user
	result := db.Query("email", email)
	// if no result came back, there was no match for the email meaning user does not exist
	if len(result) == 0 {
		log.Println("User does not exist:", email)
		session.AddFlash(utils.Flash{Message: "User " + req.FormValue("email") + " does not exit", Class: utils.FlashError})
		err := session.Save(req, res)
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		utils.RenderPage(res, req, "email")
		return
	}

	for _, v := range result {
		r := v.(map[string]interface{})
		// check if account is active
		if r["account_active"].(bool) {
			pwdMatch := utils.MatchString(r["password"].(string), req.FormValue("password"))
			if pwdMatch {
				// login successful
				utils.Empty(session)
				session.Values["email"] = email
				session.Values["fname"] = r["fname"].(string)
				session.AddFlash(utils.Flash{Message: "Welcome, " + html.UnescapeString(session.Values["fname"].(string)), Class: utils.FlashSuccess})
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
			session.AddFlash(utils.Flash{Message: "Wrong password, please try again", Class: utils.FlashError})
			err := session.Save(req, res)
			if err != nil {
				log.Println(err)
				http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			utils.RenderPage(res, req, "email")
			return

		}
		// Account is inactive
		session.AddFlash(utils.Flash{Message: "Account " + email + " is disabled", Class: utils.FlashNotice})
		err := session.Save(req, res)
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		utils.RenderPage(res, req)
		return

	}

	// renders the login page
	utils.RenderPage(res, req)

}
