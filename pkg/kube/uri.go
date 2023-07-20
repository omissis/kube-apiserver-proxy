package kube

import (
	"fmt"
	"strings"
)

const minURIComponentsCount = 3

var (
	ErrMalformedURI      = fmt.Errorf("uri has less than 3 parts in it")
	ErrURIIsNotSupported = fmt.Errorf("uri is not supported")
)

func GetGroupVersionFromURI(uri string) (string, string, error) {
	if uri == "/api" {
		return "core", "", nil
	}

	if strings.HasPrefix(uri, "/api/") {
		return "core", "v1", nil
	}

	if uri == "/apis" {
		return "apis", "", nil
	}

	if strings.HasPrefix(uri, "/apis/") {
		parts := strings.Split(strings.Trim(uri, "/"), "/")

		if len(parts) < minURIComponentsCount {
			return "", "", fmt.Errorf("%w '%s'", ErrMalformedURI, uri)
		}

		return parts[1], parts[2], nil
	}

	return "", "", fmt.Errorf("%w: '%s'", ErrURIIsNotSupported, uri)
}
