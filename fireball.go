package fireball

import (
    "fmt"
    "github.com/d3code/xlog"
    "log/slog"
    "net/http"
)

type Engine struct {
    Config   *Config
    GroupMap map[string]*GroupContext
    Logger   *slog.Logger
}

type HandlerFunc func(*Context) (*Response, error)

func New(config *Config) *Engine {
    return &Engine{
        Config:   config,
        GroupMap: make(map[string]*GroupContext),
    }
}

// Default uses a predefined configuration for the fireball engine.
// The default configuration is to listen to all hosts on port 8080
func Default() *Engine {
    config := &Config{
        Addr: ":8080",
        Log: Log{
            Level: slog.LevelInfo,
            Json:  false,
        },
    }

    return New(config)
}

func (e Engine) Run() error {
    muxBase := http.NewServeMux()

    for path, group := range e.GroupMap {
        mux := http.NewServeMux()

        for route, handler := range group.routeMap {
            mux.HandleFunc(route, group.getHandler(handler, &e))
        }

        for route, handler := range group.websocketMap {
            mux.Handle(route, handler)
        }

        var wrappedHandler http.HandlerFunc = mux.ServeHTTP
        for _, middleware := range group.middleware {
            wrappedHandler = middleware(wrappedHandler)
        }

        logMiddleware := func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if path != "/" {
                    newPath := r.URL.Path[len(path)-1:]
                    r.URL.Path = newPath
                }
                next.ServeHTTP(w, r)
            })
        }

        muxBase.Handle(path, logMiddleware(wrappedHandler))
    }

    logMiddleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            xlog.Debug(fmt.Sprintf("Request [ %s %s ]", r.Method, r.RequestURI))
            next.ServeHTTP(w, r)
        })
    }

    return http.ListenAndServe(e.Config.Addr, logMiddleware(muxBase))
}
