package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
		tmpl.Execute(w, nil)
	})

	fs := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", r)
}
