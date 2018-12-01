package handlers

import (
	"html"
	"log"
	"net/http"

	"github.com/jonases/task-manager/shared"
)

// ServeUpdateTodos handles requests for /todo endpoint
func ServeUpdateTodos(res http.ResponseWriter, req *http.Request) {

	log.Println("Invoking ServeUpdateTodos")

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

	var todo shared.Todo
	todo.State = html.EscapeString(req.FormValue("state"))
	todo.Title = html.EscapeString(req.FormValue("title"))
	todo.Email = session.Values["email"].(string)

	if todo.Title == "" {
		log.Println("no task name was sent in the request")
		session.AddFlash(shared.Flash{Message: "No task name was sent in the request, please add one", Class: shared.FlashError})
		err := session.Save(req, res)
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		shared.RenderPage(res, req)
		return
	}

	err := shared.UpdateDocument(html.EscapeString(req.FormValue("id")), html.EscapeString(req.FormValue("rev")), todo)
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

	log.Println("successfully updated todo with id=" + req.FormValue("id"))

}
