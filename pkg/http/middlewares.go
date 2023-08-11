package http

import (
	"net/http"
)

func NewServeMux(mws []MuxMiddleware) *ServeMux {
	if mws == nil {
		mws = make([]MuxMiddleware, 0)
	}

	return &ServeMux{
		ServeMux:    http.NewServeMux(),
		middlewares: mws,
	}
}

type ServeMux struct {
	*http.ServeMux
	middlewares []MuxMiddleware
}

func (mux *ServeMux) Handle(pattern string, handler http.Handler) {
	for _, mw := range mux.middlewares {
		handler = mw(handler)
	}

	mux.ServeMux.Handle(pattern, handler)
}

func (mux *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	var hh http.Handler = http.HandlerFunc(handler)

	for _, mw := range mux.middlewares {
		hh = mw(hh)
	}

	mux.ServeMux.Handle(pattern, hh)
}

func (mux *ServeMux) Use(mw MuxMiddleware) {
	if mux.middlewares == nil {
		mux.middlewares = make([]MuxMiddleware, 0)
	}

	mux.middlewares = append(mux.middlewares, mw)
}

type MuxMiddleware func(http.Handler) http.Handler

type Middleware http.Handler
