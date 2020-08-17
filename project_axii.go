package main

import (
	//"bytes"
	//"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("not-production-key______________")
	store = sessions.NewCookieStore(key)
)

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func level0(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
	tmpl.Execute(w, nil)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", level0)

	fs := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", r)
}
