package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/omissis/kube-apiserver-proxy/internal/app"
)

func NewServeCommand(ctr *app.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run kube-apiserver-proxy server",
		Args:  cobra.ExactArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println("Running the kube-apiserver-proxy server...")

			return ctr.Echo().Start(":8080")
		},
	}

	return cmd
}
