package proxy_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/omissis/kube-apiserver-proxy/pkg/kube"
	"github.com/omissis/kube-apiserver-proxy/pkg/kube/proxy"
	"go.uber.org/mock/gomock"
	"k8s.io/client-go/rest"
)

func TestHTTP_DoServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	go func() {
		http.HandleFunc("/api/v1/pods", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`{"foo":"bar"}`))
		})

		if err := http.ListenAndServe("localhost:9876", nil); err != nil {
			log.Fatalf("cannot start http server: %v", err)
		}
	}()

	conFacMock := kube.NewMockRESTConfigFactory(ctrl)
	conFacMock.
		EXPECT().
		New(gomock.Any()).
		Return(&rest.Config{
			Host: "localhost:9876",
		}, nil)

	f := kube.NewDefaultRESTClientFactory(conFacMock, nil, "")

	c, _ := f.Client("apps", "v1")
	rr := rest.NewRequest(c)

	cliFacMock := kube.NewMockRESTClientFactory(ctrl)
	cliFacMock.
		EXPECT().
		Request(gomock.Any()).
		Return(rr, nil)

	hr, err := http.NewRequest("GET", "https://api.kube-apiserver-proxy.dev/api/v1/pods", nil)
	if err != nil {
		t.Fatalf("cannot create http request: %v", err)
	}

	w := httptest.NewRecorder()
	w.Write([]byte(`{"foo":"bar"}`))

	hp := proxy.NewHTTP(
		cliFacMock,
		[]proxy.ResponseBodyTransformer{
			proxy.NewJqResponseBodyTransformer(),
		},
	)

	if err := hp.DoServeHTTP(w, *hr); err != nil {
		t.Errorf("did not expect an error, %v given", err)
	}

	if got, want := w.Body.String(), `{"foo":"bar"}`; got != want {
		t.Errorf("w.Body.String() = %s, want %s", got, want)
	}
}
