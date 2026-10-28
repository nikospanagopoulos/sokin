package web

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/nikospanagopoulos/sokin/internal/store"
)

type Server struct {
	store *store.Store
	tmpl  *template.Template
}

// New loads the HTML templates. Run the program from the repository root.
func New(s *store.Store) (*Server, error) {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Server{store: s, tmpl: tmpl}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.index)
	mux.HandleFunc("POST /entries", s.saveEntry)
	return mux
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	entries, err := s.store.ListEntries(r.Context())
	if err != nil {
		log.Printf("list entries: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Today   string
		Entries any
	}{
		Today:   time.Now().Format("2006-01-02"),
		Entries: entries,
	}
	if err := s.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("render: %v", err)
	}
}

func (s *Server) saveEntry(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("entry_date"))
	if err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	err = s.store.SaveEntry(r.Context(), date,
		r.FormValue("did_today"), r.FormValue("learned"), r.FormValue("mistakes"))
	if err != nil {
		log.Printf("save entry: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
