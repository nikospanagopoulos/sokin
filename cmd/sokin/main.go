package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	name := os.Getenv("POSTGRES_DB")
	if user == "" || pass == "" || name == "" {
		log.Fatal("POSTGRES_* variables are not set (run: set -a; source .env; set +a)")
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pass),
		Host:   "127.0.0.1:5432",
		Path:   name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, u.String())
	if err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	defer conn.Close(context.Background())

	var n int
	if err := conn.QueryRow(ctx, "SELECT count(*) FROM entries").Scan(&n); err != nil {
		log.Fatalf("query failed: %v", err)
	}
	fmt.Println("connected, entries:", n)
}
