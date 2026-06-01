package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed *.html robots.txt
var staticFiles embed.FS

type rsvpPayload struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	AdditionalGuests int    `json:"additional_guests"`
	Attending        string `json:"attending"`
	DietaryNotes     string `json:"dietary_notes"`
	SubmittedAt      string `json:"submitted_at"`
}

var db *sql.DB

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	dbPath := flag.String("db", "./rsvps.db", "path to sqlite database file")
	flag.Parse()

	absDB, err := filepath.Abs(*dbPath)
	if err != nil {
		log.Fatalf("resolve db path: %v", err)
	}
	log.Printf("using database: %s", absDB)

	db, err = sql.Open("sqlite", absDB)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("connect db: %v", err)
	}

	if _, err = db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		log.Fatalf("pragma: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS rsvps (
			id                INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name        TEXT    NOT NULL,
			last_name         TEXT,
			email             TEXT    NOT NULL,
			additional_guests INTEGER NOT NULL DEFAULT 0,
			attending         TEXT    NOT NULL,
			dietary_notes     TEXT,
			submitted_at      TEXT,
			created_at        DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("create table: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(staticFiles)))
	mux.HandleFunc("POST /api/rsvp", handleRSVP)

	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func handleRSVP(w http.ResponseWriter, r *http.Request) {
	var p rsvpPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if p.FirstName == "" || p.Email == "" {
		http.Error(w, "first_name and email are required", http.StatusUnprocessableEntity)
		return
	}

	if p.SubmittedAt == "" {
		p.SubmittedAt = time.Now().UTC().Format(time.RFC3339)
	}

	_, err := db.Exec(`
		INSERT INTO rsvps (first_name, last_name, email, additional_guests, attending, dietary_notes, submitted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, p.FirstName, p.LastName, p.Email, p.AdditionalGuests, p.Attending, p.DietaryNotes, p.SubmittedAt)
	if err != nil {
		log.Printf("db insert: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
