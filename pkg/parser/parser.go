package parser

import (
	"bufio"
	"io"

	"github.com/kkucherenkov/kindle_quotes/pkg/quotes"
)

func ParseQuotes(file io.Reader) []quotes.KindleQuote {
	scanner := bufio.NewScanner(file)
	var qs []quotes.KindleQuote
	var line string
	for scanner.Scan() {
		quote := quotes.KindleQuote{}
		line = scanner.Text()
		quote.ParseTitleAndAuthor(line)
		scanner.Scan()
		line = scanner.Text()
		quote.ParseDate(line)
		scanner.Scan()
		line = scanner.Text()
		scanner.Scan()
		line = scanner.Text()
		quote.Quote = line
		scanner.Scan()
		line = scanner.Text()
		qs = append(qs, quote)
	}
	return qs
}
