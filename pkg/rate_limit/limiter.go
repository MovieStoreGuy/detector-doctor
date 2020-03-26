package ratelimit

import "context"

// Limiter defines the basic definition of a rater limiter
type Limiter interface {
	// Consume is a blocking call once you have
	Consume(context.Context) error
}
