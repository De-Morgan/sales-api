package mid

import (
	"context"
	"net/http"
	"sales-api/business/web/v1/auth"
	"sales-api/business/web/v1/response"
	"sales-api/foundation/logger"
	"sales-api/foundation/validate"
	"sales-api/foundation/web"
)

func Errors(log *logger.Logger) web.Middleware {

	m := func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			if err := handler(ctx, w, r); err != nil {
				log.Error(ctx, "message", "msg", err)

				var er response.ErrorDocument
				var status int

				switch {

				//trusted error
				case response.IsError(err):
					reqErr := response.GetError(err)

					if validate.IsFieldErrors(reqErr.Err) {
						fieldErrors := validate.GetFieldErrors(reqErr.Err)
						er = response.ErrorDocument{
							Error:  "data validation error",
							Fields: fieldErrors.Fields(),
						}
						status = reqErr.Status
						break
					}
					er = response.ErrorDocument{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				case auth.IsAuthError(err):
					er = response.ErrorDocument{
						Error: http.StatusText(http.StatusUnauthorized),
					}
					status = http.StatusUnauthorized

				default:
					er = response.ErrorDocument{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError

				}

				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}

			}

			return nil
		}
	}

	return m
}
