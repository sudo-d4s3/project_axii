package main

import (
	//"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key            = []byte("not-production-key______________")
	store          = sessions.NewCookieStore(key)
	request_id_key = 8080
	log_file       = "server_log.log"
)

func loggingHandler(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(request_id_key).(string)
				if !ok {
					requestID = "UNKNOWN"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())

			}()

			next.ServeHTTP(w, r)
		})
	}
}

func tracingHandler(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), request_id_key, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func level0(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
	tmpl.Execute(w, nil)
}

func c341b271f5dba18dd4099435670a2c74(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "level")

	session.Values["level0"] = true
	session.Save(r, w)

	http.Redirect(w, r, "/level1", http.StatusMovedPermanently)
}
func level1(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "level")

	if auth, ok := session.Values["level0"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		fmt.Fprintf(w, session.Values["True"].(string))
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/menu.html"))
	tmpl.Execute(w, nil)
}

func main() {
	//Logging
	logfile, err := os.OpenFile(log_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	logger := log.New(logfile, "http: ", log.LstdFlags)
	defer logfile.Close()

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	//Router
	r := mux.NewRouter()

	r.HandleFunc("/", level0)
	r.HandleFunc("/level1", level1)
	r.HandleFunc("/c341b271f5dba18dd4099435670a2c74", c341b271f5dba18dd4099435670a2c74)

	fs := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))

	server := &http.Server{
		Addr:         ":8888",
		Handler:      tracingHandler(nextRequestID)(loggingHandler(logger)(r)),
		ErrorLog:     logger,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())

}
