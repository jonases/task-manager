package handlers

import (
	"html"
	"log"
	"net/http"

	"github.com/jonases/task-manager/shared"
)

// ServeDeleteTodos handles requests for /todo endpoint
func ServeDeleteTodos(res http.ResponseWriter, req *http.Request) {

	log.Println("Invoking ServeDeleteTodos")

	session := shared.NewSession(req)

	defer req.Body.Close()

	if err := req.ParseForm(); err != nil {
		log.Println(err)
		session.AddFlash(shared.Flash{Message: "error parsing form", Class: shared.FlashError})
		err := session.Save(req, res)
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		shared.RenderPage(res, req)
		return
	}

	// log.Println(fmt.Sprintf("%+v", req.Form))

	err := shared.DeleteDocument(html.EscapeString(req.FormValue("id")), html.EscapeString(req.FormValue("rev")))
	if err != nil {
		log.Println(err)
		session.AddFlash(shared.Flash{Message: "error parsing form", Class: shared.FlashError})
		err := session.Save(req, res)
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		shared.RenderPage(res, req)
		return
	}

}
