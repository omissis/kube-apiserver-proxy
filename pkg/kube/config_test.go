package kube_test

import (
	"strings"
	"testing"

	"github.com/omissis/kube-apiserver-proxy/pkg/kube"
)

func TestDefaultRESTConfigFactory_New_WithoutCluster(t *testing.T) {
	f := kube.NewDefaultRESTConfigFactory()

	cfg, err := f.New("")

	if err == nil {
		t.Error("expected error, got nil")
	}

	prefix := "cannot build config from flags nor create it in cluster"

	if err != nil && !strings.HasPrefix(err.Error(), prefix) {
		t.Errorf("expected error message prefix: %v, got: %v", prefix, err.Error())
	}

	if cfg != nil {
		t.Errorf("did not expect a config, got one: %+v", cfg)
	}
}
