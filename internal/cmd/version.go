package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

func NewVersionCommand(versions map[string]string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display version information about kube-apiserver-proxy",
		Args:  cobra.ExactArgs(0),
		RunE: func(_ *cobra.Command, _ []string) error {
			ks := make([]string, 0, len(versions))
			for k := range versions {
				ks = append(ks, k)
			}
			sort.Strings(ks)

			for _, k := range ks {
				fmt.Printf("%s: %s\n", k, versions[k])
			}

			return nil
		},
	}
}
