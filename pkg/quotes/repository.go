package quotes

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shaj13/go-guardian/v2/auth"
)

type QuotesRepository interface {
	// GetQuotes() []KindleQuote
	// GetBooks() []string
	// GetAuthors() []string
	// GetQuotesByAuthor(book string) []KindleQuote
	// GetQuotesByTitle(author string) []KindleQuote
}

type DBQuotesRepository interface {
	QuotesRepository
	GetQuotes(ctx context.Context) ([]KindleQuote, error)
	ImportQuotes(ctx context.Context, quotes []KindleQuote) (int, error)
}

type DBRepository struct {
	db *pgxpool.Pool
}

func CreateRepository(database *pgxpool.Pool) DBQuotesRepository {
	repo := DBRepository{db: database}

	return repo
}

func (qr DBRepository) GetQuotes(ctx context.Context) ([]KindleQuote, error) {
	user_id := auth.UserFromCtx(ctx).GetID()
	sqlGetQuotes := `
select tbl_authors.author_name, tbl_sources.source_title,tbl_quotes.quote, tbl_quotes.date_taken 
from tbl_quotes 
INNER JOIN tbl_sources ON tbl_sources.source_id = tbl_quotes.source_id 
INNER JOIN tbl_authors ON tbl_authors.author_id = tbl_sources.author_id 
where tbl_quotes.user_id=$1 
`
	quotes := []KindleQuote{}
	conn, err := qr.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, sqlGetQuotes, user_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		quote := KindleQuote{}
		t := time.Now()
		err = rows.Scan(&quote.Author, &quote.Title, &quote.Quote, &t)
		quote.Date = t.String()
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}
	defer conn.Release()
	return quotes, nil
}

func (qr DBRepository) ImportQuotes(ctx context.Context, quotes []KindleQuote) (int, error) {
	user_id := auth.UserFromCtx(ctx).GetID()

	sqlInsertAuthor := `
INSERT INTO tbl_authors (author_name, user_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
	`
	sqlSelectAuthor := `
SELECT author_id FROM tbl_authors  WHERE author_name = $1 AND user_id = $2
	`
	sqlInsertBook := `
INSERT INTO tbl_sources (source_title, author_id, user_id)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING
`
	sqlSelectBook := `
SELECT source_id FROM tbl_sources  WHERE source_title = $1 AND user_id = $2
	`

	sqlInsertQuote := `
INSERT INTO tbl_quotes (source_id, quote, date_taken, user_id)
VALUES ($1, $2, $3, $4)
	`

	for _, quote := range quotes {
		auth_id := -1
		source_id := -1
		conn, err := qr.db.Acquire(ctx)
		if err != nil {
			return 0, fmt.Errorf("can't acuire connection")
		}
		tr, err := conn.Begin(ctx)
		rr, err := tr.Query(ctx, sqlInsertAuthor, quote.Author, user_id)
		if err != nil {

			return 0, err
		}
		rr.Close()
		tr.Commit(ctx)
		conn.Release()
		// time.Sleep(1 * time.Second)
		conn, err = qr.db.Acquire(ctx)
		if err != nil {
			return 0, err
		}
		err = conn.QueryRow(ctx, sqlSelectAuthor, quote.Author, user_id).Scan(&auth_id)
		if err != nil {
			return 0, err
		}
		conn.Release()

		conn, err = qr.db.Acquire(ctx)
		if err != nil {
			return 0, err
		}
		tr, err = conn.Begin(ctx)
		rr, err = tr.Query(ctx, sqlInsertBook, quote.Title, auth_id, user_id)
		if err != nil {
			return 0, err
		}
		rr.Close()
		tr.Commit(ctx)
		conn.Release()
		conn, err = qr.db.Acquire(ctx)
		if err != nil {
			return 0, err
		}
		err = conn.QueryRow(ctx, sqlSelectBook, quote.Title, user_id).Scan(&source_id)
		if err != nil {
			return 0, err
		}
		conn.Release()
		conn, err = qr.db.Acquire(ctx)
		if err != nil {
			return 0, err
		}
		tr, err = conn.Begin(ctx)
		rr, err = tr.Query(ctx, sqlInsertQuote, source_id, quote.Quote, time.Now(), user_id)
		if err != nil {
			return 0, err
		}
		rr.Close()
		tr.Commit(ctx)
		conn.Release()
	}
	return len(quotes), nil
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
