package logger

import (
	"context"
	"log/slog"
	"os"

	utilCtx "go-grpc/internal/util/context"
)

type SlogHandler struct {
	slog.Handler
}

func (h *SlogHandler) Handle(ctx context.Context, r slog.Record) error {
	requestId, ok := ctx.Value(utilCtx.XRequestId).(string)
	if ok {
		r.AddAttrs(slog.Attr{Key: "requestId", Value: slog.String("requestId", requestId).Value})
	}

	requestSource, ok := ctx.Value(utilCtx.XRequestSource).(string)
	if ok {
		r.AddAttrs(slog.Attr{Key: "requestSource", Value: slog.String("requestSource", requestSource).Value})
	}

	uid, ok := ctx.Value(utilCtx.XUid).(string)
	if ok {
		r.AddAttrs(slog.Attr{Key: "uid", Value: slog.String("uid", uid).Value})
	}

	return h.Handler.Handle(ctx, r)
}

var slogHandler = &SlogHandler{
	slog.NewTextHandler(os.Stdout, nil),
}

var logger = slog.New(slogHandler)

func Info(ctx context.Context, message string) {
	env := os.Getenv("ENV")
	if env != "testing" {
		logger.InfoContext(ctx, message)
	}
}

func Warn(ctx context.Context, message string) {
	env := os.Getenv("ENV")
	if env != "testing" {
		logger.WarnContext(ctx, message)
	}
}

func Error(ctx context.Context, message string) {
	env := os.Getenv("ENV")
	if env != "testing" {
		logger.ErrorContext(ctx, message)
	}
}
