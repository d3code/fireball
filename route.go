package fireball

import (
	"golang.org/x/net/websocket"
)

func (e *Engine) Route(route string, handler HandlerFunc) {
	group := e.Group("/")
	group.Route(route, handler)
}

func (e *Engine) RouteWs(route string, action func(ws *websocket.Conn)) {
	group := e.Group("/")
	group.RouteWs(route, action)
}
