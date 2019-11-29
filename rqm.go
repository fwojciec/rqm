package rqm

import "github.com/PuerkitoBio/goquery"

// PageType is a type of the requested page. The type determines how the page
// is eventually processed.
type PageType int

// Rq represents a request to be processed.
type Rq struct {
	// URL is the request url.
	URL string
	// Retry number
	Retry int
	// PageType is a request type which will determine how it will be processed.
	PageType PageType
	// Doc is a goquery Document
	Doc *goquery.Document
}
