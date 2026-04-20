package balancer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestInvalidBalancerCreation_ShouldFail(t *testing.T) {
	_, err := NewBalancer([]string{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBalancerAlternates_ShouldPass(t *testing.T) {
	balancers := []string{"upstream-a ", "upstream-a "}
	balancer, err := NewBalancer(balancers)

	if err != nil {
		t.Fatal(err)
	}

	for _, b := range balancers {

		nextStream, err := balancer.Next()

		if err != nil {
			t.Error(err)
		}

		if nextStream == "" {
			t.Fatalf("expected next stream")
		}

		if nextStream != b {
			t.Errorf("expected next stream will be %s", b)
		}
	}
}
func TestBalancerConcurrentNext_roundRobinStaysBalanced_ShouldPass(t *testing.T) {

	const (
		workers   = 32
		perWorker = 32
	)

	wantTotal := workers * perWorker

	urlA := "http://upstream-a.test"
	urlB := "http://upstream-b.test"

	bal, err := NewBalancer([]string{urlA, urlB})
	if err != nil {
		t.Fatal(err)
	}

	var countA, countB atomic.Int64
	var wg sync.WaitGroup
	var errMu sync.Mutex
	var firstErr error

	fail := func(format string, args ...any) {
		errMu.Lock()
		if firstErr == nil {
			firstErr = fmt.Errorf(format, args...)
		}
		errMu.Unlock()
	}
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := 0; i < perWorker; i++ {
				got, err := bal.Next()

				if err != nil {
					fail("Next: %w", err)
					return
				}

				switch got {
				case urlA:
					countA.Add(1)
				case urlB:
					countB.Add(1)
				default:
					fail("unexpected target %q", got)
					return
				}
			}
		}()
	}
	wg.Wait()

	errMu.Lock()
	e := firstErr
	errMu.Unlock()

	if e != nil {
		t.Fatal(e)
	}

	a := countA.Load()
	b := countB.Load()

	if a+b != int64(wantTotal) {
		t.Fatalf("total calls %d, want %d", a+b, wantTotal)
	}

	d := a - b

	if d < 0 {
		d = -d
	}

	if d > 1 {
		t.Fatalf("round-robin imbalance under concurrency: A=%d B=%d (want |A-B| <= 1)", a, b)
	}
}
