package kube

import (
	"context"
	"fmt"

	"github.com/labstack/echo"
	"k8s.io/client-go/rest"
)

func EchoProxy(rc *rest.RESTClient, c echo.Context) error {
	req := rest.NewRequest(rc)
	req = req.
		Verb(c.Request().Method).
		RequestURI(c.Request().RequestURI).
		Body(c.Request().Body)

	res := req.Do(context.Background())

	res.StatusCode(&c.Response().Status)
	body, err := res.Raw()
	if err != nil {
		return fmt.Errorf("cannot get response of proxied response body: %w", err)
	}

	_, err = c.Response().Write(body)
	if err != nil {
		return fmt.Errorf("cannot write body to the response: %w", err)
	}

	return nil
}
