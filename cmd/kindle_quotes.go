package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kkucherenkov/kindle_quotes/pkg/parser"
	"github.com/kkucherenkov/kindle_quotes/pkg/quotes"
)

func main() {
	conn, err := pgxpool.Connect(context.Background(), "postgresql://postgres:docker@localhost:5432/kindle_quotes")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	qts := parser.ParseQuotes("../data/My Clippings.txt")
	repo := quotes.CreateRepository(conn)
	repo.ImportQuotes(qts)
	fmt.Println("done")

	conn.Close()
}
