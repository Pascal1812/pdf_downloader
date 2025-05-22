package main

import (
	"log"
	"net/http"
	"pdf/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/api/generate", handlers.DownloadPDF()).Methods("POST", "OPTIONS")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("front-end")))

	log.Printf("Server starting at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
