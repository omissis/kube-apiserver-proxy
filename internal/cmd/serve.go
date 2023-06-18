package cmd

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/spf13/cobra"

	"github.com/omissis/kube-apiserver-proxy/internal/app"
	"github.com/omissis/kube-apiserver-proxy/internal/kube"
)

const kubeApiServer = "kubernetes.default.svc"

func NewServeCommand(ctr *app.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run kube-apiserver-proxy server",
		Args:  cobra.ExactArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println("Running the kube-apiserver-proxy server...")

			ctr.Echo().Any("/*", func(c echo.Context) error {
				group, version, err := kube.GetGroupVersionFromURI(c.Request().RequestURI)
				if err != nil {
					return fmt.Errorf("cannot get group and version from request uri: %w", err)
				}

				return kube.EchoProxy(ctr.K8sRESTClient(group, version), c)
			})

			return ctr.Echo().Start(fmt.Sprintf("%s:%d", ctr.APIServerHost, ctr.APIServerPort))
		},
	}

	return cmd
}
