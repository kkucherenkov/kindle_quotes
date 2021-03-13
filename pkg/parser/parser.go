package parser

import (
	"bufio"
	"log"
	"os"

	"github.com/kkucherenkov/kindle_quotes/pkg/quotes"
)

func ParseQuotes(filePath string) []quotes.KindleQuote {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

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
