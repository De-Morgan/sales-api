package mid

import (
	"context"
	"net/http"
	"sales-api/business/web/v1/metrics"
	"sales-api/foundation/web"
)

func Metrics() web.Middleware {

	m := func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			ctx = metrics.Set(ctx)

			err = handler(ctx, w, r)
			metrics.AddRequests(ctx)
			metrics.AddGoroutines(ctx)
			if err != nil {
				metrics.AddErrors(ctx)
			}
			return
		}
	}

	return m
}
