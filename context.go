package fly

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

const (
	HeaderContentType = "Content-Type"
)

const (
	MIMEApplicationJSON            = "application/json"
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8
)

const (
	charsetUTF8 = "charset=UTF-8"
)

type Context interface {
	Next() error

	Request() *http.Request
	Response() http.ResponseWriter

	RealIp() string
	Scheme() string
	Path() string

	/*query*/
	Query(name string) string
	QueryAll() url.Values

	/*form*/
	Form(name string) string
	FormAll() url.Values

	/*cookie*/
	Cookie(name string) (*http.Cookie, error)
	CookieAll() []*http.Cookie

	// ctx data
	Get(key string) interface{}
	Set(key string, val interface{})

	// parse form data
	Bind(i interface{}) error

	// sends a JSON response with status code.
	JSON(code int, i interface{}) error
}

type context struct {
	request  *http.Request
	response http.ResponseWriter
	path     string
	query    url.Values
	Values   map[string]interface{}
	fly      *Fly
}

func NewContext(request *http.Request, response http.ResponseWriter) Context {
	return &context{
		request:  request,
		response: response,
		Values:   make(map[string]interface{}),
	}
}

func (ctx *context) Next() error {
	return nil
}

func (ctx *context) Request() *http.Request {
	return ctx.request
}

func (ctx *context) Response() http.ResponseWriter {
	return ctx.response
}

func (ctx *context) RealIp() string {
	return ctx.request.RemoteAddr
}

func (ctx *context) Scheme() string {
	return ctx.request.Proto
}

func (ctx *context) Path() string {
	return ctx.request.RequestURI
}

func (ctx *context) Query(name string) string {
	return ctx.request.URL.Query().Get(name)
}

func (ctx *context) QueryAll() url.Values {
	return ctx.request.URL.Query()
}

func (ctx *context) Form(name string) string {
	return ctx.request.FormValue(name)
}

func (ctx *context) FormAll() url.Values {
	return ctx.request.Form
}

func (ctx *context) Cookie(name string) (*http.Cookie, error) {
	return ctx.Cookie(name)
}

func (ctx *context) CookieAll() []*http.Cookie {
	return ctx.CookieAll()
}

func (ctx *context) Get(key string) interface{} {
	return ctx.Values[key]
}

func (ctx *context) Set(key string, val interface{}) {
	ctx.Values[key] = val
}

func (ctx *context) Bind(i interface{}) error {
	return errors.New("not impliments")
}

func (ctx *context) JSON(code int, i interface{}) error {
	ctx.response.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	ctx.response.WriteHeader(code)
	_, err = ctx.response.Write(b)
	return err
}
