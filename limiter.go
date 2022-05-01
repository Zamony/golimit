/*
Simple goroutine-safe rate-limiter (token-bucket algorithm).

Example:
  // Create a new rate-limiter, allowing up-to 10 calls per second
  lim := golimit.New(10, time.Second)

  for i := 0; i < 20; i++ {
    if lim.Limit(1) {
      fmt.Println("Over limit!")
    } else {
      fmt.Println("OK")
    }
  }
*/

package golimit

import (
	"sync"
	"time"
)

// Limiter is a goroutine-safe rate-limiter,
// which implements token-bucket algorithm.
type Limiter struct {
	mu     *sync.Mutex
	limit  int64
	period int64
	last   int64
	curr   int64
}

// New creates a new limiter with specified limit and period.
func New(limit int64, period time.Duration) *Limiter {
	return &Limiter{
		limit:  limit,
		period: period.Microseconds(),
		last:   0,
		curr:   limit,
		mu:     &sync.Mutex{},
	}
}

// Limit returns true if an action was rejected.
// It accepts positive weight of an action as an argument.
func (l *Limiter) Limit(n int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now().UnixMicro()
	l.curr += (now - l.last) * l.limit / l.period
	l.last = now
	if l.curr > l.limit {
		l.curr = l.limit
	}

	if l.curr < n {
		return true
	}

	l.curr -= n
	return false
}

// Up increases the current possible weight by n.
func (l *Limiter) Up(n int64) {
	l.mu.Lock()
	l.curr += n
	if l.curr > l.limit {
		l.curr = l.limit
	}
	l.mu.Unlock()
}
