//go:build unit

package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/omissis/kube-apiserver-proxy/internal/app"
)

func TestHTTPServeMux(t *testing.T) {
	t.Parallel()

	container := app.NewContainer()

	assert.NotNil(t, container.HTTPServeMux())
}

func TestHTTPServer(t *testing.T) {
	t.Parallel()

	container := app.NewContainer()

	assert.NotNil(t, container.HTTPServer())
}

// TODO: enable this test
// func TestK8sHTTPProxy(t *testing.T) {
// 	t.Parallel()

// 	container := app.NewContainer()

// 	assert.NotNil(t, container.K8sHTTPProxy())
// }

// TODO: enable this test
// func TestRESTClientFactory(t *testing.T) {
// 	t.Parallel()

// 	container := app.NewContainer()

// 	assert.NotNil(t, container.RESTClientFactory())
// }
