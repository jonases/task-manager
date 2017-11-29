package handlers

import (
	"db"
	"html"
	"log"
	"models"
	"net/http"
	"utils"
)

// SendMessage handles POST requets made from the Contact page
func SendMessage(res http.ResponseWriter, req *http.Request) {
	//log.Println("Invoking SendMessage")

	// creates session
	session := utils.NewSession(req)

	var message models.MsgData

	message.Name = html.EscapeString(req.FormValue("name"))
	message.Email = html.EscapeString(req.FormValue("email"))
	message.Message = html.EscapeString(req.FormValue("message"))

	// create db connection to "messages" db
	db.CreateDBConnection(models.MessagesDB)
	// creates the document in the database
	err := db.CreateDocument(message)
	if err != nil {
		log.Println(err)
		session.AddFlash(utils.Flash{Message: "Error when contacting database, please try again later", Class: utils.FlashError})
		err := session.Save(req, res)
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		utils.RenderPage(res, req, "email", "name", "message")
		return
	}

	session.AddFlash(utils.Flash{Message: "Message sucessfully sent!", Class: utils.FlashSuccess})
	err = session.Save(req, res)
	if err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// renders contact page, after the request
	utils.RenderPage(res, req)

}
