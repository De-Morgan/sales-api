package hackgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"sales-api/business/web/v1/response"
	"sales-api/foundation/web"
)

func hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if rand.Intn(2) == 1 {
		return response.NewError(errors.New("TRUSTED ERROR"), http.StatusBadRequest)
	}
	success := struct {
		Status string
	}{
		Status: "OK",
	}
	return web.Respond(ctx, w, success, http.StatusOK)
}
