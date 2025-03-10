package utils

import (
	"context"
	"log/slog"
	"os"
)

type MyHandler struct {
	slog.Handler
}

func (h *MyHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, ok := ctx.Value("request_id").(string); ok {
		r.AddAttrs(slog.String("request_id", id))
	}

	return h.Handler.Handle(ctx, r)
}

func (h *MyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *MyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MyHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *MyHandler) WithGroud(name string) slog.Handler {
	return &MyHandler{Handler: h.Handler.WithGroup(name)}
}

func InitLogging() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: getLogLevel()})
	myHandler := MyHandler{Handler: jsonHandler}
	logger := slog.New(&myHandler)

	slog.SetDefault(logger)
}

func getLogLevel() slog.Level {
    logLevelStr := os.Getenv("LOG_LEVEL")
    switch logLevelStr {
    case "debug":
        return slog.LevelDebug
    case "info":
        return slog.LevelInfo
    case "warn":
        return slog.LevelWarn
    case "error":
        return slog.LevelError
    default:
        return slog.LevelInfo
    }
}
