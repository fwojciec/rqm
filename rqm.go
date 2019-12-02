package rqm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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
	Pop(context.Context) (*Rq, error)
	Push(context.Context, *Rq) error
	IsEmpty(context.Context) (bool, error)
}

// Processor processes fetched requests.
type Processor interface {
	Process(context.Context, *Rq, Queuer) error
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
			isEmpty, err := rm.IsEmpty(ctx)
			if err != nil {
				log.Printf("database error: %s", err)
				return err
			}
			if isEmpty {
				continue
			}
			rm.makeRequest(ctx)
		}
	}
}

func (rm *RequestMaker) makeRequest(ctx context.Context) error {
	r, err := rm.Pop(ctx)
	if err != nil {
		return err
	}
	b, err := getPageCtx(ctx, r.URL)
	if err != nil {
		if r.retryCounter < 2 {
			// try three times
			r.retryCounter++
			err := rm.Push(ctx, r)
			if err != nil {
				return err
			}
		} else {
			// TODO: come up with a better way to handle errors
			log.Printf("failed to fetch %q: %s", r.URL, err)
		}
	}
	r.Body = b
	if err = rm.Process(ctx, r, rm); err != nil {
		// TODO: come up with a better way to handle errors
		log.Printf("failed to process the body for %q: %s", r.URL, err)
	}
	return nil
}

// NewRequestMaker returns a pointer to a new instance of RequestMaker.
func NewRequestMaker(p Processor, q Queuer, minDelay, maxDelay int, rs ...*Rq) (*RequestMaker, error) {
	rm := &RequestMaker{q, p, minDelay, maxDelay}
	ctx := context.Background()
	for _, r := range rs {
		err := rm.Push(ctx, r)
		if err != nil {
			return nil, err
		}
	}
	return rm, nil
}

func getPageCtx(ctx context.Context, url string) (*bytes.Buffer, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
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
		return nil, fmt.Errorf("status code error: %s", res.Status)
	}
	var b bytes.Buffer
	_, err = io.Copy(&b, res.Body)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
