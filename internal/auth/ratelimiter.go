package auth

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu        sync.Mutex
	ips       map[string]int
	limit     int
	resetTime time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration) *RateLimiter {
	return &RateLimiter{
		ips:       make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}
}

// Check if the IP exceeded the rate limit
func (r *RateLimiter) CheckRateLimit(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.ips[ip]; exists {
		if r.ips[ip] >= r.limit {
			return true // Exceeded limit
		}
		r.ips[ip]++
	} else {
		r.ips[ip] = 1
	}

	// Reset rate limit after the defined period
	go func() {
		time.Sleep(r.resetTime)
		r.mu.Lock()
		delete(r.ips, ip)
		r.mu.Unlock()
	}()
	return false // Allowed request
}