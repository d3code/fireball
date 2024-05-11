package fireball

import (
    "log/slog"
    "net/http"
)

type Context struct {
    r      *http.Request
    body   []byte
    Logger *slog.Logger
}

func (c Context) GetAuthorization() string {
    return c.r.Header.Get("Authorization")
}

func (c Context) GetHeader(header string) string {
    return c.r.Header.Get(header)
}

func (c Context) GetHeaders() http.Header {
    return c.r.Header
}

func (c Context) GetMethod() string {
    return c.r.Method
}

func (c Context) GetPath() string {
    return c.r.URL.Path
}

func (c Context) GetPathParam(param string) string {
    return c.r.PathValue(param)
}

func (c Context) GetQueryString(query string) string {
    return c.r.URL.Query().Get(query)
}

func (c Context) GetCookie(name string) *http.Cookie {
    cookie, err := c.r.Cookie(name)
    if err != nil {
        return nil
    }
    return cookie
}

func (c Context) GetBody() []byte {
    if c.body != nil {
        return c.body
    }

    body := make([]byte, c.r.ContentLength)
    _, err := c.r.Body.Read(body)
    if err != nil {
        c.Logger.Error(err.Error())
        return nil
    }

    c.body = body
    return body
}

func (c Context) GetBodyString() string {
    return string(c.GetBody())
}
