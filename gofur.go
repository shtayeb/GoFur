package gofur

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

type Plug func(Handler) Handler

type Handler func(c *Context) error

type ErrorHandler func(error, *Context) error

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	ctx      context.Context
	params   httprouter.Params
}

func newContext(w http.ResponseWriter, r *http.Request, params httprouter.Params) *Context {
	return &Context{
		Response: w,
		Request:  r,
		ctx:      context.Background(),
		params:   params,
	}
}

func (c *Context) Param(name string) string {
	return c.params.ByName(name)
}

func (c *Context) Query(name string) string {
	return c.Request.URL.Query().Get(name)
}

func (c *Context) FormValue(name string) string {
	return c.Request.FormValue(name)
}

func (c *Context) Render(component templ.Component) error {
	return component.Render(c.ctx, c.Response)
}

func (c *Context) Redirect(url string, code int) error {
	if code < http.StatusMultipleChoices || code > http.StatusTemporaryRedirect {
		return errors.New("invalid redirect code")
	}
	http.Redirect(c.Response, c.Request, url, code)
	return nil
}

func (c *Context) JSON(status int, v any) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(status)
	return json.NewDecoder(c.Request.Body).Decode(&v)
}

func (c *Context) Text(status int, t string) error {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(status)
	_, err := c.Response.Write([]byte(t))
	return err
}

func (c *Context) Set(key string, value any) {
	c.ctx = context.WithValue(c.ctx, key, value)
}

func (c *Context) Get(key string) any {
	return c.ctx.Value(key)
}

type Gofur struct {
	ErrorHandler ErrorHandler
	router       *httprouter.Router
	plugs        []Plug
}

func New() *Gofur {
	return &Gofur{
		router:       httprouter.New(),
		ErrorHandler: defaultErrorHandler,
	}
}

type methodNotAllowedHandler struct {
	handler Handler
}

func (h methodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(w, r, httprouter.Params{})
	h.handler(ctx)
}

func (gf *Gofur) MethodNotAllowed(h Handler) {
	gf.router.MethodNotAllowed = methodNotAllowedHandler{h}
}

func (gf *Gofur) Plug(plugs ...Plug) {
	gf.plugs = append(gf.plugs, plugs...)
}

func (gf *Gofur) Start() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	// Retrieve and sanitize listen address from env
	listenAddr := os.Getenv("GOFUR_HTTP_LISTEN_ADDR")
	listenAddr = strings.TrimSpace(listenAddr)

	// If listen address is not set, use default host and port
	if listenAddr == "" {
		listenAddr = ":3000"
	}

	// Print the URL where the app is running
	browsableURL := listenAddr
	if strings.HasPrefix(browsableURL, ":") {
		browsableURL = "localhost" + browsableURL
	}

	fmt.Printf("Gofur app running at http://%s\n", browsableURL)

	// Start the HTTP server
	return http.ListenAndServe(listenAddr, gf.router)
}

func (gf *Gofur) add(method, path string, h Handler, plugs ...Plug) {
	gf.router.Handle(method, path, gf.makeHTTPRouterHandle(h, plugs...))
}

func (gf *Gofur) Get(path string, h Handler, plugs ...Plug) {
	gf.add("GET", path, h, plugs...)
}

func (gf *Gofur) Post(path string, h Handler, plugs ...Plug) {
	gf.add("POST", path, h, plugs...)
}

func (gf *Gofur) Put(path string, h Handler, plugs ...Plug) {
	gf.add("PUT", path, h, plugs...)
}

func (gf *Gofur) Delete(path string, h Handler, plugs ...Plug) {
	gf.add("DELETE", path, h, plugs...)
}

func (gf *Gofur) Head(path string, h Handler, plugs ...Plug) {
	gf.add("HEAD", path, h, plugs...)
}

func (gf *Gofur) Options(path string, h Handler, plugs ...Plug) {
	gf.add("OPTIONS", path, h, plugs...)
}

func (gf *Gofur) makeHTTPRouterHandle(h Handler, plugs ...Plug) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := newContext(w, r, params)
		for i := len(plugs) - 1; i >= 0; i-- {
			h = plugs[i](h)
		}
		for i := len(gf.plugs) - 1; i >= 0; i-- {
			h = gf.plugs[i](h)
		}
		if err := h(ctx); err != nil {
			// todo: handle the error from the error handler huh?
			gf.ErrorHandler(err, ctx)
		}
	}
}

func defaultErrorHandler(err error, c *Context) error {
	slog.Error("error", "err", err)
	return nil
}
