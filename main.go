package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kkucherenkov/kindle_quotes/pkg/quotes"
	"github.com/kkucherenkov/kindle_quotes/pkg/transport"
	"github.com/kkucherenkov/kindle_quotes/pkg/users"
	_ "github.com/shaj13/libcache/fifo"
)

func main() {
	// conn, err := pgxpool.Connect(context.Background(), "postgresql://postgres:docker@localhost:5432/kindle_quotes")
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }

	//qts := parser.ParseQuotes("../data/My Clippings.txt")
	qRepo := quotes.CreateRepository(nil)
	uRepo := users.NewUserRepository(nil)
	//repo.ImportQuotes(qts)
	fmt.Println(qRepo)

	hanlder := transport.New(qRepo, uRepo)

	// defer conn.Close()

	router := mux.NewRouter()
	router.HandleFunc("/v1/auth/login", hanlder.Login()).Methods("GET")
	// router.HandleFunc("/v1/quotes", middleware(http.HandlerFunc(getQuotes))).Methods("GET")
	log.Println("server started and listening on http://127.0.0.1:8080")
	http.ListenAndServe("127.0.0.1:8080", router)

}
