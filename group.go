package fireball

import (
	"net/http"
	"strings"

	"github.com/d3code/xlog"
	"golang.org/x/net/websocket"
)

type GroupContext struct {
    prefix       string
    middleware   []func(http.Handler) http.HandlerFunc
    routeMap     map[string]HandlerFunc
    websocketMap map[string]websocket.Server
}

func (e *Engine) Group(prefix string) *GroupContext {
    if !strings.HasPrefix(prefix, "/") {
        prefix = "/" + prefix
    }

    if !strings.HasSuffix(prefix, "/") {
        prefix += "/"
    }

    if group, ok := e.GroupMap[prefix]; ok {
        return group
    }

    group := &GroupContext{
        prefix:       prefix,
        middleware:   make([]func(http.Handler) http.HandlerFunc, 0),
        routeMap:     make(map[string]HandlerFunc),
        websocketMap: make(map[string]websocket.Server),
    }

    e.GroupMap[prefix] = group
    return group
}

func (g *GroupContext) CORS() {
    corsMiddleware := func(h http.Handler) http.HandlerFunc {
        interceptor := func(w http.ResponseWriter, r *http.Request) {
            if origin := r.Header.Get("Origin"); origin != "" {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
                    headers := []string{"Content-Type", "Accept", "Authorization", "X-Access-Token", "X-Refresh-Token"}
                    methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}

                    w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
                    w.Header().Set("Access-Control-Expose-Headers", strings.Join(headers, ","))
                    w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))

                    return
                }
            }
            h.ServeHTTP(w, r)
        }
        return interceptor
    }
    g.Use(corsMiddleware)
}

func (g *GroupContext) Use(middleware ...func(http.Handler) http.HandlerFunc) {
    g.middleware = append(g.middleware, middleware...)
}

func (g *GroupContext) Route(route string, handler HandlerFunc) {
    g.routeMap[route] = handler
}

func (g *GroupContext) RouteWs(route string, action func(ws *websocket.Conn)) {
    config := &websocket.Config{
        Origin: nil,
    }

    wsHandler := websocket.Server{
        Handler: action,
        Config:  *config,
    }

    g.websocketMap[route] = wsHandler
}

func (g *GroupContext) getHandler(handler HandlerFunc, e *Engine) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        c := &Context{
            r: req,
        }

        r, handlerErr := handler(c)
        if handlerErr != nil {
            _, writeError := w.Write([]byte(handlerErr.Error()))
            if writeError != nil {
                xlog.Error(writeError.Error())
            }
            return
        }

        if r == nil {
            return
        }

        for key, value := range r.Headers {
            w.Header().Set(key, value)
        }

        for _, cookie := range r.Cookies {
            http.SetCookie(w, cookie)
        }

        w.WriteHeader(r.StatusCode)

        _, writeError := w.Write(r.Content)
        if writeError != nil {
            xlog.Error(writeError.Error())
        }
    }
}
