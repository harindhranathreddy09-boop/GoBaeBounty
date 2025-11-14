package ratelimit

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

// Limiter controls request rate limiting and concurrency
type Limiter struct {
	limiter *rate.Limiter
	workers int
	sem     chan struct{}
	mu      sync.Mutex
}

// NewLimiter creates new Limiter with maxRate requests per second and concurrency workers
func NewLimiter(maxRate, workers int) *Limiter {
	return &Limiter{
		limiter: rate.NewLimiter(rate.Limit(maxRate), maxRate),
		workers: workers,
		sem:     make(chan struct{}, workers),
	}
}

// Wait blocks until a rate token and a worker slot are available
func (l *Limiter) Wait(ctx context.Context) error {
	if err := l.limiter.Wait(ctx); err != nil {
		return err
	}

	select {
	case l.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release frees a worker slot
func (l *Limiter) Release() {
	select {
	case <-l.sem:
	default:
	}
}
