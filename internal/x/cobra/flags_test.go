package cobra_test

import (
	"testing"

	"github.com/omissis/kube-apiserver-proxy/internal/x/cobra"
)

func TestInitEnvs(t *testing.T) {
	t.Parallel()

	v := cobra.InitEnvs("test")

	if v == nil {
		t.Error("InitEnvs() returned nil")
	}
}
