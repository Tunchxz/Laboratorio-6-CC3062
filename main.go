package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}
	defer db.Close()

	createSeriesTable(db)

	router := mux.NewRouter()

	router.HandleFunc("/api/series", getSeries(db)).Methods("GET")
	router.HandleFunc("/api/series/{id}", getSerie(db)).Methods("GET")
	router.HandleFunc("/api/series/{id}", updateSerie(db)).Methods("PUT")
	router.HandleFunc("/api/series", createSerie(db)).Methods("POST")
	router.HandleFunc("/api/series/{id}", deleteSerie(db)).Methods("DELETE")
	router.HandleFunc("/api/series/{id}/status", updateSerieStatus(db)).Methods("PATCH")
	router.HandleFunc("/api/series/{id}/episode", updateSerieLastEpisodeWatched(db)).Methods("PATCH")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	wrappedRouter := jsonContentTypeMiddleware(corsHandler(router))

	log.Fatal(http.ListenAndServe(":8080", wrappedRouter))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
