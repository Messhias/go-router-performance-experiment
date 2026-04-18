package balancer

import "testing"

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
