package resilience

import (
	"log"
	"sync"
)

// Simple Circuit Breaker to protect downstream AI APIs
type CircuitBreaker struct {
	mu           sync.Mutex
	FailureCount int
	Threshold    int
	IsOpen       bool
}

func (cb *CircuitBreaker) Execute(task func() error) error {
	cb.mu.Lock()
	if cb.IsOpen {
		cb.mu.Unlock()
		return errors.New("circuit_breaker_open: rejecting_request_to_protect_system")
	}
	cb.mu.Unlock()

	err := task()
	if err != nil {
		cb.recordFailure()
		return err
	}
	return nil
}

func (cb *CircuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.FailureCount++
	if cb.FailureCount >= cb.Threshold {
		cb.IsOpen = true
		log.Println("CRITICAL: Circuit Breaker TRIPPED. Protecting downstream API.")
	}
}
