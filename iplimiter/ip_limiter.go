package iplimiter

import (
	"sync"

	"go.uber.org/ratelimit"
)

type IPLimiter struct {
	ips map[string]ratelimit.Limiter
	mu  sync.Mutex
	rps int
}

type IPLimiterOption func(*IPLimiter)

func RPS(rps int) IPLimiterOption {
	return func(l *IPLimiter) {
		l.rps = rps
	}
}

func NewLimiter(opts ...IPLimiterOption) *IPLimiter {
	l := &IPLimiter{
		ips: make(map[string]ratelimit.Limiter),
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *IPLimiter) newLimiter(ip string) ratelimit.Limiter {
	l.ips[ip] = ratelimit.New(l.rps)
	return l.ips[ip]
}

func (l *IPLimiter) GetLimiter(ip string) ratelimit.Limiter {
	l.mu.Lock()

	limiter, ok := l.ips[ip]
	if ok {
		l.mu.Unlock()
		return limiter
	}

	nl := l.newLimiter(ip)
	l.mu.Unlock()

	return nl
}
