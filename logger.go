package fireball

import (
    "log/slog"
    "os"
)

func createLogger(level slog.Level, logJson bool) *slog.Logger {
    handlerOpts := &slog.HandlerOptions{
        AddSource: level == slog.LevelDebug,
        Level:     level,
    }

    var handler slog.Handler
    if logJson {
        handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
    } else {
        handler = slog.NewTextHandler(os.Stdout, handlerOpts)
    }

    logger := slog.New(handler).With(
        slog.String("service", "fireball"),
        slog.String("user", os.Getenv("USER")),
    )
    return logger
}
