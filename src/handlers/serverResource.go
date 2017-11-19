package handlers

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
)

// ServeResource retrieves and serves the static pages such as .js and .css files
func ServeResource(res http.ResponseWriter, req *http.Request) {
	log.Println("Invoking serveResource")
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := "static" + req.URL.Path
	log.Println("PATH:", path)

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
		log.Println("Extension Not Found:", path)
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
