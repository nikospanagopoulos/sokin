package web

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/nikospanagopoulos/sokin/internal/store"
)

const dateLayout = "2006-01-02"

type Server struct {
	store *store.Store
	pages map[string]*template.Template
}

// Page holds what the shared layout needs. Every page data struct embeds it.
type Page struct {
	Active string
}

type dashboardData struct {
	Page
	TodayEntry   *store.Entry
	Entries      []store.Entry
	DaysWeek     int
	WeekElapsed  int
	WeekPct      int
	DaysMonth    int
	MonthElapsed int
	MonthPct     int
}

type writeData struct {
	Page
	Date  string
	Entry *store.Entry
}

// New loads the templates. Each page is parsed together with the layout.
// Run the program from the repository root.
func New(s *store.Store) (*Server, error) {
	pages := make(map[string]*template.Template)
	for _, name := range []string{"dashboard.html", "write.html"} {
		t, err := template.ParseFiles("templates/layout.html", "templates/"+name)
		if err != nil {
			return nil, err
		}
		pages[name] = t
	}
	return &Server{store: s, pages: pages}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.dashboard)
	mux.HandleFunc("GET /write", s.writeForm)
	mux.HandleFunc("POST /entries", s.saveEntry)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	return mux
}

// render writes into a buffer first, so a template error never sends half a page.
func (s *Server) render(w http.ResponseWriter, page string, data any) {
	var buf bytes.Buffer
	if err := s.pages[page].ExecuteTemplate(&buf, "layout", data); err != nil {
		serverError(w, "render "+page, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

func serverError(w http.ResponseWriter, what string, err error) {
	log.Printf("%s: %v", what, err)
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func today() time.Time {
	n := time.Now()
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
}

func percent(part, whole int) int {
	if whole == 0 {
		return 0
	}
	return part * 100 / whole
}

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	t := today()

	weekElapsed := (int(t.Weekday())+6)%7 + 1 // Monday = 1 ... Sunday = 7
	weekStart := t.AddDate(0, 0, -(weekElapsed - 1))
	monthStart := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)

	daysWeek, err := s.store.CountEntries(ctx, weekStart, t)
	if err != nil {
		serverError(w, "count week", err)
		return
	}
	daysMonth, err := s.store.CountEntries(ctx, monthStart, t)
	if err != nil {
		serverError(w, "count month", err)
		return
	}
	todayEntry, err := s.store.GetEntry(ctx, t)
	if err != nil {
		serverError(w, "get today", err)
		return
	}
	entries, err := s.store.ListEntries(ctx)
	if err != nil {
		serverError(w, "list entries", err)
		return
	}

	s.render(w, "dashboard.html", dashboardData{
		Page:         Page{Active: "dashboard"},
		TodayEntry:   todayEntry,
		Entries:      entries,
		DaysWeek:     daysWeek,
		WeekElapsed:  weekElapsed,
		WeekPct:      percent(daysWeek, weekElapsed),
		DaysMonth:    daysMonth,
		MonthElapsed: t.Day(),
		MonthPct:     percent(daysMonth, t.Day()),
	})
}

func (s *Server) writeForm(w http.ResponseWriter, r *http.Request) {
	date := today()
	if q := r.URL.Query().Get("date"); q != "" {
		d, err := time.Parse(dateLayout, q)
		if err != nil {
			http.Error(w, "invalid date", http.StatusBadRequest)
			return
		}
		date = d
	}

	entry, err := s.store.GetEntry(r.Context(), date)
	if err != nil {
		serverError(w, "get entry", err)
		return
	}

	s.render(w, "write.html", writeData{
		Page:  Page{Active: "write"},
		Date:  date.Format(dateLayout),
		Entry: entry,
	})
}

func (s *Server) saveEntry(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(dateLayout, r.FormValue("entry_date"))
	if err != nil {
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	err = s.store.SaveEntry(r.Context(), date,
		r.FormValue("did_today"), r.FormValue("learned"), r.FormValue("mistakes"))
	if err != nil {
		serverError(w, "save entry", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
