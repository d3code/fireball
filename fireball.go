package fireball

import (
	"fmt"
	"net/http"
	"time"

	"github.com/d3code/xlog"
	"github.com/google/uuid"
)

type Engine struct {
	Config   *Config
	GroupMap map[string]*GroupContext
}

type HandlerFunc func(*Context) (*Response, error)

func New(config *Config) *Engine {
	return &Engine{
		Config:   config,
		GroupMap: make(map[string]*GroupContext),
	}
}

// Default uses a predefined configuration for the fireball engine.
// The default configuration is to listen to all hosts at a free port as reported by the operating system or 8080 if no free port is found.
func Default() *Engine {
	port, err := GetFreePort()
	if err != nil {
		port = 8080
	}

	config := &Config{
		Host: "",
		Port: port,
		Log: Log{
			Request: true,
		},
	}

	return New(config)
}

func (e *Engine) Addr() string {
	return fmt.Sprintf("%s:%d", e.Config.Host, e.Config.Port)
}

func (e *Engine) Run() error {
	muxBase := http.NewServeMux()

	for path, group := range e.GroupMap {
		mux := http.NewServeMux()

		for route, handler := range group.routeMap {
			mux.HandleFunc(route, group.getHandler(handler, e))
		}

		for route, handler := range group.websocketMap {
			mux.Handle(route, handler)
		}

		var wrappedHandler http.HandlerFunc = mux.ServeHTTP
		for _, middleware := range group.middleware {
			wrappedHandler = middleware(wrappedHandler)
		}

		groupMiddleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if path != "/" {
					newPath := r.URL.Path[len(path)-1:]
					r.URL.Path = newPath
				}
				next.ServeHTTP(w, r)
			})
		}

		muxBase.Handle(path, groupMiddleware(wrappedHandler))
	}

	logMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := uuid.New().String()
			r.Header.Set("X-Request-Id", requestId)

			start := time.Now()
			responseWriter := NewLoggingResponseWriter(w)
			next.ServeHTTP(responseWriter, r)

			duration := time.Since(start)
			if e.Config.Log.Request {
				xlog.Debugf("[http] %s %s | %s | %d | %v", r.Method, r.RequestURI, requestId, responseWriter.statusCode, duration)
			}
		})
	}

	return http.ListenAndServe(e.Addr(), logMiddleware(muxBase))
}
