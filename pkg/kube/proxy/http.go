package proxy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/omissis/kube-apiserver-proxy/pkg/kube"
)

var (
	ErrCannotApplyResponseTransformers = errors.New("cannot apply response transformers")
	ErrCannotCreateRESTClient          = errors.New("cannot create rest client")
	ErrCannotGetProxiedResponseBody    = errors.New("cannot get response of proxied response body")
	ErrCannotParseRequestURI           = errors.New("cannot parse request URI")
	ErrCannotTransformResponseBody     = errors.New("cannot transform response body")
	ErrCannotWriteResponseBody         = errors.New("cannot write body to the response")
	ErrContextIsNil                    = errors.New("context is nil")
	ErrResponseWriterIsNil             = errors.New("response writer is nil")
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

	if w == nil {
		log.Printf("error: no response writer was passed\n")

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
	if ctx == nil {
		return ErrContextIsNil
	}

	if w == nil {
		return ErrResponseWriterIsNil
	}

	req, err := h.restClientFactory.Request(r)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreateRESTClient, err)
	}

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	res := req.Do(cctx)

	body, err := res.Raw()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotGetProxiedResponseBody, err)
	}

	body, err = h.applyTransformers(r, body)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotApplyResponseTransformers, err)
	}

	sc := 0
	res.StatusCode(&sc)
	w.WriteHeader(sc)

	_, err = w.Write(body)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotWriteResponseBody, err)
	}

	return nil
}

func (h *HTTP) applyTransformers(r http.Request, body []byte) ([]byte, error) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		return body, fmt.Errorf("%w: %w", ErrCannotParseRequestURI, err)
	}

	for _, rt := range h.responseTransformers {
		src := u.Query().Get(rt.Name())

		if src != "" {
			body, err = rt.Run(body, map[string]any{"src": src})
			if err != nil {
				return body, fmt.Errorf("%w: %w", ErrCannotTransformResponseBody, err)
			}

			return body, nil
		}
	}

	return body, nil
}
