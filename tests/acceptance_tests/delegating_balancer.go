package acceptance_tests

import (
	"errors"
)

type delegatingBalancer struct{}

func (*delegatingBalancer) Next() (string, error) {
	roundRobinState.mu.Lock()
	b := roundRobinState.bal
	roundRobinState.mu.Unlock()

	if b == nil {
		return "", errors.New("balancer not configured")
	}
	return b.Next()
}
