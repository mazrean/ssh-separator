package domain

import (
	"errors"
	"sync/atomic"
)

var (
	ErrTooManyConnections = errors.New("too many connections")
)

// ConnectionLimiter manages global connection limits
type ConnectionLimiter struct {
	maxConnections     int64
	currentConnections int64
}

// NewConnectionLimiter creates a new connection limiter with the specified maximum
func NewConnectionLimiter(maxConnections int64) *ConnectionLimiter {
	return &ConnectionLimiter{
		maxConnections:     maxConnections,
		currentConnections: 0,
	}
}

// TryAcquire attempts to acquire a connection slot
// Returns an error if the limit has been reached
func (cl *ConnectionLimiter) TryAcquire() error {
	for {
		current := atomic.LoadInt64(&cl.currentConnections)
		if current >= cl.maxConnections {
			return ErrTooManyConnections
		}
		if atomic.CompareAndSwapInt64(&cl.currentConnections, current, current+1) {
			return nil
		}
	}
}

// Release releases a connection slot
func (cl *ConnectionLimiter) Release() {
	atomic.AddInt64(&cl.currentConnections, -1)
}

// CurrentConnections returns the current number of connections
func (cl *ConnectionLimiter) CurrentConnections() int64 {
	return atomic.LoadInt64(&cl.currentConnections)
}

// MaxConnections returns the maximum number of connections allowed
func (cl *ConnectionLimiter) MaxConnections() int64 {
	return cl.maxConnections
}
