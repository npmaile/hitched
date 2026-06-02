package main

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"embed"
	"encoding/hex"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed *.html *.css robots.txt *.png
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

type rsvpRow struct {
	FirstName        string
	LastName         string
	Email            string
	AdditionalGuests int
	Attending        string
	DietaryNotes     string
	SubmittedAt      string
}

var (
	db        *sql.DB
	adminUser string
	adminPass string
	sessions  sync.Map // token(string) → expiry(time.Time)
)

func newToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func isValidSession(token string) bool {
	v, ok := sessions.Load(token)
	if !ok {
		return false
	}
	if time.Now().After(v.(time.Time)) {
		sessions.Delete(token)
		return false
	}
	return true
}

func isAuthenticated(r *http.Request) bool {
	c, err := r.Cookie("admin_session")
	return err == nil && isValidSession(c.Value)
}

func main() {
	addr   := flag.String("addr", ":8080", "listen address")
	dbPath := flag.String("db", "./rsvps.db", "path to sqlite database file")
	aUser  := flag.String("admin-user", "admin", "admin username")
	aPass  := flag.String("admin-pass", "", "admin password (admin route disabled if unset)")
	flag.Parse()

	adminUser = *aUser
	adminPass = *aPass
	if adminPass == "" {
		log.Println("warning: -admin-pass not set; /admin will return 404")
	}

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
	mux.HandleFunc("GET /admin", handleAdminPage)
	mux.HandleFunc("POST /admin/login", handleAdminLogin)
	mux.HandleFunc("POST /admin/logout", handleAdminLogout)

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

// ── Admin ─────────────────────────────────────────────────────────────────────

var loginTmpl = template.Must(template.New("login").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Admin — Nate &amp; Susanna</title>
<link href="https://fonts.googleapis.com/css2?family=Fraunces:ital,wght@1,700&family=Space+Grotesk:wght@300;400;500&display=swap" rel="stylesheet">
<link rel="stylesheet" href="/shared.css">
<style>
.login-wrap{min-height:100vh;display:flex;align-items:center;justify-content:center;background:var(--ink)}
.login-box{background:var(--white);border:2.5px solid var(--ink);padding:48px 40px;width:100%;max-width:380px}
.login-title{font-family:'Fraunces',serif;font-size:32px;font-weight:700;font-style:italic;margin-bottom:32px;color:var(--ink)}
.field{display:flex;flex-direction:column;gap:7px;margin-bottom:16px}
.field label{font-size:10px;letter-spacing:0.25em;text-transform:uppercase;font-weight:500;color:rgba(26,16,8,0.55)}
.field input{background:var(--white);border:2px solid var(--ink);font-family:'Space Grotesk',sans-serif;font-size:14px;padding:12px 14px;color:var(--ink);outline:none;width:100%;-webkit-appearance:none;border-radius:0}
.field input:focus{background:var(--lime)}
.login-btn{width:100%;padding:14px;background:var(--coral);border:2px solid var(--ink);color:var(--white);font-family:'Space Grotesk',sans-serif;font-size:11px;font-weight:500;letter-spacing:0.28em;text-transform:uppercase;cursor:pointer;margin-top:8px;transition:background 0.15s}
.login-btn:hover{background:var(--ink)}
.error{color:var(--coral);font-size:12px;margin-bottom:16px;letter-spacing:0.05em}
</style>
</head>
<body>
<div class="login-wrap">
  <div class="login-box">
    <div class="login-title">Admin</div>
    {{if .Error}}<p class="error">{{.Error}}</p>{{end}}
    <form method="POST" action="/admin/login">
      <div class="field">
        <label for="u">Username</label>
        <input type="text" id="u" name="username" autocomplete="username" required>
      </div>
      <div class="field">
        <label for="p">Password</label>
        <input type="password" id="p" name="password" autocomplete="current-password" required>
      </div>
      <button class="login-btn" type="submit">Sign in →</button>
    </form>
  </div>
</div>
</body>
</html>`))

type adminData struct {
	RSVPs       []rsvpRow
	TotalYes    int
	TotalNo     int
	TotalMaybe  int
	TotalGuests int
}

var dashTmpl = template.Must(template.New("dash").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Admin — Nate &amp; Susanna</title>
<link href="https://fonts.googleapis.com/css2?family=Fraunces:ital,wght@1,700&family=Space+Grotesk:wght@300;400;500&display=swap" rel="stylesheet">
<link rel="stylesheet" href="/shared.css">
<style>
.admin-bar{background:var(--ink);padding:16px 32px;display:flex;align-items:center;justify-content:space-between}
.admin-bar-title{font-family:'Fraunces',serif;font-size:22px;font-style:italic;font-weight:700;color:var(--lime)}
.logout-form button{background:none;border:1px solid rgba(255,254,249,0.25);color:rgba(255,254,249,0.6);font-family:'Space Grotesk',sans-serif;font-size:10px;letter-spacing:0.2em;text-transform:uppercase;padding:6px 14px;cursor:pointer;transition:all 0.15s}
.logout-form button:hover{border-color:var(--lime);color:var(--lime)}
.stats{display:grid;grid-template-columns:repeat(4,1fr);border-bottom:2.5px solid var(--ink)}
.stat{padding:28px 24px;border-right:2.5px solid var(--ink)}
.stat:last-child{border-right:none}
.stat:nth-child(1){background:var(--lime)}
.stat:nth-child(2){background:var(--coral)}
.stat:nth-child(3){background:var(--violet)}
.stat:nth-child(4){background:var(--sky)}
.stat-label{font-size:10px;letter-spacing:0.22em;text-transform:uppercase;font-weight:500;color:rgba(26,16,8,0.55);margin-bottom:8px}
.stat-num{font-family:'Fraunces',serif;font-size:48px;font-weight:700;font-style:italic;color:var(--ink);line-height:1}
.table-wrap{padding:40px 32px;overflow-x:auto}
table{width:100%;border-collapse:collapse;font-size:13px;min-width:600px}
th{font-size:10px;letter-spacing:0.2em;text-transform:uppercase;font-weight:500;color:rgba(26,16,8,0.5);text-align:left;padding:10px 16px;border-bottom:2.5px solid var(--ink)}
td{padding:12px 16px;border-bottom:1px solid rgba(26,16,8,0.1);vertical-align:top}
tr:hover td{background:var(--yellow)}
.badge{display:inline-block;font-size:9px;letter-spacing:0.15em;text-transform:uppercase;font-weight:500;padding:3px 8px;border:1.5px solid var(--ink)}
.badge-yes{background:var(--lime)}
.badge-no{background:var(--coral);color:var(--white)}
.badge-maybe{background:var(--violet)}
.empty{text-align:center;padding:48px;color:rgba(26,16,8,0.35);font-weight:300;letter-spacing:0.05em}
@media(max-width:600px){.stats{grid-template-columns:1fr 1fr}.stat:nth-child(2){border-right:none}.stat:nth-child(3){border-top:2.5px solid var(--ink)}.stat:nth-child(4){border-top:2.5px solid var(--ink);border-right:none}.table-wrap{padding:24px 16px}}
</style>
</head>
<body>
<div class="admin-bar">
  <span class="admin-bar-title">RSVPs</span>
  <form class="logout-form" method="POST" action="/admin/logout">
    <button type="submit">Sign out</button>
  </form>
</div>
<div class="stats">
  <div class="stat"><div class="stat-label">Attending</div><div class="stat-num">{{.TotalYes}}</div></div>
  <div class="stat"><div class="stat-label">Declined</div><div class="stat-num">{{.TotalNo}}</div></div>
  <div class="stat"><div class="stat-label">Maybe</div><div class="stat-num">{{.TotalMaybe}}</div></div>
  <div class="stat"><div class="stat-label">Total Guests</div><div class="stat-num">{{.TotalGuests}}</div></div>
</div>
<div class="table-wrap">
  <table>
    <thead>
      <tr>
        <th>Name</th>
        <th>Email</th>
        <th>Attending</th>
        <th>+Guests</th>
        <th>Notes</th>
        <th>Submitted</th>
      </tr>
    </thead>
    <tbody>
      {{range .RSVPs}}
      <tr>
        <td>{{.FirstName}} {{.LastName}}</td>
        <td>{{.Email}}</td>
        <td><span class="badge badge-{{.Attending}}">{{.Attending}}</span></td>
        <td>{{.AdditionalGuests}}</td>
        <td>{{.DietaryNotes}}</td>
        <td>{{.SubmittedAt}}</td>
      </tr>
      {{else}}
      <tr><td class="empty" colspan="6">No RSVPs yet.</td></tr>
      {{end}}
    </tbody>
  </table>
</div>
</body>
</html>`))

