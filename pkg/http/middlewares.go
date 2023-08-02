package http

import (
	"fmt"
	"net/http"
	"net/url"
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", corsMethod(conf))
		w.Header().Set("Access-Control-Allow-Origin", corsOrigin(conf, r))
		w.Header().Set("Access-Control-Allow-Headers", corsHeaders(conf))
		w.Header().Set("Access-Control-Allow-Credentials", corsCredentials(conf))

		next.ServeHTTP(w, r)
	})
}

func corsMethod(conf CORSConfig) string {
	methods := []string{"*"}
	if conf.AllowMethods != nil {
		methods = conf.AllowMethods
	}

	return strings.Join(methods, ", ")
}

//nolint:cyclop // leave this alone
func corsOrigin(conf CORSConfig, req *http.Request) string {
	if req == nil {
		return ""
	}

	origins := []string{"*"}
	if conf.AllowOrigins != nil {
		origins = conf.AllowOrigins
	}

	if len(origins) == 1 && origins[0] == "*" {
		return origins[0]
	}

	origin := req.Header.Get("Origin")

	u, err := url.Parse(origin)
	if err != nil {
		return ""
	}

	for _, o := range origins {
		if o == origin {
			return origin
		}

		v, err := url.Parse(o)
		if err != nil {
			return ""
		}

		if v.Scheme == u.Scheme && v.Host == u.Host {
			return fmt.Sprintf("%s://%s", v.Scheme, v.Host)
		}
	}

	return ""
}

func corsHeaders(conf CORSConfig) string {
	headers := []string{"Origin", "Content-Type", "Accept"}
	if conf.AllowHeaders != nil {
		headers = conf.AllowHeaders
	}

	return strings.Join(headers, ", ")
}

func corsCredentials(conf CORSConfig) string {
	credentials := "false"
	if conf.AllowCredentials {
		credentials = "true"
	}

	return credentials
}
