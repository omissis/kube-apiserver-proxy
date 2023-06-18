package kube

import (
	"fmt"
	"strings"
)

func GetGroupVersionFromURI(uri string) (string, string, error) {
	parts := strings.Split(uri, "/")

	if len(parts) < 3 {
		return "", "", fmt.Errorf("uri '%s' has less than 3 parts in it", uri)
	}

	return parts[1], parts[2], nil
}
