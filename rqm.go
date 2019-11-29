package rqm

import (
	"context"
	"io"
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/fwojciec/rqm/http"
)

// PageType identifies the variant of the requested page. The type determines
// how the page will be eventually processed.
type PageType int

// Rq represents a request to be processed.
type Rq struct {
	URL          string
	PageType     PageType
	Doc          *goquery.Document
	retryCounter int
}

// AddDoc reads document from io.Reader (e.g. response body) and adds it to
// the Rq struct.
func (r *Rq) AddDoc(b io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(b)
	if err != nil {
		return err
	}
	r.Doc = doc
	return nil
}

// Queue manages the queue of requests waiting to be processed.
type Queue interface {
	Pop() *Rq
	Push(*Rq)
	IsEmpty() bool
}

// Processor processes fetched requests.
type Processor interface {
	Process(*Rq) error
}

// RequestMaker makes requests.
type RequestMaker struct {
	Queue
	Processor
}

// Run runs the RequestMaker.
func (rm *RequestMaker) Run(ctx context.Context) error {
	ticker := NewRandomTicker(2000, 3000)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return ctx.Err()
		case <-ticker.C:
			if rm.IsEmpty() {
				continue
			}
			rm.makeRequest(ctx)
		}
	}
}

func (rm *RequestMaker) makeRequest(ctx context.Context) {
	r := rm.Pop()
	b, err := http.GetPageCtx(ctx, r.URL)
	if err != nil {
		if r.retryCounter < 3 {
			// retry three times
			r.retryCounter++
			rm.Push(r)
		} else {
			// TODO: come up with a better way to handle errors
			log.Printf("failed to fetch %q: %s", r.URL, err)
		}
	}
	err = r.AddDoc(b)
	if err != nil {
		// TODO: come up with a better way to handle errors
		log.Printf("failed to generate doc for %q: %s", r.URL, err)
	}
	err = rm.Process(r)
	if err != nil {
		// TODO: come up with a better way to handle errors
		log.Printf("failed to process the body for %q: %s", r.URL, err)
	}
}

// NewRequestMaker returns a pointer to a new instance of RequestMaker.
func NewRequestMaker(r *Rq, p Processor, q Queue) *RequestMaker {
	rm := &RequestMaker{q, p}
	rm.Push(r)
	return rm
}
