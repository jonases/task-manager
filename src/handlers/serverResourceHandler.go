package handlers

import (
	"bufio"
	"log"
	"models"
	"net/http"
	"os"
	"strings"
)

// ServeResource retrieves and serves the static pages such as .js and .css files
func ServeResource(res http.ResponseWriter, req *http.Request) {
	//log.Println("Invoking serveResource")

	// return method not allowed if attempt to use a method different than GET
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}

	path := models.Path + models.Public + "static" + req.URL.Path

	var contentType string

	if strings.HasSuffix(path, ".css") {
		contentType = "text/css; charset=utf-8"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "application/javascript; charset=utf-8"
	} else {
		contentType = "text/plain; charset=utf-8"
	}

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()

	res.Header().Add("Content-Type", contentType)

	buff := bufio.NewReader(file)
	_, err = buff.WriteTo(res)

	if err != nil {
		log.Fatalln(err)
	}
}
