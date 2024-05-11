package main

import (
    "github.com/d3code/fireball"
)

func main() {

    engine := fireball.Default()

    engine.Route("GET /", func(c *fireball.Context) (*fireball.Response, error) {

        name := c.GetPathParam("name")
        query := c.GetQueryString("query")

        c.Logger.Info("Hello " + name)
        c.Logger.Info("Query " + query)

        return fireball.ResponseJson(rootResponse{
            Message: "Hello " + name,
            Query:   query,
        })
    })

    engine.Logger.Info("Server is running on " + engine.Config.Addr)
    err := engine.Run()
    if err != nil {
        engine.Logger.Error(err.Error())
        return
    }
}

type rootResponse struct {
    Message string `json:"message"`
    Query   string `json:"query"`
}
