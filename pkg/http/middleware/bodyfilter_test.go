package middleware_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/omissis/kube-apiserver-proxy/pkg/config"
	"github.com/omissis/kube-apiserver-proxy/pkg/http/middleware"
)

func TestBodyFilter(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc           string
		conf           []config.BodyFilterConfig
		httpMethod     string
		path           string
		body           string
		wantBody       string
		wantStatusCode int
	}{
		{
			desc:           "default config",
			conf:           []config.BodyFilterConfig{},
			httpMethod:     "PATCH",
			path:           "/api/v1/namespaces/default/pods",
			body:           `{"metadata":{"name":"foo"}}`,
			wantBody:       `{"metadata":{"name":"foo"}}`,
			wantStatusCode: http.StatusOK,
		},
		{
			desc: "match method, not path",
			conf: []config.BodyFilterConfig{
				{
					Methods: []string{"PATCH"},
					Paths: []config.BodyFilterConfigPaths{
						{
							Path: "/api/v1/namespaces/default/deployments",
							Type: "glob",
						},
					},
				},
			},
			httpMethod:     "PATCH",
			path:           "/api/v1/namespaces/default/pods",
			body:           `{"metadata":{"name":"foo"}}`,
			wantBody:       `{"metadata":{"name":"foo"}}`,
			wantStatusCode: http.StatusOK,
		},
		{
			desc: "match method, match path -- success",
			conf: []config.BodyFilterConfig{
				{
					Methods: []string{"PATCH"},
					Paths: []config.BodyFilterConfigPaths{
						{
							Path: "/api/v1/namespaces/default/pods",
							Type: "glob",
						},
					},
					Filter: `{"metadata":{"name":"*"}}`,
				},
			},
			httpMethod:     "PATCH",
			path:           "/api/v1/namespaces/default/pods",
			body:           `{"metadata":{"name":"foo"}}`,
			wantBody:       `{"metadata":{"name":"foo"}}`,
			wantStatusCode: http.StatusOK,
		},
		{
			desc: "match method, match path, slightly more complex -- success",
			conf: []config.BodyFilterConfig{
				{
					Methods: []string{"PATCH"},
					Paths: []config.BodyFilterConfigPaths{
						{
							Path: "/apis/rbac.authorization.k8s.io/v1/*/*",
							Type: "glob",
						},
					},
					Filter: `{"metadata":{"labels":{"app": "*"}}}`,
				},
			},
			httpMethod:     "PATCH",
			path:           "/apis/rbac.authorization.k8s.io/v1/clusterrolebindings/cluster-admin",
			body:           `{"metadata":{"annotations":{"rbac.authorization.kubernetes.io/autoupdate":"true"},"creationTimestamp":"2023-08-10T17:19:25.000Z","labels":{"kubernetes.io/bootstrapping":"rbac-defaults","app":"true"},"name":"cluster-admin","resourceVersion":"732","uid":"c22faac5-d2f0-4447-8c15-e5edb13edeab"},"roleRef":{"apiGroup":"rbac.authorization.k8s.io","kind":"ClusterRole","name":"cluster-admin"},"subjects":[{"apiGroup":"rbac.authorization.k8s.io","kind":"Group","name":"system:masters"}]}`,
			wantBody:       `{"metadata":{"labels":{"app":"true"}}}`,
			wantStatusCode: http.StatusOK,
		},
		{
			desc: "match method, match path, body/filter mismatch -- failure",
			conf: []config.BodyFilterConfig{
				{
					Methods: []string{"PATCH"},
					Paths: []config.BodyFilterConfigPaths{
						{
							Path: "/api/v1/namespaces/default/pods",
							Type: "glob",
						},
					},
					Filter: `{"metadata":{"namae":"foo"}}`,
				},
			},
			httpMethod:     "PATCH",
			path:           "/api/v1/namespaces/default/pods",
			body:           `{"metadata":{"name":"foo"}}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			desc: "match method, match path, array -- success",
			conf: []config.BodyFilterConfig{
				{
					Methods: []string{"PATCH"},
					Paths: []config.BodyFilterConfigPaths{
						{
							Path: "/api/v1/namespaces/default/test",
							Type: "glob",
						},
					},
					Filter: `{"metadata":{"name":"*", "tests": [{"foo": "*"}, {"foo": "*"}]}}`,
				},
			},
			httpMethod:     "PATCH",
			path:           "/api/v1/namespaces/default/test",
			body:           `{"metadata":{"name":"foo", "tests": [{"foo": "bar"}, {"foo": "baz"}]}}`,
			wantBody:       `{"metadata":{"name":"foo","tests":[{"foo":"bar"},{"foo":"baz"}]}}`,
			wantStatusCode: http.StatusOK,
		},
		{
			desc: "match method, match path, array -- failure",
			conf: []config.BodyFilterConfig{
				{
					Methods: []string{"PATCH"},
					Paths: []config.BodyFilterConfigPaths{
						{
							Path: "/api/v1/namespaces/default/test",
							Type: "glob",
						},
					},
					Filter: `{"metadata":{"name":"*", "tests": [{"foo": "*"}, {"foo": "*"}]}}`,
				},
			},
			httpMethod:     "PATCH",
			path:           "/api/v1/namespaces/default/test",
			body:           `{"metadata":{"name":"foo", "tests": [{"foo": "bar"}]}}`,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			handler := middleware.BodyFilter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}

				defer r.Body.Close()

				w.Write(body)
			}), tC.conf)

			url := "https://api.kube-apiserver-proxy.dev" + tC.path

			req := httptest.NewRequest(tC.httpMethod, url, nil)

			req.Body = io.NopCloser(bytes.NewBufferString(tC.body))

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			resp := w.Result()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()

			assert.Equal(t, tC.wantStatusCode, resp.StatusCode)

			if tC.wantStatusCode == http.StatusOK {
				assert.Equal(t, tC.wantBody, string(body))
			}
		})
	}
}

func TestMatchConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc          string
		conf          []config.BodyFilterConfig
		httpMethod    string
		path          string
		body          string
		matchedConfig config.BodyFilterConfig
		matches       bool
	}{
		{
			desc: "does not match method -- does not match",
			conf: []config.BodyFilterConfig{
				{
					Methods: []string{"PATCH"},
					Paths: []config.BodyFilterConfigPaths{
						{
							Path: "/api/v1/namespaces/default/pods",
							Type: "glob",
						},
					},
				},
			},
			httpMethod: "POST",
			path:       "/api/v1/namespaces/default/pods",
			body:       `{"metadata":{"name":"foo"}}`,
			matches:    false,
		},
		{
			desc: "match method, does not match path -- does not match",
			conf: []config.BodyFilterConfig{
				{
					Methods: []string{"PATCH"},
					Paths: []config.BodyFilterConfigPaths{
						{
							Path: "/api/v1/namespaces/default/pods",
							Type: "glob",
						},
					},
					Filter: `{"metadata":{"name":"*"}}`,
				},
			},
			httpMethod: "PATCH",
			path:       "/api/v1/namespaces/default/deployments",
			body:       `{"metadata":{"name":"foo"}}`,
			matches:    false,
		},
		{
			desc: "match method, match path -- matches",
			conf: []config.BodyFilterConfig{
				{
					Methods: []string{"PATCH"},
					Paths: []config.BodyFilterConfigPaths{
						{
							Path: "/api/v1/namespaces/default/deployments",
							Type: "glob",
						},
					},
					Filter: `{"metadata":{"name":"*"}}`,
				},
				{
					Methods: []string{"PATCH"},
					Paths: []config.BodyFilterConfigPaths{
						{
							Path: "/api/v1/namespaces/default/pods",
							Type: "glob",
						},
					},
					Filter: `{"metadata":{"name":"*"}}`,
				},
			},
			httpMethod: "PATCH",
			path:       "/api/v1/namespaces/default/pods",
			body:       `{"metadata":{"name":"foo"}}`,
			matchedConfig: config.BodyFilterConfig{
				Methods: []string{"PATCH"},
				Paths: []config.BodyFilterConfigPaths{
					{
						Path: "/api/v1/namespaces/default/pods",
						Type: "glob",
					},
				},
				Filter: `{"metadata":{"name":"*"}}`,
			},
			matches: true,
		},
	}

	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			url := "https://api.kube-apiserver-proxy.dev" + tC.path

			req := httptest.NewRequest(tC.httpMethod, url, nil)

			req.Body = io.NopCloser(bytes.NewBufferString(tC.body))

			matchedConfig, matches := middleware.MatchConfig(req, tC.conf)

			assert.Equal(t, tC.matchedConfig, matchedConfig)
			assert.Equal(t, tC.matches, matches)
		})
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc    string
		body    map[string]any
		target  map[string]any
		wantErr bool
		err     string
	}{
		{
			desc: "fill target from body with mismatch - failure",
			body: map[string]any{
				"metadata": map[string]any{
					"labellos": map[string]any{
						"name": "foo",
					},
				},
			},
			target: map[string]any{
				"metadata": map[string]any{
					"labels": map[string]any{
						"name": "*",
					},
				},
			},
			wantErr: true,
			err:     "error during body filter: key labels not found in base map",
		},
		{
			desc: "fill target from body - success",
			body: map[string]any{
				"metadata": map[string]any{
					"name": "foo",
				},
			},
			target: map[string]any{
				"metadata": map[string]any{
					"name": "*",
				},
			},
			wantErr: false,
		},
		{
			desc: "fill target from body slightly more complex - success",
			body: map[string]any{
				"metadata": map[string]any{
					"name": "foo",
					"labels": map[string]any{
						"app": "bar",
					},
					"annotations": map[string]any{
						"app": "bar",
					},
				},
				"spec": map[string]any{
					"template": map[string]any{
						"metadata": map[string]any{
							"labels": map[string]any{
								"app": "bar",
							},
						},
					},
				},
			},
			target: map[string]any{
				"metadata": map[string]any{
					"labels": map[string]any{
						"app": "*",
					},
				},
			},
			wantErr: false,
		},
		{
			desc: "fill target from body array - success",
			body: map[string]any{
				"metadata": map[string]any{
					"test": []any{
						map[string]any{
							"foo": "bar",
						},
						map[string]any{
							"foo": "baz",
						},
						map[string]any{
							"foo": "qux",
						},
					},
				},
			},
			target: map[string]any{
				"metadata": map[string]any{
					"test": []any{
						map[string]any{
							"foo": "*",
						},
						map[string]any{
							"foo": "*",
						},
					},
				},
			},
		},
		{
			desc: "fill target from body array mismatch - failure",
			body: map[string]any{
				"metadata": map[string]any{
					"tests": []any{
						map[string]any{
							"foo": "bar",
						},
					},
				},
			},
			target: map[string]any{
				"metadata": map[string]any{
					"tests": []any{
						map[string]any{
							"foo": "*",
						},
						map[string]any{
							"foo": "*",
						},
					},
				},
			},
			wantErr: true,
			err:     "error during body filter: key tests target array is bigger than body array",
		},
		{
			desc: "fill target from body array mismatch keys - failure",
			body: map[string]any{
				"metadata": map[string]any{
					"tests": []any{
						map[string]any{
							"fooz": "bar",
						},
					},
				},
			},
			target: map[string]any{
				"metadata": map[string]any{
					"tests": []any{
						map[string]any{
							"foo": "*",
						},
					},
				},
			},
			wantErr: true,
			err:     "error during body filter: key foo not found in base map",
		},
		{
			desc: "fill target from body nested array - success",
			body: map[string]any{
				"metadata": map[string]any{
					"tests": []any{
						[]any{
							map[string]any{
								"foo": "bar",
							},
						},
					},
				},
			},
			target: map[string]any{
				"metadata": map[string]any{
					"tests": []any{
						[]any{
							map[string]any{
								"foo": "*",
							},
						},
					},
				},
			},
		},
		{
			desc: "fill target from body nested array - failure",
			body: map[string]any{
				"metadata": map[string]any{
					"tests": []any{
						[]any{
							map[string]any{
								"foo": "bar",
							},
						},
					},
				},
			},
			target: map[string]any{
				"metadata": map[string]any{
					"tests": []any{
						[]any{
							map[string]any{
								"foo": "*",
							},
							map[string]any{
								"foo": "*",
							},
						},
					},
				},
			},
			wantErr: true,
			err:     "error during body filter: key tests target array is bigger than body array",
		},
	}

	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			err := middleware.Filter(tC.body, tC.target)

			assert.Equal(t, tC.wantErr, err != nil)

			if err != nil {
				assert.Equal(t, tC.err, err.Error())
			}
		})
	}
}
