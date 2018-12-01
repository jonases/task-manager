package handlers

import (
	"html"
	"log"
	"net/http"

	"github.com/jonases/task-manager/shared"
)

// ServeContent retrieves and serves the HTML pages
func ServeContent(res http.ResponseWriter, req *http.Request) {
	//log.Println("Invoking serveContent")
	session := shared.NewSession(req)

	if req.URL.EscapedPath() == "/register" {

		// if user is authenticated, do not let to register
		if session.Values["email"] != nil {
			session.AddFlash(shared.Flash{Message: "Can't create an account while logged in.", Class: shared.FlashNotice})
			err := session.Save(req, res)
			if err != nil {
				log.Println(err)
				http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			http.Redirect(res, req, "/", http.StatusFound)
			return
		}

	}

	if req.URL.EscapedPath() == "/logout" {

		// if user is authenticated
		if session.Values["email"] != nil {
			shared.Empty(session)
			log.Println("sess:", session.Values)
			log.Println("in here")
			session.AddFlash(shared.Flash{Message: "Successfully logged out!", Class: shared.FlashNotice})
			err := session.Save(req, res)
			if err != nil {
				log.Println(err)
				http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		log.Println("sess:", session.Values)

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

	if req.URL.EscapedPath() == "/todo" {

		// if user is not authenticated, do not allow access to "/todo" endpoint
		if session.Values["email"] == nil {
			http.Redirect(res, req, "/", http.StatusFound)
			return
		}

		log.Println("session:", session.Values)

		// creates db connection to "tasks" document
		shared.CreateDBConnection(shared.TasksDocumentName)

		// get all content in the "tasks" document
		todoList, err := shared.QueryByFieldAndValue("email", session.Values["email"].(string))
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// log.Println("todoList=\n", fmt.Sprintf("%+v", todoList))

		var todoSlice []shared.Todo
		var todo shared.Todo

		for _, v := range todoList {

			val := v.(map[string]interface{})
			todo.Title = html.UnescapeString(val["title"].(string))
			todo.State = val["state"].(string)
			todo.ID = val["_id"].(string)
			todo.Rev = val["_rev"].(string)

			todoSlice = append(todoSlice, todo)
		}

		session.Values["todos"] = true
		shared.Todos = todoSlice

		err = session.Save(req, res)
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	shared.RenderPage(res, req)

}
