package app

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/omissis/kube-apiserver-proxy/pkg/kube"
	"github.com/omissis/kube-apiserver-proxy/pkg/kube/proxy"
	"k8s.io/client-go/rest"
)

var ErrCannotCreateContainer = fmt.Errorf("cannot create container")

const (
	apiServerPort = 8080
)

type ContainerFactoryFunc func() (*Container, error)

func NewDefaultParameters() Parameters {
	return Parameters{
		APIServerHost:     "0.0.0.0",
		APIServerPort:     apiServerPort,
		APIAllowedOrigins: []string{"http://localhost:3000", "http://kasp.dev"},
	}
}

type Parameters struct {
	APIServerHost     string
	APIServerPort     uint16
	APIAllowedOrigins []string
	KubeconfigPath    string
}

type services struct {
	echo                 *echo.Echo
	k8sRESTClientFactory *kube.DefaultK8sRESTClientFactory
	k8sHTTProxy          *proxy.HTTP
	k8sHttpClient        *http.Client
	k8sRESTConfigFactory *kube.DefaultRESTConfigFactory
}

func NewContainer() *Container {
	return &Container{
		Parameters: NewDefaultParameters(),
	}
}

type Container struct {
	Parameters
	services
}

func (c *Container) Echo() *echo.Echo {
	if c.echo == nil {
		c.echo = echo.New()
		c.echo.Debug = true
		if len(c.Parameters.APIAllowedOrigins) > 0 {
			c.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
				AllowOrigins: c.Parameters.APIAllowedOrigins,
				AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			}))
		}
	}

	return c.echo
}

func (c *Container) K8sHTTPProxy() *proxy.HTTP {
	if c.k8sHTTProxy == nil {
		c.k8sHTTProxy = proxy.NewHTTP(c.K8sRESTClientFactory(), []proxy.ResponseBodyTransformer{
			proxy.NewJqResponseBodyTransformer(),
		})
	}

	return c.k8sHTTProxy
}

func (c *Container) K8sRESTClientFactory() *kube.DefaultK8sRESTClientFactory {
	if c.k8sRESTClientFactory == nil {
		c.k8sRESTClientFactory = kube.NewDefaultK8sRESTClientFactory(
			c.K8sRESTConfigFactory(),
			c.K8sHttpClient(),
			c.KubeconfigPath,
		)
	}

	return c.k8sRESTClientFactory
}

func (c *Container) K8sRESTConfigFactory() *kube.DefaultRESTConfigFactory {
	if c.k8sRESTConfigFactory == nil {
		c.k8sRESTConfigFactory = kube.NewDefaultRESTConfigFactory()
	}

	return c.k8sRESTConfigFactory
}

func (c *Container) K8sHttpClient() *http.Client {
	if c.k8sHttpClient == nil {
		k8sHttpClient, err := rest.HTTPClientFor(c.K8sRestConfig())
		if err != nil {
			panic(err)
		}

		c.k8sHttpClient = k8sHttpClient
	}

	return c.k8sHttpClient
}

func (c *Container) K8sRestConfig() *rest.Config {
	config, err := c.K8sRESTConfigFactory().New(c.KubeconfigPath)
	if err != nil {
		panic(err)
	}

	return config
}
