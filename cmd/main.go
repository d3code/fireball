package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/d3code/fireball"
	"github.com/d3code/xlog"
)

func main() {

	engine := fireball.Default()

	group := engine.Group("/go")
	group.Route("GET /", func(c *fireball.Context) (*fireball.Response, error) {

		name := c.GetQueryString("name")
		query := c.GetQueryString("query")

		response, _ := fireball.ResponseJson(rootResponse{
			Message: "Hello " + name,
			Query:   query,
		})

		response.StatusCode = 201
		return response, nil
	})

	go func() {
		time.Sleep(1 * time.Second)
		x, err := http.Get(fmt.Sprintf("http://localhost:%v/go?name=world&query=example", engine.Config.Port))
		if err != nil {
			xlog.Error(err.Error())
			return
		}

		response, err := io.ReadAll(x.Body)
		if err != nil {
			xlog.Error(err.Error())
			return
		}

		xlog.Info(string(response))
	}()

	xlog.Info("Server running at address " + engine.Addr())
	err := engine.Run()
	if err != nil {
		xlog.Error(err.Error())
		return
	}
}

type rootResponse struct {
	Message string `json:"message"`
	Query   string `json:"query"`
}
