package fireball

import (
    "encoding/json"
    "net/http"
)

type Response struct {
    Content    []byte
    StatusCode int
    Headers    map[string]string
    Cookies    []*http.Cookie
}

func (r Response) SetHeader(key string, value string) {
    if r.Headers == nil {
        r.Headers = make(map[string]string)
    }
    r.Headers[key] = value
}

func (r Response) SetContentType(contentType string) {
    r.SetHeader("Content-Type", contentType)
}

func (r Response) AddCookie(cookie *http.Cookie) {
    if r.Cookies == nil {
        r.Cookies = make([]*http.Cookie, 0)
    }
    r.Cookies = append(r.Cookies, cookie)
}

func (r Response) SetStatusCode(code int) {
    r.StatusCode = code
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
        Content:    data,
        StatusCode: 200,
        Headers:    headers,
        Cookies:    make([]*http.Cookie, 0),
    }
}
