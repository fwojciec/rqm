package rqm_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/fwojciec/rqm"
)

type TestProcessor struct {
	Body  *bytes.Buffer
	Calls int
}

func (tp *TestProcessor) Process(r *rqm.Rq) error {
	tp.Body = r.Body
	tp.Calls++
	return nil
}

func TestRequestMaker(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		body       *bytes.Buffer
		statusCode int
		addRequest bool
		expBody    *bytes.Buffer
		expCalls   int
	}{
		{
			name:       "single valid request",
			body:       bytes.NewBufferString(`<a href="https://test.com">test</a>`),
			statusCode: 200,
			addRequest: true,
			expBody:    bytes.NewBufferString(`<a href="https://test.com">test</a>`),
			expCalls:   1,
		},
		{
			name:       "empty queue",
			body:       nil,
			statusCode: 200,
			addRequest: false,
			expBody:    nil,
			expCalls:   0,
		},
		{
			name:       "non-OK response",
			body:       nil,
			statusCode: 404,
			addRequest: true,
			expBody:    nil,
			expCalls:   3,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				fmt.Fprintf(w, tc.body.String())
			}))
			defer ts.Close()
			q := &rqm.Queue{}
			p := &TestProcessor{}
			var rm *rqm.RequestMaker
			if tc.addRequest {
				rm = rqm.NewRequestMaker(p, q, 1, 2, &rqm.Rq{URL: ts.URL})
			} else {
				rm = rqm.NewRequestMaker(p, q, 1, 2)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			rm.Run(ctx)
			if !reflect.DeepEqual(tc.expBody, p.Body) {
				t.Fatalf("expected %s, got %s", tc.expBody, p.Body)
			}
			if p.Calls != tc.expCalls {
				t.Fatalf("expected %d calls, got: %d", tc.expCalls, p.Calls)
			}
		})
	}
}
