package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

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

	//qts := parser.ParseQuotes("../data/My Clippings.txt")
	qRepo := quotes.CreateRepository(conn)
	uRepo := users.NewUserRepository(conn)
	//repo.ImportQuotes(qts)
	// fmt.Println(qRepo)

	hanlder := transport.New(qRepo, uRepo)

	defer conn.Close()

	router := mux.NewRouter()
	router.HandleFunc("/v1/auth/login", hanlder.Login()).Methods("POST")
	router.HandleFunc("/v1/auth/register", hanlder.Registration()).Methods("POST")
	router.HandleFunc("/v1/quotes", hanlder.GetQuotes()).Methods("GET")
	log.Println("server started and listening on http://127.0.0.1:8080")
	err = http.ListenAndServe("127.0.0.1:8080", router)
	fmt.Println(err)
}
