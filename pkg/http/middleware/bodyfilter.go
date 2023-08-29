package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"

	"github.com/omissis/kube-apiserver-proxy/pkg/config"
	kaspHttp "github.com/omissis/kube-apiserver-proxy/pkg/http"
)

var ErrDuringBodyFilter = errors.New("error during body filter")

func BodyFilterMux(conf []config.BodyFilterConfig) kaspHttp.MuxMiddleware {
	return func(next http.Handler) http.Handler {
		return BodyFilter(next, conf)
	}
}

func BodyFilter(next http.Handler, conf []config.BodyFilterConfig) kaspHttp.Middleware {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			slog.Warn("empty request")

			http.Error(w, "Empty request", http.StatusBadRequest)

			return
		}

		c, match := MatchConfig(r, conf)
		if !match {
			next.ServeHTTP(w, r)

			return
		}

		if r.Body == nil {
			slog.Warn("empty request body")

			http.Error(w, "Empty request body", http.StatusBadRequest)

			return
		}

		body, err := decodeBody(r.Body)
		if err != nil {
			slog.Error("cannot decode body", "error", err, "body", r.Body)

			http.Error(w, "Decoding Error", http.StatusBadRequest)

			return
		}

		filteredBody, err := getFilteredBody(body, c.Filter)
		if err != nil {
			slog.Error("cannot get filtered body", "error", err, "body", body, "filter", c.Filter)

			statusCode := http.StatusInternalServerError

			if errors.Is(err, ErrDuringBodyFilter) {
				statusCode = http.StatusBadRequest
			}

			http.Error(w, "Filter Error", statusCode)

			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(filteredBody))

		next.ServeHTTP(w, r)
	})
}

func MatchConfig(r *http.Request, conf []config.BodyFilterConfig) (config.BodyFilterConfig, bool) {
	for _, c := range conf {
		if !slices.Contains(c.Methods, strings.ToUpper(r.Method)) {
			continue
		}

		for _, p := range c.Paths {
			switch p.Type {
			case "glob":
				if m, err := filepath.Match(p.Path, r.URL.Path); m && err == nil {
					return c, true
				}

			default:
				slog.Warn("unknown path type", "type", p.Type)

				return config.BodyFilterConfig{}, false
			}
		}
	}

	return config.BodyFilterConfig{}, false
}

func Filter(body, filteredBody map[string]any) error {
	for k, v := range filteredBody {
		if _, ok := body[k]; !ok {
			return fmt.Errorf("%w: key %s not found in base map", ErrDuringBodyFilter, k)
		}

		if v == "*" {
			filteredBody[k] = body[k]

			continue
		}

		targetKMap, ok := filteredBody[k].(map[string]any)
		if ok {
			bodyKMap, bOk := body[k].(map[string]any)
			if !bOk {
				return fmt.Errorf("%w: key %s is not a map", ErrDuringBodyFilter, k)
			}

			if err := Filter(bodyKMap, targetKMap); err != nil {
				return err
			}

			continue
		}

		targetKArr, ok := filteredBody[k].([]any)
		if ok {
			bodyKArr, ok := body[k].([]any)
			if !ok {
				return fmt.Errorf("%w: key %s is not an array", ErrDuringBodyFilter, k)
			}

			if err := filterArrayHelper(bodyKArr, targetKArr, k); err != nil {
				return err
			}

			continue
		}

		slog.Debug("key is not a map or an array, skipped", "key", k)
	}

	return nil
}

func filterArrayHelper(body, target []any, key string) error {
	if len(target) > len(body) {
		return fmt.Errorf("%w: key %s target array is bigger than body array", ErrDuringBodyFilter, key)
	}

	for i := range target {
		targetIMap, ok := target[i].(map[string]any)
		if ok {
			bodyIMap, bOk := body[i].(map[string]any)
			if !bOk {
				return fmt.Errorf("%w: key %s is not an array of maps", ErrDuringBodyFilter, key)
			}

			if err := Filter(bodyIMap, targetIMap); err != nil {
				return err
			}

			continue
		}

		targetIArr, ok := target[i].([]any)
		if ok {
			bodyIArr, ok := body[i].([]any)
			if !ok {
				return fmt.Errorf("%w: key %s is not an array of arrays", ErrDuringBodyFilter, key)
			}

			if err := filterArrayHelper(bodyIArr, targetIArr, key); err != nil {
				return err
			}

			continue
		}

		slog.Debug("key is not a map or an array, skipped", "key", key)
	}

	return nil
}

func decodeBody(body io.ReadCloser) (map[string]any, error) {
	filteredBody := map[string]any{}

	bodyDecoder := json.NewDecoder(body)

	if err := bodyDecoder.Decode(&filteredBody); err != nil {
		return nil, err
	}

	return filteredBody, nil
}

func getFilteredBody(body map[string]any, filter string) ([]byte, error) {
	filteredBody := map[string]any{}

	if err := json.Unmarshal([]byte(filter), &filteredBody); err != nil {
		return nil, err
	}

	if err := Filter(body, filteredBody); err != nil {
		return nil, err
	}

	bodyFromTarget, err := json.Marshal(filteredBody)
	if err != nil {
		return nil, err
	}

	return bodyFromTarget, nil
}
