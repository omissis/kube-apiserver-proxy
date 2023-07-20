package app

import (
	"fmt"
	"net/http"
	"time"

	"k8s.io/client-go/rest"

	httpx "github.com/omissis/kube-apiserver-proxy/pkg/http"
	"github.com/omissis/kube-apiserver-proxy/pkg/kube"
	"github.com/omissis/kube-apiserver-proxy/pkg/kube/proxy"
)

var ErrCannotCreateContainer = fmt.Errorf("cannot create container")

const (
	apiServerPort = 8080
)

type ContainerFactoryFunc func() (*Container, error)

func NewDefaultParameters() Parameters {
	return Parameters{
		APIServerTimeout:  5 * time.Second,
		APIServerHost:     "0.0.0.0",
		APIServerPort:     apiServerPort,
		APIAllowedOrigins: []string{"http://localhost:3000", "http://kasp.dev"},
	}
}

type Parameters struct {
	APIServerTimeout  time.Duration
	APIServerHost     string
	APIServerPort     uint16
	APIAllowedOrigins []string
	KubeconfigPath    string
}

type services struct {
	httpServer           *http.Server
	httpServeMux         *httpx.ServeMux
	k8sRESTClientFactory *kube.DefaultRESTClientFactory
	k8sHTTProxy          *proxy.HTTP
	k8sHTTPClient        *http.Client
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

func (c *Container) HTTPServeMux() *httpx.ServeMux {
	if c.httpServeMux == nil {
		c.httpServeMux = httpx.NewServeMux(nil)

		if len(c.Parameters.APIAllowedOrigins) > 0 {
			c.httpServeMux.Use(
				httpx.CORSMuxMiddleware(httpx.CORSConfig{
					AllowOrigins: c.Parameters.APIAllowedOrigins,
				}),
			)
		}
	}

	return c.httpServeMux
}

func (c *Container) HTTPServer() *http.Server {
	if c.httpServer == nil {
		c.httpServer = &http.Server{
			Addr:              fmt.Sprintf("%s:%d", c.APIServerHost, c.APIServerPort),
			Handler:           c.HTTPServeMux(),
			ReadHeaderTimeout: c.APIServerTimeout,
		}
	}

	return c.httpServer
}

func (c *Container) K8sHTTPProxy() *proxy.HTTP {
	if c.k8sHTTProxy == nil {
		c.k8sHTTProxy = proxy.NewHTTP(c.RESTClientFactory(), []proxy.ResponseBodyTransformer{
			proxy.NewJqResponseBodyTransformer(),
		})
	}

	return c.k8sHTTProxy
}

func (c *Container) RESTClientFactory() *kube.DefaultRESTClientFactory {
	if c.k8sRESTClientFactory == nil {
		c.k8sRESTClientFactory = kube.NewDefaultRESTClientFactory(
			c.RESTConfigFactory(),
			c.K8sHTTPClient(),
			c.KubeconfigPath,
		)
	}

	return c.k8sRESTClientFactory
}

func (c *Container) RESTConfigFactory() *kube.DefaultRESTConfigFactory {
	if c.k8sRESTConfigFactory == nil {
		c.k8sRESTConfigFactory = kube.NewDefaultRESTConfigFactory()
	}

	return c.k8sRESTConfigFactory
}

func (c *Container) K8sHTTPClient() *http.Client {
	if c.k8sHTTPClient == nil {
		k8sHTTPClient, err := rest.HTTPClientFor(c.K8sRestConfig())
		if err != nil {
			panic(err)
		}

		c.k8sHTTPClient = k8sHTTPClient
	}

	return c.k8sHTTPClient
}

func (c *Container) K8sRestConfig() *rest.Config {
	config, err := c.RESTConfigFactory().New(c.KubeconfigPath)
	if err != nil {
		panic(err)
	}

	return config
}
