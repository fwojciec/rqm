package rqm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

// PageType identifies the variant of the requested page. The type determines
// how the page will be eventually processed.
type PageType int

// Rq represents a request to be processed.
type Rq struct {
	URL          string
	PageType     PageType
	Body         *bytes.Buffer
	retryCounter int
}

// Queuer manages the queue of requests waiting to be processed.
type Queuer interface {
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
	Queuer
	Processor
	minDelay int
	maxDelay int
}

// Run runs the RequestMaker.
func (rm *RequestMaker) Run(ctx context.Context) error {
	ticker := NewRandomTicker(rm.minDelay, rm.maxDelay)
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
	b, err := getPageCtx(ctx, r.URL)
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
	r.Body = b
	if err = rm.Process(r); err != nil {
		// TODO: come up with a better way to handle errors
		log.Printf("failed to process the body for %q: %s", r.URL, err)
	}
}

// NewRequestMaker returns a pointer to a new instance of RequestMaker.
func NewRequestMaker(r *Rq, p Processor, q Queuer, minDelay, maxDelay int) *RequestMaker {
	rm := &RequestMaker{q, p, minDelay, maxDelay}
	rm.Push(r)
	return rm
}

func getPageCtx(ctx context.Context, url string) (*bytes.Buffer, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error %d %s", res.StatusCode, res.Status)
	}
	var b bytes.Buffer
	_, err = io.Copy(&b, res.Body)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
