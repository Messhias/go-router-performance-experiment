package balancer

import (
	"errors"
	"sync"
)

type RoundRobin interface {
	Next() (string, error)
}

type roundRobin struct {
	next    int
	targets []string
	mu      sync.Mutex
}

func (r *roundRobin) Next() (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.targets) == 0 {
		return "", errors.New("no targets")
	}

	chosen := r.targets[r.next]
	r.next++

	if r.next == len(r.targets) {
		r.next = 0
	}

	return chosen, nil
}

func NewBalancer(balancers []string) (RoundRobin, error) {
	if len(balancers) == 0 {
		return nil, errors.New("no balancers provided")
	}

	return &roundRobin{targets: balancers}, nil
}
