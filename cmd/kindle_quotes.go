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
	_ "github.com/shaj13/libcache/fifo"
)

func main() {
	conn, err := pgxpool.Connect(context.Background(), "postgresql://postgres:docker@localhost:5432/kindle_quotes")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	//qts := parser.ParseQuotes("../data/My Clippings.txt")
	repo := quotes.CreateRepository(conn)
	//repo.ImportQuotes(qts)
	fmt.Println(repo)

	defer conn.Close()

	router := mux.NewRouter()
	// router.HandleFunc("/v1/auth/token", middleware(http.HandlerFunc(createToken))).Methods("GET")
	// router.HandleFunc("/v1/quotes", middleware(http.HandlerFunc(getQuotes))).Methods("GET")
	log.Println("server started and listening on http://127.0.0.1:8080")
	http.ListenAndServe("127.0.0.1:8080", router)

}

func getQuotes(w http.ResponseWriter, r *http.Request) {

	body := fmt.Sprintf("Author: %s \n", "")
	w.Write([]byte(body))
}
