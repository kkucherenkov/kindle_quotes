package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kkucherenkov/kindle_quotes/pkg/quotes"
	"github.com/kkucherenkov/kindle_quotes/pkg/transport"
	"github.com/kkucherenkov/kindle_quotes/pkg/users"
	_ "github.com/shaj13/libcache/fifo"
)

func main() {
	conn, err := pgxpool.Connect(context.Background(), "postgresql://postgres:docker@localhost:5432/kindle_quotes")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	qRepo := quotes.CreateRepository(conn)
	uRepo := users.NewUserRepository(conn)

	handler := transport.New(qRepo, uRepo)

	defer conn.Close()

	router := mux.NewRouter()
	router.HandleFunc("/v1/auth/login", handler.Login()).Methods("POST")
	router.HandleFunc("/v1/auth/register", handler.Registration()).Methods("POST")
	router.HandleFunc("/v1/quotes", handler.GetQuotes()).Methods("GET")
	router.HandleFunc("/v1/upload", handler.UploadQuotes()).Methods("POST")
	log.Println("server started and listening on http://127.0.0.1:8080")

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	err = http.ListenAndServe("127.0.0.1:8080", loggedRouter)
	fmt.Println(err)
}
