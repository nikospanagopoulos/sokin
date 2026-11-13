package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikospanagopoulos/sokin/internal/store"
	"github.com/nikospanagopoulos/sokin/internal/web"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run holds the real logic so that deferred cleanups (like closing the
// database pool) still run. log.Fatal would skip them.
func run() error {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	st, err := store.Open(dbCtx)
	cancel()
	if err != nil {
		return err
	}
	defer st.Close()

	srv, err := web.New(st)
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	// Cancelled when the process receives Ctrl+C or SIGTERM (e.g. docker stop).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Println("Sokin listening on http://localhost:8080")
		errCh <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return err // the server failed to start (e.g. port already in use)
	case <-ctx.Done():
		log.Println("shutting down...")
	}

	// Give in-flight requests up to 10 seconds to finish.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return err
	}
	log.Println("stopped cleanly")
	return nil
}
