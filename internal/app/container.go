package app

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/omissis/kube-apiserver-proxy/internal/graph"
	"github.com/omissis/kube-apiserver-proxy/internal/graph/generated"
)

var ErrCannotCreateContainer = fmt.Errorf("cannot create container")

type ContainerFactoryFunc func() (*Container, error)

func NewDefaultParameters() Parameters {
	return Parameters{}
}

type Parameters struct{}

type services struct {
	gqlServer            *handler.Server
	gqlPlaygroundHandler http.HandlerFunc
	clientset            *kubernetes.Clientset
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

func (c *Container) GQLServerHandler() *handler.Server {
	if c.gqlServer == nil {
		c.gqlServer = handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
			Resolvers: &graph.Resolver{
				Clientset: c.Clientset(),
			},
		}))
	}

	return c.gqlServer
}

func (c *Container) GQLPlaygroundHandler() http.HandlerFunc {
	if c.gqlPlaygroundHandler == nil {
		c.gqlPlaygroundHandler = playground.Handler("GraphQL playground", "/query")
	}

	return c.gqlPlaygroundHandler
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
