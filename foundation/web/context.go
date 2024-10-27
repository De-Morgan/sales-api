package web

import (
	"context"
	"time"
)

type ctxKey int

const key ctxKey = 1

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// SetValues sets the specified Values in the context.
func SetValues(ctx context.Context, v *Values) context.Context {
	return context.WithValue(ctx, key, v)
}

// GetValues returns the values from the context.
func GetValues(ctx context.Context) *Values {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return &Values{
			TraceID: "00000000-0000-0000-0000-000000000000",
			Now:     time.Now(),
		}
	}
	return v
}

func GetTraceID(ctx context.Context) string {
	val := GetValues(ctx)
	return val.TraceID
}

// GetTime returns the time from the context.
func GetTime(ctx context.Context) time.Time {
	val := GetValues(ctx)
	return val.Now
}

// SetStatusCode sets the status code back into the context.
func SetStatusCode(ctx context.Context, statusCode int) {
	val := GetValues(ctx)
	val.StatusCode = statusCode
}
