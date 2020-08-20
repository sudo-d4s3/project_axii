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

	fs := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      tracingHandler(nextRequestID)(loggingHandler(logger)(r)),
		ErrorLog:     logger,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())

}
