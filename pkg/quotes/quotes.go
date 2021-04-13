package quotes

import (
	"fmt"
	"strings"
)

type KindleQuote struct {
	Title  string `db:"source_title"`
	Author string `db:"author_name"`
	Date   string `db:"date_taken"`
	Quote  string `db:"quote"`
}

func (quote *KindleQuote) ParseTitleAndAuthor(titleLine string) {
	startAuthorPos := strings.LastIndex(titleLine, "(")
	endAuthorPos := strings.LastIndex(titleLine, ")")
	quote.Title = titleLine[0:startAuthorPos]
	quote.Author = titleLine[startAuthorPos+1 : endAuthorPos]
}

func (quote *KindleQuote) ParseDate(dateLine string) {
	divPosition := strings.LastIndex(dateLine, "|")
	quote.Date = dateLine[divPosition+22:]
}

func (quote KindleQuote) String() {
	fmt.Println("Title: " + quote.Title)
	fmt.Println("Author: " + quote.Author)
	fmt.Println("Quote: " + quote.Quote)
	fmt.Println()
	fmt.Println("Date: " + quote.Date)
	fmt.Println("==========================================")
}
