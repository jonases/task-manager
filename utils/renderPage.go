package utils

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jonases/cybersecuryproject/models"

	"github.com/josephspurrier/csrfbanana"
)

// RenderPage renders a specific web page
func RenderPage(res http.ResponseWriter, req *http.Request, list ...string) {

	sess := NewSession(req)

	// create a view
	v := NewView(req)
	// creates a token to protect against CSRF
	v.Vars["token"] = csrfbanana.Token(res, req, sess)

	// send back the attributes sent from the client. eg. Email, First Name
	// it repopulates the form fields if the validation did not pass
	if len(list) > 0 {
		Repopulate(list, req.Form, v.Vars)
	}

	// Get the flashes for the template
	if flashes := sess.Flashes(); len(flashes) > 0 {
		v.Vars["flashes"] = make([]Flash, len(flashes))
		for i, f := range flashes {
			v.Vars["flashes"].([]Flash)[i] = f.(Flash)
		}
		err := sess.Save(req, res)
		if err != nil {
			log.Println(err)
			http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}

	path := req.URL.EscapedPath()[1:]

	if path == "" || path == "home" {
		path = "index"
	}

	// populates with the templates and static HTML pages
	// so that we do not need to restart the service at every file update
	staticPages := PopulateStaticPages()
	staticPage := staticPages.Lookup(path + ".html")

	if staticPage == nil {
		log.Println("Page Not Found:", path)
		staticPage = staticPages.Lookup("404.html")
		res.WriteHeader(http.StatusNotFound)
	}

	v.Vars["Section"], v.Vars["Title"] = CreateContext(path)

	err := staticPage.Execute(res, v.Vars)
	if err != nil {
		log.Println(err)
	}
}

// PopulateStaticPages finds all HTML pages in the given directory
func PopulateStaticPages() *template.Template {
	//log.Println("Invoking populateStaticPages")

	tmpl := template.New("templates")
	templatePaths := new([]string)

	basePath := models.Path + models.Public + models.TemplatesPath

	templateFolder, err := os.Open(basePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer templateFolder.Close()

	templatePathsRaw, err := templateFolder.Readdir(0)
	if err != nil {
		log.Fatalln(err)
	}

	for _, pathInfo := range templatePathsRaw {
		//log.Println(pathInfo.Name())
		*templatePaths = append(*templatePaths, basePath+"/"+pathInfo.Name())
	}

	tmpl, err = tmpl.ParseFiles(*templatePaths...)
	if err != nil {
		log.Fatalln(err)
	}
	return tmpl
}

// CreateContext creates the context that will be used in the HTML page/template
func CreateContext(page string) (section, title string) {
	//log.Println("Invoking createContext")
	switch {
	case page == "index":
		title = "Home"
		section = page
	case page == "about":
		title = "About"
		section = page
	case page == "contact":
		title = "Contact"
		section = page
	case page == "login":
		title = "Login"
		section = page
	case page == "register":
		title = "Register"
		section = page
	default:
		title = page
		section = page
	}

	return
}
