package httpx

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"golang.org/x/time/rate"
)

// RateLimiter limits requests per client IP with a token bucket per IP.
// The map of buckets is bounded: idle buckets are evicted.
type RateLimiter struct {
	perMinute  int
	maxClients int
	idle       time.Duration
	now        func() time.Time

	mu      sync.Mutex
	clients map[string]*client
}

type client struct {
	lim  *rate.Limiter
	seen time.Time
}

// NewRateLimiter allows perMinute requests per IP per minute, with bursts up to perMinute.
func NewRateLimiter(perMinute int) *RateLimiter {
	return &RateLimiter{perMinute: perMinute, maxClients: 10000, idle: 10 * time.Minute, now: time.Now, clients: map[string]*client{}}
}

// Allow reports whether a request from ip may proceed, and if not, how long to wait.
func (l *RateLimiter) Allow(ip string) (bool, time.Duration) {
	now := l.now()
	l.mu.Lock()
	c, ok := l.clients[ip]
	if !ok {
		if len(l.clients) >= l.maxClients {
			l.evict(now)
		}
		c = &client{lim: rate.NewLimiter(rate.Limit(float64(l.perMinute)/60), l.perMinute)}
		l.clients[ip] = c
	}
	c.seen = now
	l.mu.Unlock()

	r := c.lim.ReserveN(now, 1)
	if !r.OK() {
		return false, time.Minute
	}
	if d := r.DelayFrom(now); d > 0 {
		r.CancelAt(now)
		return false, d
	}
	return true, 0
}

// evict drops idle buckets; if none are idle, it drops the oldest half. Caller holds mu.
func (l *RateLimiter) evict(now time.Time) {
	for ip, c := range l.clients {
		if now.Sub(c.seen) > l.idle {
			delete(l.clients, ip)
		}
	}
	if len(l.clients) < l.maxClients {
		return
	}
	cutoff := now
	for _, c := range l.clients {
		if c.seen.Before(cutoff) {
			cutoff = c.seen
		}
	}
	mid := cutoff.Add(now.Sub(cutoff) / 2)
	for ip, c := range l.clients {
		if !c.seen.After(mid) {
			delete(l.clients, ip)
		}
	}
}

// Len returns the number of tracked clients.
func (l *RateLimiter) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.clients)
}

// HumaMiddleware rejects over-limit requests with 429, Retry-After, and a problem body.
// The client IP comes from Gin, which honors forwarding headers only from trusted proxies.
func (l *RateLimiter) HumaMiddleware(ctx huma.Context, next func(huma.Context)) {
	c := humagin.Unwrap(ctx)
	if ok, wait := l.Allow(c.ClientIP()); !ok {
		c.Header("Retry-After", strconv.Itoa(int(math.Ceil(wait.Seconds()))))
		WriteProblem(c, http.StatusTooManyRequests, "too many search requests; try again later")
		return
	}
	next(ctx)
}
