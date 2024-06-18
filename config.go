package fireball

import "log/slog"

type Config struct {
    Addr string
    Log  Log
}

type Log struct {
    Level slog.Level
    Json  bool
}
