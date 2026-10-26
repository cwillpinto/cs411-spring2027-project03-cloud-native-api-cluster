package main

import (
	"sync"
	"time"
)

// Limiter is per API replica; discuss the implications of local limits when
// an HPA adds replicas. The public tests use injected times for reproducibility.
type Limiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	used   int
	reset  time.Time
}

func NewLimiter(limit int, window time.Duration) *Limiter {
	return &Limiter{limit: limit, window: window}
}

// TODO: allow at most limit requests in each fixed window. Reset when now
// reaches the window boundary. Protect all state with mu.
func (l *Limiter) Allow(now time.Time) bool {
	_ = now
	return true
}
