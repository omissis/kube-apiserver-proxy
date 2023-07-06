package app_test

import (
	"testing"

	"github.com/omissis/kube-apiserver-proxy/internal/app"
	"github.com/stretchr/testify/assert"
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

func TestK8sRESTClientFactory(t *testing.T) {
	t.Parallel()

	container := app.NewContainer()

	assert.NotNil(t, container.K8sRESTClientFactory())
}
