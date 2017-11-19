package handlers

import (
	"html/template"
	"log"
	"models"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// ServeContent retrieves and serves the HTML pages
func ServeContent(res http.ResponseWriter, req *http.Request) {
	log.Println("Invoking serveContent")
	urlParams := mux.Vars(req)
	log.Println("urlParams:", urlParams)
	pageAlias := urlParams["pageAlias"]
	//log.Println("pageAlias: ", pageAlias)
	if pageAlias == "" || pageAlias == "home" {
		pageAlias = "index"
	}

	// updates the templates and static HTML pages
	// so that we do not need to restart the service at every file update
	staticPages := populateStaticPages()
	staticPage := staticPages.Lookup(pageAlias + ".html")
	if staticPage == nil {
		log.Println("Page Not Found:", pageAlias)
		staticPage = staticPages.Lookup("404.html")
		res.WriteHeader(http.StatusNotFound)
	}

	context := createContext(pageAlias)

	err := staticPage.Execute(res, context)
	if err != nil {
		log.Fatalln(err)
	}
}

func populateStaticPages() *template.Template {
	log.Println("Invoking populateStaticPages")

	tmpl := template.New("templates")
	templatePaths := new([]string)

	basePath := models.TemplatesPath
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
		log.Println(pathInfo.Name())
		*templatePaths = append(*templatePaths, basePath+"/"+pathInfo.Name())

	}

	basePath = models.StaticHTML
	templateFolder, err = os.Open(basePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer templateFolder.Close()
	templatePathsRaw, err = templateFolder.Readdir(0)
	if err != nil {
		log.Fatalln(err)
	}
	for _, pathInfo := range templatePathsRaw {
		log.Println(pathInfo.Name())
		*templatePaths = append(*templatePaths, basePath+"/"+pathInfo.Name())

	}

	tmpl, err = tmpl.ParseFiles(*templatePaths...)
	if err != nil {
		log.Fatalln(err)
	}
	return tmpl
}

func createContext(page string) (context models.Context) {
	log.Println("Invoking createContext")

	switch {
	case page == "index":
		context.Title = "Home"
		context.Section = page
	case page == "about":
		context.Title = "About"
		context.Section = page
	case page == "contact":
		context.Title = "Contact"
		context.Section = page
	default:
		context.Title = page
		context.Section = page
	}

	return
}
