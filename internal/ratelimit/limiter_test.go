package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestLimiterBasic(t *testing.T) {
	lim := NewLimiter(10, 3)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		err := lim.Wait(ctx)
		if err != nil {
			t.Fatalf("Wait failed: %v", err)
		}
		lim.Release()
	}
}

func TestLimiterRateLimit(t *testing.T) {
	lim := NewLimiter(5, 3)
	ctx := context.Background()

	start := time.Now()

	for i := 0; i < 10; i++ {
		err := lim.Wait(ctx)
		if err != nil {
			t.Fatalf("Rate limit wait error: %v", err)
		}
		lim.Release()
	}

	duration := time.Since(start)
	if duration < 2*time.Second {
		t.Errorf("Expected rate limit delay, got duration %v", duration)
	}
}
