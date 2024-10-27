package logger

import (
	"context"
	"log/slog"
)

// logHandler provides a wrapper around the slog handler to capture which
// log level is being logged for event handling.
type logHandler struct {
	slog.Handler
	events Events
}

func newLogHander(handler slog.Handler, events Events) *logHandler {
	return &logHandler{
		Handler: handler,
		events:  events,
	}
}

func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logHandler{Handler: h.Handler.WithAttrs(attrs), events: h.events}
}

func (h *logHandler) WithGroup(name string) slog.Handler {
	return &logHandler{Handler: h.Handler.WithGroup(name), events: h.events}
}

func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	switch r.Level {
	case slog.LevelDebug:
		if h.events.Debug != nil {
			h.events.Debug(ctx, toRecord(r))
		}
	case slog.LevelError:
		if h.events.Error != nil {
			h.events.Error(ctx, toRecord(r))
		}

	case slog.LevelWarn:
		if h.events.Warn != nil {
			h.events.Warn(ctx, toRecord(r))
		}

	case slog.LevelInfo:
		if h.events.Info != nil {
			h.events.Info(ctx, toRecord(r))
		}
	}

	return h.Handler.Handle(ctx, r)
}
