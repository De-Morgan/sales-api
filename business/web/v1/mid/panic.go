package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"sales-api/business/web/v1/metrics"
	"sales-api/foundation/web"
)

func Panics() web.Middleware {

	m := func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
					metrics.AddPanics(ctx)
				}

			}()

			return handler(ctx, w, r)

		}
	}

	return m
}
