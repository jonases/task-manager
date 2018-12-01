package handlers

import (
	"html"
	"log"
	"net/http"

	"github.com/jonases/task-manager/shared"
)

// ServeAddTodos handles requests for /todo endpoint
func ServeAddTodos(res http.ResponseWriter, req *http.Request) {

	log.Println("Invoking ServeTodos")

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

	todo.Title = html.EscapeString(req.FormValue("task"))
	todo.State = "in-progress"
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

	err := shared.InsertTodo(todo)
	if err != nil {
		log.Println(err)
		http.Error(res, "error inserting todo into database", http.StatusInternalServerError)
		return
	}

	session.Values["todos"] = true
	shared.Todos = []shared.Todo{todo}

	err = session.Save(req, res)
	if err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(res, req, "/todo", http.StatusFound)

}
