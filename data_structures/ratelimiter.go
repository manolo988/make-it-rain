package data_structures

import (
	"sync"
	"time"
)

// TokenBucket implements token bucket rate limiting algorithm
type TokenBucket struct {
	capacity     int           // Maximum number of tokens
	tokens       int           // Current tokens
	refillRate   int           // Tokens added per interval
	refillPeriod time.Duration // Interval for refill
	lastRefill   time.Time
	mu           sync.Mutex
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity, refillRate int, refillPeriod time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:     capacity,
		tokens:       capacity, // Start with full bucket
		refillRate:   refillRate,
		refillPeriod: refillPeriod,
		lastRefill:   time.Now(),
	}
}

// Allow checks if request is allowed (consumes token)
func (tb *TokenBucket) Allow() bool {
	return tb.AllowN(1)
}

// AllowN checks if n requests are allowed
func (tb *TokenBucket) AllowN(n int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}
	return false
}

// refill adds tokens based on time elapsed
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int(elapsed / tb.refillPeriod) * tb.refillRate

	if tokensToAdd > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+tokensToAdd)
		tb.lastRefill = now
	}
}

// AvailableTokens returns current token count
func (tb *TokenBucket) AvailableTokens() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	return tb.tokens
}

// SlidingWindowCounter implements sliding window rate limiting
type SlidingWindowCounter struct {
	windowSize time.Duration
	maxCount   int
	requests   []time.Time
	mu         sync.Mutex
}

// NewSlidingWindowCounter creates a sliding window counter
func NewSlidingWindowCounter(windowSize time.Duration, maxCount int) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		windowSize: windowSize,
		maxCount:   maxCount,
		requests:   make([]time.Time, 0),
	}
}

// Allow checks if request is allowed
func (swc *SlidingWindowCounter) Allow() bool {
	swc.mu.Lock()
	defer swc.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-swc.windowSize)

	// Remove old requests outside window
	validRequests := make([]time.Time, 0)
	for _, reqTime := range swc.requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}
	swc.requests = validRequests

	// Check if we can allow new request
	if len(swc.requests) < swc.maxCount {
		swc.requests = append(swc.requests, now)
		return true
	}

	return false
}

// Count returns current request count in window
func (swc *SlidingWindowCounter) Count() int {
	swc.mu.Lock()
	defer swc.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-swc.windowSize)

	count := 0
	for _, reqTime := range swc.requests {
		if reqTime.After(windowStart) {
			count++
		}
	}

	return count
}

// FixedWindowCounter implements fixed window rate limiting
type FixedWindowCounter struct {
	windowSize   time.Duration
	maxCount     int
	currentCount int
	windowStart  time.Time
	mu           sync.Mutex
}

// NewFixedWindowCounter creates a fixed window counter
func NewFixedWindowCounter(windowSize time.Duration, maxCount int) *FixedWindowCounter {
	return &FixedWindowCounter{
		windowSize:   windowSize,
		maxCount:     maxCount,
		currentCount: 0,
		windowStart:  time.Now(),
	}
}

// Allow checks if request is allowed
func (fwc *FixedWindowCounter) Allow() bool {
	fwc.mu.Lock()
	defer fwc.mu.Unlock()

	now := time.Now()

	// Check if we need to reset window
	if now.Sub(fwc.windowStart) >= fwc.windowSize {
		fwc.windowStart = now
		fwc.currentCount = 0
	}

	// Check if we can allow request
	if fwc.currentCount < fwc.maxCount {
		fwc.currentCount++
		return true
	}

	return false
}

// Count returns current count in window
func (fwc *FixedWindowCounter) Count() int {
	fwc.mu.Lock()
	defer fwc.mu.Unlock()

	now := time.Now()
	if now.Sub(fwc.windowStart) >= fwc.windowSize {
		return 0
	}
	return fwc.currentCount
}

// LeakyBucket implements leaky bucket algorithm
type LeakyBucket struct {
	capacity  int           // Bucket capacity
	leakRate  int           // Requests processed per interval
	leakPeriod time.Duration // Interval for leak
	queue     []time.Time   // Request queue
	lastLeak  time.Time
	mu        sync.Mutex
}

