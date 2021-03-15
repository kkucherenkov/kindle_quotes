package quotes

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type QuotesRepository interface {
	GetBooks() []string
	GetAuthors() []string
	GetQuotesByAuthor(book string) []KindleQuote
	GetQuotesByTitle(author string) []KindleQuote
}

type DBQuotesRepository interface {
	ImportQuotes(quotes []KindleQuote)
}

type DBRepository struct {
	db *pgxpool.Pool
}

func CreateRepository(database *pgxpool.Pool) DBQuotesRepository {
	repo := DBRepository{db: database}

	return repo
}

func (qr DBRepository) ImportQuotes(quotes []KindleQuote) {

	sqlInsertAuthor := `
INSERT INTO tbl_authors (author_name)
VALUES ($1)
ON CONFLICT DO NOTHING
	`
	sqlSelectAuthor := `
SELECT author_id FROM tbl_authors  WHERE author_name = $1
	`
	sqlInsertBook := `
INSERT INTO tbl_sources (source_title, author_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`
	sqlSelectBook := `
SELECT source_id FROM tbl_sources  WHERE source_title = $1
	`

	sqlInsertQuote := `
INSERT INTO tbl_quotes (source_id, quote, date_taken)
VALUES ($1, $2, $3)
	`

	for i, quote := range quotes {
		auth_id := -1
		source_id := -1
		conn, err := qr.db.Acquire(context.Background())
		if err != nil {
			fmt.Println("can't acuire connection")
			return
		}
		tr, err := conn.Begin(context.Background())
		rr, err := tr.Query(context.Background(), sqlInsertAuthor, quote.Author)
		if err != nil {
			fmt.Println("transaction error", err)
			return
		}
		rr.Close()
		tr.Commit(context.Background())
		conn.Release()
		// time.Sleep(1 * time.Second)
		conn, err = qr.db.Acquire(context.Background())
		if err != nil {
			fmt.Println("can't acuire connection")
			return
		}
		err = conn.QueryRow(context.Background(), sqlSelectAuthor, quote.Author).Scan(&auth_id)
		if err != nil {
			fmt.Println("Error in authors")
			return
		}
		conn.Release()

		conn, err = qr.db.Acquire(context.Background())
		if err != nil {
			fmt.Println("can't acuire connection")
			return
		}
		tr, err = conn.Begin(context.Background())
		rr, err = tr.Query(context.Background(), sqlInsertBook, quote.Title, auth_id)
		if err != nil {
			fmt.Println("transaction error", err)
			return
		}
		rr.Close()
		tr.Commit(context.Background())
		conn.Release()
		// time.Sleep(1 * time.Second)
		conn, err = qr.db.Acquire(context.Background())
		if err != nil {
			fmt.Println("can't acuire connection")
			return
		}
		err = conn.QueryRow(context.Background(), sqlSelectBook, quote.Title).Scan(&source_id)
		if err != nil {
			fmt.Println("error in sources")
			return
		}
		conn.Release()
		conn, err = qr.db.Acquire(context.Background())
		if err != nil {
			fmt.Println("can't acuire connection")
			return
		}
		tr, err = conn.Begin(context.Background())
		rr, err = tr.Query(context.Background(), sqlInsertQuote, source_id, quote.Quote, time.Now())
		if err != nil {
			fmt.Println("transaction error", err)
			return
		}
		rr.Close()
		tr.Commit(context.Background())
		conn.Release()
		// time.Sleep(1 * time.Second)

		fmt.Println(fmt.Sprint(i+1) + " quote inserted")
	}

}

type InmemRepo struct {
	QuotesRepository
	qs []KindleQuote
}

func New(qs []KindleQuote) QuotesRepository {
	return InmemRepo{qs: qs}
}

func (r InmemRepo) GetBooks() []string {
	var result map[string]bool
	var books []string
	for _, quote := range r.qs {
		result[quote.Title] = true
	}
	for book := range result {
		books = append(books, book)
	}
	return books
}

func (r InmemRepo) GetAuthors() []string {
	var result map[string]bool
	var authors []string
	for _, quote := range r.qs {
		result[quote.Author] = true
	}
	for author := range result {
		authors = append(authors, author)
	}
	return authors
}

func (r InmemRepo) GetQuotesByTitle(book string) []KindleQuote {
	var result []KindleQuote
	for _, quote := range r.qs {
		if quote.Title == book {
			result = append(result, quote)
		}
	}
	return result
}

func (r InmemRepo) GetQuotesByAuthor(author string) []KindleQuote {
	var result []KindleQuote
	for _, quote := range r.qs {
		if quote.Title == author {
			result = append(result, quote)
		}
	}
	return result
}
