package fireball

import (
    "log"
    "log/slog"
    "net/http"
    "strings"
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
    return http.ListenAndServe(e.Config.Addr, middlewareCORS(e.mux))
}

// middlewareCORS is a middlewareLog that adds CORS headers to the response if the Origin header is set
func middlewareCORS(h http.Handler) http.Handler {
    interceptor := func(w http.ResponseWriter, r *http.Request) {
        if origin := r.Header.Get("Origin"); origin != "" {

            w.Header().Set("Access-Control-Allow-Origin", origin)

            if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
                headers := []string{"Content-Type", "Accept", "Authorization", "X-Access-Token", "X-Refresh-Token"}
                methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}

                w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
                w.Header().Set("Access-Control-Expose-Headers", strings.Join(headers, ","))

                w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))

                log.Println("CORS preflight request")

                return
            }
        }
        h.ServeHTTP(w, r)
    }

    return http.HandlerFunc(interceptor)
}
