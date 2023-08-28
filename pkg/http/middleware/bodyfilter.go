package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"golang.org/x/exp/slices"

	"github.com/omissis/kube-apiserver-proxy/pkg/config"
	kaspHttp "github.com/omissis/kube-apiserver-proxy/pkg/http"
)

var ErrDuringBodyFilter = fmt.Errorf("error during body filter")

func BodyFilterMux(conf []config.BodyFilterConfig) http2.MuxMiddleware {
	return func(next http.Handler) http.Handler {
		return BodyFilter(next, conf)
	}
}

func BodyFilter(next http.Handler, conf []config.BodyFilterConfig) http2.Middleware {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, match := MatchConfig(r, conf)
		if !match {
			next.ServeHTTP(w, r)

			return
		}

		if r.Body == nil {
			slog.Warning("empty request body")

			http.Error(w, "Empty request body", http.StatusBadRequest)

			return
		}

		bodyDecoder := json.NewDecoder(r.Body)

		filteredBody := map[string]any{}
		filterTarget := map[string]any{}

		if err := bodyDecoder.Decode(&filteredBody); err != nil {
			slog.Error("cannot decode body", "error", err, "body", r.Body)
		
			http.Error(w, "Decoding Error", http.StatusBadRequest)

			return
		}

		if err := json.Unmarshal([]byte(c.Filter), &filterTarget); err != nil {
			slog.Error("cannot unmarshal filter config", "error", err, "filterConfig", c.Filter)

			http.Error(w, "Unmarshalling Error", http.StatusInternalServerError)

			return
		}

		if err := FillTargetFromBody(filteredBody, filterTarget); err != nil {
			slog.Error("cannot create filtered body", "error", err, "filteredBody", filteredBody, "filterTarget", filterTarget)

			http.Error(w, "Mapping Error", http.StatusInternalServerError)

			return
		}

		bodyFromTarget, err := json.Marshal(filterTarget)
		if err != nil {
			slog.Error("cannot marshal filtered body", "error", err, "filterTarget", filterTarget)

			http.Error(w, "Marshalling Error", http.StatusInternalServerError)

			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyFromTarget))

		next.ServeHTTP(w, r)
	})
}

func MatchConfig(r *http.Request, conf []config.BodyFilterConfig) (config.BodyFilterConfig, bool) {
	for _, c := range conf {
		if !slices.Contains(c.Methods, r.Method) {
			continue
		}

		for _, p := range c.Paths {
			switch p.Type {
			case "glob":
				if m, err := filepath.Match(p.Path, r.URL.Path); m && err == nil {
					return c, true
				}

			default:
				return config.BodyFilterConfig{}, false
			}
		}
	}

	return config.BodyFilterConfig{}, false
}

func FillTargetFromBody(body, target map[string]any) error {
	var err error

	for k, v := range target {
		if _, ok := body[k]; !ok {
			err = fmt.Errorf("%w: key %s not found in base map", ErrDuringBodyFilter, k)

			break
		}

		if v == "*" {
			target[k] = body[k]

			continue
		}

		targetKMap, ok := target[k].(map[string]any)
		if ok {
			bodyKMap, bOk := body[k].(map[string]any)
			if !bOk {
				err = fmt.Errorf("%w: key %s is not a map", ErrDuringBodyFilter, k)

				break
			}

			return FillTargetFromBody(bodyKMap, targetKMap)
		}

		targetKArr, ok := target[k].([]any)
		if ok {
			bodyKArr, ok := body[k].([]any)
			if !ok {
				err = fmt.Errorf("%w: key %s is not an array", ErrDuringBodyFilter, k)

				break
			}

			return fillTargetArrayHelper(bodyKArr, targetKArr, k)
		}
	}

	return err
}

func fillTargetArrayHelper(body, target []any, key string) error {
	for i, v := range target {
		targetIMap, ok := v.(map[string]any)
		if ok {
			bodyIMap, ok := body[i].(map[string]any)
			if !ok {
				return fmt.Errorf("%w: key %s is not an array of maps", ErrDuringBodyFilter, key)
			}

			return FillTargetFromBody(bodyIMap, targetIMap)
		}
	}

	return nil
}
