package app

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/omissis/kube-apiserver-proxy/internal/graph"
)

var ErrCannotCreateContainer = fmt.Errorf("cannot create container")

type ContainerFactoryFunc func() (*Container, error)

func NewDefaultParameters() Parameters {
	return Parameters{}
}

type Parameters struct{}

type services struct {
	echo                 *echo.Echo
	gqlServer            *handler.Server
	gqlPlaygroundHandler http.HandlerFunc
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
		c.gqlServer = handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	}

	return c.gqlServer
}

func (c *Container) GQLPlaygroundHandler() http.HandlerFunc {
	if c.gqlPlaygroundHandler == nil {
		c.gqlPlaygroundHandler = playground.Handler("GraphQL playground", "/query")
	}

	return c.gqlPlaygroundHandler
}

func (c *Container) Echo() *echo.Echo {
	if c.echo == nil {
		c.echo = echo.New()
		c.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"localhost", "kube-apiserver-proxy.dev"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}))
	}

	return c.echo
}
