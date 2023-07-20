package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/omissis/kube-apiserver-proxy/pkg/kube"
)

func NewHTTP(restClientFactory kube.RESTClientFactory, responseTransformers []ResponseBodyTransformer) *HTTP {
	return &HTTP{
		restClientFactory:    restClientFactory,
		responseTransformers: responseTransformers,
	}
}

type HTTP struct {
	responseTransformers []ResponseBodyTransformer
	restClientFactory    kube.RESTClientFactory
}

// ServeHTTP implements http.Handler interface
func (h *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		log.Printf("error: no request was passed\n")

		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	if err := h.DoServeHTTP(ctx, w, *r); err != nil {
		log.Printf("error: %s\n", err)
	}
}

// DoServeHTTP does the actual job of ServeHTTP, but it returns an error
//
// This method is useful when you want to integrate the handler with a different http server, and it helps
// to avoid the log.Printf in ServeHTTP, leaving the responsibility of the error handling to the caller.
func (h *HTTP) DoServeHTTP(ctx context.Context, w http.ResponseWriter, r http.Request) error {
	req, err := h.restClientFactory.Request(r)
	if err != nil {
		return fmt.Errorf("cannot create rest client: %w", err)
	}

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	res := req.Do(cctx)

	body, err := res.Raw()
	if err != nil {
		return fmt.Errorf("cannot get response of proxied response body: %w", err)
	}

	body, err = h.applyTransformers(r, body)
	if err != nil {
		return fmt.Errorf("cannot apply response transformers: %w", err)
	}

	sc := 0
	res.StatusCode(&sc)
	w.WriteHeader(sc)

	_, err = w.Write(body)
	if err != nil {
		return fmt.Errorf("cannot write body to the response: %w", err)
	}

	return nil
}

func (h *HTTP) applyTransformers(r http.Request, body []byte) ([]byte, error) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		return body, fmt.Errorf("cannot parse request uri: %w", err)
	}

	for _, rt := range h.responseTransformers {
		src := u.Query().Get(rt.Name())

		if src != "" {
			body, err = rt.Run(body, map[string]any{"src": src})
			if err != nil {
				return body, fmt.Errorf("cannot transform response body: %w", err)
			}

			return body, nil
		}
	}

	return body, nil
}
