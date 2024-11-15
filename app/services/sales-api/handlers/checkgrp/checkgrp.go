package checkgrp

import (
	"context"
	"net/http"
	"os"
	"sales-api/business/data/dbsql/pgx"
	"sales-api/foundation/logger"
	"sales-api/foundation/web"
	"time"

	"github.com/jmoiron/sqlx"
)

type Handlers struct {
	build string
	log   *logger.Logger
	db    *sqlx.DB
}

func New(build string, logger *logger.Logger, db *sqlx.DB) *Handlers {
	return &Handlers{build: build, log: logger, db: db}
}

// Readiness checks if the database is ready and if not will return a 500 status.
// Do not respond by just returning an error because further up in the call
// stack it will interpret that as a non-trusted error.
func (h *Handlers) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := "ok"
	statusCode := http.StatusOK

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := pgx.StatusCheck(ctx, h.db); err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
		h.log.Info(ctx, "readiness failure", "status", status)
	}

	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}
	return web.Respond(ctx, w, data, statusCode)
}

// Liveness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (h *Handlers) liviness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//Todo log this when it fails

	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}
	data := struct {
		Status     string `json:"status,omitempty"`
		Build      string `json:"build,omitempty"`
		Host       string `json:"host,omitempty"`
		Name       string `json:"name,omitempty"`
		PodIP      string `json:"podIP,omitempty"`
		Node       string `json:"node,omitempty"`
		Namespace  string `json:"namespace,omitempty"`
		GOMAXPROCS string `json:"GOMAXPROCS,omitempty"`
	}{
		Status:     "up",
		Build:      h.build,
		Host:       host,
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: os.Getenv("GOMAXPROCS"),
	}
	//h.log.Info(ctx, "liveness", "status", "OK")

	return web.Respond(ctx, w, data, http.StatusOK)
}
