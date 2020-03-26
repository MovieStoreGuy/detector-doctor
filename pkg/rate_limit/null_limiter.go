package ratelimit

import "context"

// NullLimiter satisfies the Limiter interface and should only be used
// for testing purposes
type NullLimiter struct{}

// Consume will allow anything and always return nil
func (*NullLimiter) Consume(_ context.Context) error { return nil }
