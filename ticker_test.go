package rqm_test

import (
	"testing"
	"time"

	"github.com/fwojciec/rqm"
)

func TestRandomTicker(t *testing.T) {
	t.Parallel()
	rt := rqm.NewRandomTicker(100, 200)
	for i := 0; i < 3; i++ {
		t0 := time.Now()
		<-rt.C
		td := time.Since(t0)
		if td < 100*time.Millisecond {
			t.Fatalf("tick was shorter than expected: %s", td)
		} else if td > 200*time.Millisecond {
			t.Fatalf("tick was longer than expected: %s", td)
		}
	}
	rt.Stop()
	time.Sleep(200 * time.Millisecond)
	select {
	case <-rt.C:
		t.Fatal("Ticker did not shut down")
	default:
		// ok
	}
}
