package cmd

import (
	"fmt"
	"log"
	"net/http"

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

			http.Handle("/", ctr.GQLPlaygroundHandler())
			http.Handle("/query", ctr.GQLServerHandler())

			log.Printf("connect to http://localhost:%s/ for GraphQL playground", "8080")

			return http.ListenAndServe(":8080", nil)
		},
	}

	return cmd
}
