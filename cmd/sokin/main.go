package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/nikospanagopoulos/sokin/internal/store"
	"github.com/nikospanagopoulos/sokin/internal/web"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	st, err := store.Open(ctx)
	cancel()
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer st.Close()

	srv, err := web.New(st)
	if err != nil {
		log.Fatalf("web: %v", err)
	}

	httpServer := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	log.Println("Sokin listening on http://localhost:8080")
	log.Fatal(httpServer.ListenAndServe())
}
