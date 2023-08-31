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
			slog.Debug("body filter did not match request", "config paths", c.Paths, "path", r.URL.Path)

			next.ServeHTTP(w, r)

			return
		}

		slog.Debug("body filter matched request", "config paths", c.Paths, "path", r.URL.Path)

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

			case "prefix":
				if strings.HasPrefix(r.URL.Path, p.Path) {
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

		match, err := filterHandlerByKey(body, filteredBody, Filter, k)
		if err != nil {
			return err
		}

		if match {
			continue
		}

		match, err = filterHandlerByKey(
			body,
			filteredBody,
			func(b, f []any) error { return filterArrayHelper(b, f, k) },
			k,
		)
		if err != nil {
			return err
		}

		if match {
			continue
		}

		slog.Debug("key is not a map or an array, skipped", "key", k)
	}

	return nil
}

func filterArrayHelper(body, filteredBody []any, key string) error {
	if len(filteredBody) > len(body) {
		return fmt.Errorf("%w: key %s filteredBody array is bigger than body array", ErrDuringBodyFilter, key)
	}

	for i := range filteredBody {
		match, err := filterHandlerByIndex(body, filteredBody, Filter, i)
		if err != nil {
			return err
		}

		if match {
			continue
		}

		match, err = filterHandlerByIndex(
			body,
			filteredBody,
			func(b, f []any) error { return filterArrayHelper(b, f, key) },
			i,
		)
		if err != nil {
			return err
		}

		if match {
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

func filterHandlerByKey[T any](body, filteredBody map[string]any, f func(T, T) error, k string) (bool, error) {
	match := false

	targetKMap, ok := filteredBody[k].(T)
	if ok {
		bodyKMap, bOk := body[k].(T)
		if !bOk {
			return false, fmt.Errorf("%w: key %s type mismatch", ErrDuringBodyFilter, k)
		}

		if err := f(bodyKMap, targetKMap); err != nil {
			return false, err
		}

		match = true
	}

	return match, nil
}

func filterHandlerByIndex[T any](body, filteredBody []any, f func(T, T) error, i int) (bool, error) {
	match := false

	targetKMap, ok := filteredBody[i].(T)
	if ok {
		bodyKMap, bOk := body[i].(T)
		if !bOk {
			return false, fmt.Errorf("%w: index %d type mismatch", ErrDuringBodyFilter, i)
		}

		if err := f(bodyKMap, targetKMap); err != nil {
			return false, err
		}

		match = true
	}

	return match, nil
}
