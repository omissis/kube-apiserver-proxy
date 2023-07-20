//go:build unit

package proxy_test

import (
	"testing"

	"github.com/omissis/kube-apiserver-proxy/pkg/kube/proxy"
)

func TestJqResponseBodyTransformer_Name(t *testing.T) {
	tf := proxy.NewJqResponseBodyTransformer()

	if got, want := tf.Name(), "jq"; got != want {
		t.Errorf("tf.Name() = %s, want %s", got, want)
	}
}

func TestJqResponseBodyTransformer_Run(t *testing.T) {
	t.Parallel()

	tf := proxy.NewJqResponseBodyTransformer()

	testCases := []struct {
		desc    string
		src     string
		want    []byte
		body    []byte
		wantErr bool
	}{
		{
			desc:    "no output transformation",
			src:     ".",
			want:    []byte(`{"foo":"bar"}`),
			body:    []byte(`{"foo":"bar"}`),
			wantErr: false,
		},
		{
			desc:    "simple output transformation",
			src:     ".foo",
			want:    []byte(`"bar"`),
			body:    []byte(`{"foo":"bar"}`),
			wantErr: false,
		},
		{
			desc:    "nested output transformation",
			src:     ".foo.bar",
			want:    []byte(`"baz"`),
			body:    []byte(`{"foo":{"bar":"baz"}}`),
			wantErr: false,
		},
	}
	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			res, err := tf.Run(tC.body, map[string]any{"src": tC.src})
			if (err != nil) != tC.wantErr {
				t.Errorf("wanted error %v, got %v", tC.wantErr, err)
			}

			if string(res) != string(tC.want) {
				t.Errorf("wanted %s, got %s", tC.want, res)
			}
		})
	}
}
