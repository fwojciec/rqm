package rqm_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fwojciec/rqm"
)

type TestProcessor struct {
	Body  *bytes.Buffer
	Calls int
	Error error
}

func (tp *TestProcessor) Process(r *rqm.Rq, q rqm.Queuer) error {
	tp.Body = r.Body
	tp.Calls++
	return tp.Error
}

func TestRequestMaker(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		body           *bytes.Buffer
		statusCode     int
		processorError error
		addRequest     bool
		expBody        *bytes.Buffer
		expCalls       int
		expLog         string
	}{
		{
			name:           "single valid request",
			body:           bytes.NewBufferString(`<a href="https://test.com">test</a>`),
			statusCode:     200,
			processorError: nil,
			addRequest:     true,
			expBody:        bytes.NewBufferString(`<a href="https://test.com">test</a>`),
			expCalls:       1,
			expLog:         "",
		},
		{
			name:           "empty queue",
			body:           nil,
			statusCode:     200,
			processorError: nil,
			addRequest:     false,
			expBody:        nil,
			expCalls:       0,
			expLog:         "",
		},
		{
			name:           "non-OK response",
			body:           nil,
			statusCode:     404,
			processorError: nil,
			addRequest:     true,
			expBody:        nil,
			expCalls:       3,
			expLog:         "status code error: 404 Not Found",
		},
		{
			name:           "processor error",
			body:           nil,
			statusCode:     200,
			processorError: fmt.Errorf("test error"),
			addRequest:     true,
			expBody:        nil,
			expCalls:       1,
			expLog:         "test error",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			stdout := bytes.NewBuffer(nil)
			log.SetOutput(stdout)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				fmt.Fprintf(w, tc.body.String())
			}))
			defer ts.Close()
			q := &rqm.Queue{}
			p := &TestProcessor{Error: tc.processorError}
			var rm *rqm.RequestMaker
			if tc.addRequest {
				rm = rqm.NewRequestMaker(p, q, 1, 2, &rqm.Rq{URL: ts.URL})
			} else {
				rm = rqm.NewRequestMaker(p, q, 1, 2)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()
			rm.Run(ctx)
			if tc.expBody.String() != p.Body.String() {
				t.Fatalf("expected %s, got %s", tc.expBody, p.Body)
			}
			if p.Calls != tc.expCalls {
				t.Fatalf("expected %d calls, got: %d", tc.expCalls, p.Calls)
			}
			if !strings.HasSuffix(strings.TrimSpace(stdout.String()), tc.expLog) {
				t.Fatalf("expected log ending with %s, got: %s", tc.expLog, stdout.String())
			}
		})
	}
}
