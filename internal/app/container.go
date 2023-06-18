package app

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/scheme"
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
}

type services struct {
	clientset      *kubernetes.Clientset
	echo           *echo.Echo
	k8sRESTClients map[string]map[string]*rest.RESTClient
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

func (c *Container) K8sRESTClient(group, version string) *rest.RESTClient {
	if c.k8sRESTClients == nil {
		c.k8sRESTClients = make(map[string]map[string]*rest.RESTClient, 1)
	}

	restClient, ok := c.k8sRESTClients[group][version]
	if !ok {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}

		config.GroupVersion = &schema.GroupVersion{
			Group:   "apps",
			Version: "v1",
		}

		config.APIPath = "/apis"
		config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

		restClient, err = rest.RESTClientFor(config)
		if err != nil {
			panic(err.Error())
		}
	}

	return restClient
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

func (c *Container) Clientset() *kubernetes.Clientset {
	if c.clientset == nil {
		// creates the in-cluster config
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}

		// creates the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		c.clientset = clientset
	}

	return c.clientset
}
