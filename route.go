package fireball

import (
    "encoding/json"
    "github.com/google/uuid"
    "log/slog"
    "net/http"
)

type Response struct {
    content    []byte
    statusCode int
    headers    map[string]string
    cookies    []*http.Cookie
}

func (r Response) SetHeader(key string, value string) {
    if r.headers == nil {
        r.headers = make(map[string]string)
    }
    r.headers[key] = value
}

func (r Response) SetContentType(contentType string) {
    r.SetHeader("Content-Type", contentType)
}

func (r Response) AddCookie(cookie *http.Cookie) {
    if r.cookies == nil {
        r.cookies = make([]*http.Cookie, 0)
    }
    r.cookies = append(r.cookies, cookie)
}

func (r Response) SetStatusCode(code int) {
    r.statusCode = code
}

// ResponseJson returns a response with a JSON content type
func ResponseJson(data any) (*Response, error) {
    marshal, marshalErr := json.Marshal(data)
    if marshalErr != nil {
        return &Response{}, marshalErr
    }

    return ResponseBytes(marshal, "application/json"), nil
}

// ResponseText returns a response with a text content type
func ResponseText(data string) (*Response, error) {
    bytes := []byte(data)
    return ResponseBytes(bytes, "text/plain"), nil
}

// ResponseHtml returns a response with a text/html content type
func ResponseHtml(data string) (*Response, error) {
    bytes := []byte(data)
    return ResponseBytes(bytes, "text/html"), nil
}

func ResponseBytes(data []byte, contentType string) *Response {
    headers := make(map[string]string)
    headers["Content-Type"] = contentType

    return &Response{
        content:    data,
        statusCode: 200,
        headers:    headers,
        cookies:    make([]*http.Cookie, 0),
    }
}

func (e Engine) Route(route string, handler HandlerFunc) {
    newRoute := func(w http.ResponseWriter, req *http.Request) {
        traceId := uuid.New().String()
        logger := createLogger(e.Config.Log.Level, e.Config.Log.Json).With(
            slog.String("request_id", req.Header.Get("X-Request-Id")),
            slog.String("trace_id", traceId),
            slog.String("remote_addr", req.RemoteAddr),
            slog.String("method", req.Method),
            slog.String("path", req.URL.Path),
        )

        c := &Context{
            r:      req,
            Logger: logger,
        }

        r, handlerErr := handler(c)
        if handlerErr != nil {
            _, writeError := w.Write([]byte(handlerErr.Error()))
            if writeError != nil {
                logger.Error(writeError.Error())
            }
            return
        }

        if r == nil {
            return
        }

        for key, value := range r.headers {
            w.Header().Set(key, value)
        }

        for _, cookie := range r.cookies {
            http.SetCookie(w, cookie)
        }

        w.WriteHeader(r.statusCode)

        _, writeError := w.Write(r.content)
        if writeError != nil {
            logger.Error(writeError.Error())
        }
        return
    }

    e.mux.HandleFunc(route, newRoute)
}
