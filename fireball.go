package fireball

import (
    "log/slog"
    "net/http"
)

type Config struct {
    Addr string
    Log  Log
}

type Log struct {
    Level slog.Level
    Json  bool
}

type Engine struct {
    Config *Config
    Logger *slog.Logger
    mux    *http.ServeMux
}

type HandlerFunc func(*Context) (*Response, error)

func New(config *Config) *Engine {
    mux := http.NewServeMux()
    logger := createLogger(config.Log.Level, config.Log.Json)

    return &Engine{
        Config: config,
        Logger: logger,
        mux:    mux,
    }
}

// Default uses a predefined configuration for the fireball engine
// The default configuration is to listen to all hosts on port 8080
func Default() *Engine {
    config := &Config{
        Addr: ":8080",
        Log: Log{
            Level: slog.LevelInfo,
            Json:  true,
        },
    }

    return New(config)
}

func (e Engine) Run() error {
    return http.ListenAndServe(e.Config.Addr, e.mux)
}
