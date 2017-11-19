package main

import (
	"html/template"
	"io"
	"net/http"
	"time"
)

var tpl *template.Template

// For this code to run, you will need this package:
// go get github.com/satori/go.uuid

func main() {
	tpl = template.Must(template.ParseGlob("./../public/html/*.html"))
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	/*
		t, err := template.ParseFiles("./../public/html/index.html")
		log.Println(err)
		err = t.Execute(w, "")
		//err := tpl.ExecuteTemplate(w, "index.html", "")
		log.Println(err)
		//http.Error(w, "asdasdasd", http.StatusBadRequest)
		//http.ServeContent(w, req, "../public/html/index.html", time.Now(), bytes.NewReader([]byte("allalalal")))
	*/
	f, err := http.Dir("../public/html/").Open("index.html")
	if err == nil {
		content := io.ReadSeeker(f)
		http.ServeContent(w, req, "index.html", time.Now(), content)
		return
	}
}
