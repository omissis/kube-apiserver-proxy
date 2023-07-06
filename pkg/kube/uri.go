package kube

import (
	"fmt"
	"strings"
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

		if len(parts) < 3 {
			return "", "", fmt.Errorf("uri '%s' has less than 3 parts in it", uri)
		}

		return parts[1], parts[2], nil
	}

	return "", "", fmt.Errorf("uri '%s' is not supported", uri)
}
