package mid

import (
	"context"
	"fmt"
	"net/http"
	"sales-api/foundation/logger"
	"sales-api/foundation/web"
	"time"
)

func Logger(log *logger.Logger) web.Middleware {

	m := func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v := web.GetValues(ctx)
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}
			log.Info(ctx, "request started", "method", r.Method, "path", path,
				"remoteaddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Info(ctx, "request completed", "method", r.Method, "path", path,
				"remoteaddr", r.RemoteAddr, "statusCode", v.StatusCode, "since", fmt.Sprintf("%d", time.Since(v.Now).Round(time.Millisecond)))

			return err
		}
	}

	return m
}
