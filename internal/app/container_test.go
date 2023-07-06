package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/omissis/kube-apiserver-proxy/internal/app"
)

func TestEcho(t *testing.T) {
	t.Parallel()

	container := app.NewContainer()

	assert.NotNil(t, container.Echo())
}

func TestK8sHTTPProxy(t *testing.T) {
	t.Parallel()

	container := app.NewContainer()

	assert.NotNil(t, container.K8sHTTPProxy())
}

func TestRESTClientFactory(t *testing.T) {
	t.Parallel()

	container := app.NewContainer()

	assert.NotNil(t, container.RESTClientFactory())
}
