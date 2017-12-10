package handlers

import (
	"db"
	"html"
	"log"
	"models"
	"net/http"
	"utils"
)

// ServeContent retrieves and serves the HTML pages
func ServeContent(res http.ResponseWriter, req *http.Request) {
	//log.Println("Invoking serveContent")
	session := utils.NewSession(req)

	if req.URL.EscapedPath() == "/logout" {

		// if user is authenticated
		if session.Values["email"] != nil {
			utils.Empty(session)
			session.AddFlash(utils.Flash{Message: "Successfully logged out!", Class: utils.FlashNotice})
			err := session.Save(req, res)
			if err != nil {
				log.Println(err)
				http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(res, req, "/", http.StatusFound)
		return
	}

	if req.URL.EscapedPath() == "/login" {

		// if user is authenticated, do not allow access to "/login" endpoint
		if session.Values["email"] != nil {
			http.Redirect(res, req, "/", http.StatusFound)
			return
		}
	}

	if req.URL.EscapedPath() == "/messages" {

		// if user is not authenticated, do not allow access to "/messages" endpoint
		if session.Values["email"] == nil {
			http.Redirect(res, req, "/", http.StatusFound)
			return
		}
		// creates db connection to "messages" document
		db.CreateDBConnection(models.MessagesDB)
		// get all content in the "messages" document
		err := db.GetAllDocs()
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var msgSlice []models.MsgData
		var msg models.MsgData

		for k := range db.Alldocs.Rows {
			doc := db.Alldocs.Rows[k]

			msg.Email = html.UnescapeString(doc["doc"].(map[string]interface{})["email"].(string))
			msg.Message = html.UnescapeString(doc["doc"].(map[string]interface{})["msg"].(string))
			msg.Name = html.UnescapeString(doc["doc"].(map[string]interface{})["name"].(string))
			msgSlice = append(msgSlice, msg)
		}

		session.Values["messages"] = true
		utils.Messages = msgSlice
		err = session.Save(req, res)
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	utils.RenderPage(res, req)

}
