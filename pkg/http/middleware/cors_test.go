package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/omissis/kube-apiserver-proxy/pkg/http/middleware"
)

func TestCORS(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc            string
		conf            middleware.CORSConfig
		origin          string
		wantOrigin      string
		methods         string
		wantMethods     string
		credentials     string
		wantCredentials string
	}{
		{
			desc:            "default config",
			conf:            middleware.CORSConfig{},
			origin:          "*",
			wantOrigin:      "*",
			methods:         "*",
			wantMethods:     "*",
			credentials:     "false",
			wantCredentials: "false",
		},
		{
			desc: "multiple origins",
			conf: middleware.CORSConfig{
				AllowOrigins: []string{"https://api.kube-apiserver-proxy.dev", "https://api.kube-apiserver-proxy.test"},
			},
			origin:          "https://api.kube-apiserver-proxy.dev",
			wantOrigin:      "https://api.kube-apiserver-proxy.dev",
			methods:         "*",
			wantMethods:     "*",
			credentials:     "false",
			wantCredentials: "false",
		},
		{
			desc: "origin with authentication",
			conf: middleware.CORSConfig{
				AllowOrigins: []string{"https://api.kube-apiserver-proxy.dev", "https://api.kube-apiserver-proxy.test"},
			},
			origin:          "https://foo:bar@api.kube-apiserver-proxy.dev",
			wantOrigin:      "https://api.kube-apiserver-proxy.dev",
			methods:         "*",
			wantMethods:     "*",
			credentials:     "false",
			wantCredentials: "false",
		},
		{
			desc: "multiple methods",
			conf: middleware.CORSConfig{
				AllowMethods: []string{http.MethodGet, "POST"},
			},
			origin:          "*",
			wantOrigin:      "*",
			methods:         "GET, POST",
			wantMethods:     "GET, POST",
			credentials:     "false",
			wantCredentials: "false",
		},
	}
	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			handler := middleware.CORS(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					io.WriteString(w, "OK")
				}),
				tC.conf,
			)

			url := "https://api.kube-apiserver-proxy.dev/api/v1/namespaces/kube-system/pods"

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Origin", tC.origin)

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			resp := w.Result()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()

			assert.Equal(t, "OK", string(body))
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, tC.wantCredentials, resp.Header.Get("Access-Control-Allow-Credentials"))
			assert.Equal(t, "Origin, Content-Type, Accept", resp.Header.Get("Access-Control-Allow-Headers"))
			assert.Equal(t, tC.wantOrigin, resp.Header.Get("Access-Control-Allow-Origin"))
			assert.Equal(t, tC.wantMethods, resp.Header.Get("Access-Control-Allow-Methods"))
		})
	}
}
