package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	kaspHttp "github.com/omissis/kube-apiserver-proxy/pkg/http"
)

type CORSConfig struct {
	AllowMethods     []string
	AllowOrigins     []string
	AllowHeaders     []string
	AllowCredentials bool
}

func CORSMux(conf CORSConfig) http2.MuxMiddleware {
	return func(next http.Handler) http.Handler {
		return CORS(next, conf)
	}
}

func CORS(next http.Handler, conf CORSConfig) http2.Middleware {
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
			break
		}

		v, err := url.Parse(o)
		if err != nil {
			origin = ""

			break
		}

		if v.Scheme == u.Scheme && v.Host == u.Host {
			origin = fmt.Sprintf("%s://%s", v.Scheme, v.Host)

			break
		}
	}

	return origin
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
