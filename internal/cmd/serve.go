package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"

	"github.com/omissis/kube-apiserver-proxy/internal/app"
)

const kubeApiServer = "kubernetes.default.svc"

var ErrParsingFlag = errors.New("cannot parse command-line flag")

type ServeCommandFlags struct {
	Kubeconfig string
}

func NewServeCommand(ctr *app.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run kube-apiserver-proxy server",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Println("Running the kube-apiserver-proxy server...")

			flags, err := getServeCommandFlags(cmd)
			if err != nil {
				return err
			}

			ctr.KubeconfigPath = flags.Kubeconfig

			ctr.Echo().Any("/*", func(c echo.Context) error {
				return ctr.K8sHTTPProxy().DoServeHTTP(c.Response().Writer, *c.Request())
			})

			return ctr.Echo().Start(fmt.Sprintf("%s:%d", ctr.APIServerHost, ctr.APIServerPort))
		},
	}

	setupServeCommandFlags(cmd)

	return cmd
}

func setupServeCommandFlags(cmd *cobra.Command) {
	var kubeconfigDefault string = ""

	if home := homedir.HomeDir(); home != "" {
		kubeconfigDefault = filepath.Join(home, ".kube", "config")
	}

	cmd.Flags().String(
		"kubeconfig",
		kubeconfigDefault,
		"(optional) absolute path to the kubeconfig file",
	)
}

func getServeCommandFlags(cmd *cobra.Command) (ServeCommandFlags, error) {
	kubeconfig, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		return ServeCommandFlags{}, fmt.Errorf("%w '%s': %w", ErrParsingFlag, "kubeconfig", err)
	}

	return ServeCommandFlags{
		Kubeconfig: kubeconfig,
	}, nil
}
