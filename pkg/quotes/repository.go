package quotes

type QuotesRepository interface {
	GetBooks() []string
	GetAuthors() []string
	GetQuotesByAuthor(book string) []KindleQuote
	GetQuotesByTitle(author string) []KindleQuote
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
	for book, _ := range result {
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
	for author, _ := range result {
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
