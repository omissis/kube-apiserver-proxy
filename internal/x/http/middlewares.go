package http

import (
	"net/http"
	"strings"
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

type CORSConfig struct {
	AllowMethods     []string
	AllowOrigins     []string
	AllowHeaders     []string
	AllowCredentials bool
}

func CORSMuxMiddleware(conf CORSConfig) MuxMiddleware {
	return func(next http.Handler) http.Handler {
		return CORSMiddleware(next, conf)
	}
}

func CORSMiddleware(next http.Handler, conf CORSConfig) Middleware {
	methods := []string{"*"}
	if conf.AllowMethods != nil {
		methods = conf.AllowMethods
	}

	origins := []string{"*"}
	if conf.AllowOrigins != nil {
		origins = conf.AllowOrigins
	}

	headers := []string{"Origin", "Content-Type", "Accept"}
	if conf.AllowHeaders != nil {
		headers = conf.AllowHeaders
	}

	credentials := "false"
	if conf.AllowCredentials {
		credentials = "true"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
		w.Header().Set("Access-Control-Allow-Origin", strings.Join(origins, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ", "))
		w.Header().Set("Access-Control-Allow-Credentials", credentials)

		next.ServeHTTP(w, r)
	})
}
