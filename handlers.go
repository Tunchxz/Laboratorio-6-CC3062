package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func getSeries(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, status, last_episode_watched, total_episodes, ranking FROM series")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var series []Serie
		for rows.Next() {
			var s Serie
			if err := rows.Scan(&s.ID, &s.Title, &s.Status, &s.LastEpisodeWatched, &s.TotalEpisodes, &s.Ranking); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			series = append(series, s)
		}
		json.NewEncoder(w).Encode(series)
	}
}

func getSerie(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var s Serie
		err := db.QueryRow("SELECT id, title, status, last_episode_watched, total_episodes, ranking FROM series WHERE id=$1", id).
			Scan(&s.ID, &s.Title, &s.Status, &s.LastEpisodeWatched, &s.TotalEpisodes, &s.Ranking)
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		json.NewEncoder(w).Encode(s)
	}
}

func createSerie(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var s Serie
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := db.QueryRow("INSERT INTO series (title, status, last_episode_watched, total_episodes, ranking) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			s.Title, s.Status, s.LastEpisodeWatched, s.TotalEpisodes, s.Ranking).Scan(&s.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(s)
	}
}

func updateSerie(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var s Serie
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := db.Exec("UPDATE series SET title=$1, status=$2, last_episode_watched=$3, total_episodes=$4, ranking=$5 WHERE id=$6",
			s.Title, s.Status, s.LastEpisodeWatched, s.TotalEpisodes, s.Ranking, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = db.QueryRow("SELECT id, title, status, last_episode_watched, total_episodes, ranking FROM series WHERE id=$1", id).
			Scan(&s.ID, &s.Title, &s.Status, &s.LastEpisodeWatched, &s.TotalEpisodes, &s.Ranking)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(s)
	}
}

func deleteSerie(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("DELETE FROM series WHERE id=$1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func updateSerieStatus(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var payload struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := db.Exec("UPDATE series SET status=$1 WHERE id=$2", payload.Status, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var updatedSerie Serie
		err = db.QueryRow("SELECT id, title, status, last_episode_watched, total_episodes, ranking FROM series WHERE id=$1", id).
			Scan(&updatedSerie.ID, &updatedSerie.Title, &updatedSerie.Status, &updatedSerie.LastEpisodeWatched, &updatedSerie.TotalEpisodes, &updatedSerie.Ranking)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(updatedSerie)
	}
}

func updateSerieLastEpisodeWatched(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var s Serie
		err := db.QueryRow("SELECT last_episode_watched, total_episodes FROM series WHERE id=$1", id).
			Scan(&s.LastEpisodeWatched, &s.TotalEpisodes)
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if s.LastEpisodeWatched < s.TotalEpisodes {
			s.LastEpisodeWatched++
			_, err := db.Exec("UPDATE series SET last_episode_watched=$1 WHERE id=$2", s.LastEpisodeWatched, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		err = db.QueryRow("SELECT id, title, status, last_episode_watched, total_episodes, ranking FROM series WHERE id=$1", id).
			Scan(&s.ID, &s.Title, &s.Status, &s.LastEpisodeWatched, &s.TotalEpisodes, &s.Ranking)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(s)
	}
}

func createSeriesTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS series (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		status TEXT NOT NULL,
		last_episode_watched INT,
		total_episodes INT,
		ranking INT
	)`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}
