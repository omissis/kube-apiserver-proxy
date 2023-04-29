package app

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var ErrCannotCreateContainer = fmt.Errorf("cannot create container")

type ContainerFactoryFunc func() (*Container, error)

func NewDefaultParameters() Parameters {
	return Parameters{}
}

type Parameters struct{}

type services struct {
	echo *echo.Echo
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
		c.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"localhost", "kube-apiserver-proxy.dev"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}))
	}

	return c.echo
}