func handleAdminPage(w http.ResponseWriter, r *http.Request) {
	if adminPass == "" {
		http.NotFound(w, r)
		return
	}
	if !isAuthenticated(r) {
		loginTmpl.Execute(w, struct{ Error string }{})
		return
	}

	rows, err := db.Query(`
		SELECT first_name, last_name, email, additional_guests, attending,
		       COALESCE(dietary_notes,''), COALESCE(submitted_at,'')
		FROM rsvps ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data adminData
	for rows.Next() {
		var row rsvpRow
		if err := rows.Scan(&row.FirstName, &row.LastName, &row.Email,
			&row.AdditionalGuests, &row.Attending, &row.DietaryNotes, &row.SubmittedAt); err != nil {
			continue
		}
		switch row.Attending {
		case "yes":
			data.TotalYes++
			data.TotalGuests += 1 + row.AdditionalGuests
		case "no":
			data.TotalNo++
		default:
			data.TotalMaybe++
		}
		data.RSVPs = append(data.RSVPs, row)
	}

	dashTmpl.Execute(w, data)
}

func handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	if adminPass == "" {
		http.NotFound(w, r)
		return
	}

	userOK := subtle.ConstantTimeCompare([]byte(r.FormValue("username")), []byte(adminUser)) == 1
	passOK := subtle.ConstantTimeCompare([]byte(r.FormValue("password")), []byte(adminPass)) == 1
	if !userOK || !passOK {
		w.WriteHeader(http.StatusUnauthorized)
		loginTmpl.Execute(w, struct{ Error string }{"Invalid credentials."})
		return
	}

	token := newToken()
	sessions.Store(token, time.Now().Add(24*time.Hour))
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    token,
		Path:     "/admin",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func handleAdminLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("admin_session"); err == nil {
		sessions.Delete(c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "admin_session",
		Value:  "",
		Path:   "/admin",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
