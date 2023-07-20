package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"

	"github.com/omissis/kube-apiserver-proxy/internal/app"
)

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

			ctr.HTTPServeMux().HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				if r == nil {
					http.Error(w, "request is empty", http.StatusInternalServerError)

					return
				}

				if w == nil {
					http.Error(w, "response is empty", http.StatusInternalServerError)

					return
				}

				ctx, cancel := context.WithCancel(r.Context())
				defer cancel()

				if err := ctr.K8sHTTPProxy().DoServeHTTP(ctx, w, *r); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			})

			return ctr.HTTPServer().ListenAndServe()
		},
	}

	setupServeCommandFlags(cmd)

	return cmd
}

func setupServeCommandFlags(cmd *cobra.Command) {
	kubeconfigDefault := ""

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