// NewLeakyBucket creates a leaky bucket
func NewLeakyBucket(capacity, leakRate int, leakPeriod time.Duration) *LeakyBucket {
	lb := &LeakyBucket{
		capacity:   capacity,
		leakRate:   leakRate,
		leakPeriod: leakPeriod,
		queue:      make([]time.Time, 0),
		lastLeak:   time.Now(),
	}

	// Start leak goroutine
	go lb.startLeaking()

	return lb
}

// Allow tries to add request to bucket
func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.queue) < lb.capacity {
		lb.queue = append(lb.queue, time.Now())
		return true
	}
	return false
}

// startLeaking continuously leaks requests
func (lb *LeakyBucket) startLeaking() {
	ticker := time.NewTicker(lb.leakPeriod)
	defer ticker.Stop()

	for range ticker.C {
		lb.leak()
	}
}

// leak removes requests from bucket
func (lb *LeakyBucket) leak() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// Remove up to leakRate requests
	toRemove := min(lb.leakRate, len(lb.queue))
	if toRemove > 0 {
		lb.queue = lb.queue[toRemove:]
	}
}

// Size returns current queue size
func (lb *LeakyBucket) Size() int {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return len(lb.queue)
}

// MultiTierRateLimiter combines multiple rate limiters
type MultiTierRateLimiter struct {
	limiters []RateLimiter
	mu       sync.RWMutex
}

// RateLimiter interface for different implementations
type RateLimiter interface {
	Allow() bool
}

// NewMultiTierRateLimiter creates multi-tier rate limiter
func NewMultiTierRateLimiter(limiters ...RateLimiter) *MultiTierRateLimiter {
	return &MultiTierRateLimiter{
		limiters: limiters,
	}
}

// Allow checks all rate limiters
func (mt *MultiTierRateLimiter) Allow() bool {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	// All limiters must allow
	for _, limiter := range mt.limiters {
		if !limiter.Allow() {
			return false
		}
	}
	return true
}

// UserRateLimiter implements per-user rate limiting
type UserRateLimiter struct {
	limiters map[string]RateLimiter
	factory  func() RateLimiter
	mu       sync.RWMutex
}

// NewUserRateLimiter creates per-user rate limiter
func NewUserRateLimiter(factory func() RateLimiter) *UserRateLimiter {
	return &UserRateLimiter{
		limiters: make(map[string]RateLimiter),
		factory:  factory,
	}
}

// Allow checks if user request is allowed
func (url *UserRateLimiter) Allow(userID string) bool {
	url.mu.Lock()
	limiter, exists := url.limiters[userID]
	if !exists {
		limiter = url.factory()
		url.limiters[userID] = limiter
	}
	url.mu.Unlock()

	return limiter.Allow()
}

// Remove removes user's rate limiter
func (url *UserRateLimiter) Remove(userID string) {
	url.mu.Lock()
	defer url.mu.Unlock()
	delete(url.limiters, userID)
}

// CircuitBreaker implements circuit breaker pattern
type CircuitBreaker struct {
	maxFailures      int
	resetTimeout     time.Duration
	halfOpenRequests int
	failureCount     int
	successCount     int
	lastFailureTime  time.Time
	state            CircuitState
	mu               sync.Mutex
}

type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// NewCircuitBreaker creates a circuit breaker
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:      maxFailures,
		resetTimeout:     resetTimeout,
		halfOpenRequests: maxFailures / 2,
		state:            StateClosed,
	}
}

// Allow checks if request is allowed
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true

	case StateOpen:
		// Check if we should transition to half-open
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			cb.state = StateHalfOpen
			cb.successCount = 0
			cb.failureCount = 0
			return true
		}
		return false

	case StateHalfOpen:
		// Allow limited requests
		totalRequests := cb.successCount + cb.failureCount
		return totalRequests < cb.halfOpenRequests
	}

	return false
}

// RecordSuccess records successful request
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateHalfOpen:
		cb.successCount++
		if cb.successCount >= cb.halfOpenRequests {
			cb.state = StateClosed
			cb.failureCount = 0
		}

	case StateClosed:
		cb.failureCount = 0
	}
}

// RecordFailure records failed request
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount++
	cb.lastFailureTime = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failureCount >= cb.maxFailures {
			cb.state = StateOpen
		}

	case StateHalfOpen:
		cb.state = StateOpen
	}
}

// GetState returns current circuit state
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}