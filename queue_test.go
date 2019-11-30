package rqm_test

import (
	"testing"

	"github.com/fwojciec/rqm"
)

func TestQueue(t *testing.T) {
	q := &rqm.Queue{}
	if !q.IsEmpty() {
		t.Fatalf("expected a freshly initialized queue to be empty")
	}
	q.Push(&rqm.Rq{URL: "test 1"})
	q.Push(&rqm.Rq{URL: "test 2"})
	if q.IsEmpty() {
		t.Fatalf("expected a queue with two items in it to be not empty")
	}
	x1 := q.Pop()
	if x1.URL != "test 1" {
		t.Fatalf(`expected the URL to be "test 1", received %q`, x1.URL)
	}
	if q.IsEmpty() {
		t.Fatalf("expected a queue with one item in it to be not empty")
	}
	x2 := q.Pop()
	if x2.URL != "test 2" {
		t.Fatalf(`expected the URL to be "test 2", received %q`, x2.URL)
	}
	if !q.IsEmpty() {
		t.Fatalf("expected a queue with no items in it to be empty")
	}
}
