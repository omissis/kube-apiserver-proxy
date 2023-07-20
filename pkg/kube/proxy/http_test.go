//go:build unit

package proxy_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	rest "k8s.io/client-go/rest"
	utiltesting "k8s.io/client-go/util/testing"
	"k8s.io/kubectl/pkg/scheme"

	"github.com/omissis/kube-apiserver-proxy/pkg/kube"
	"github.com/omissis/kube-apiserver-proxy/pkg/kube/proxy"
)

func TestHTTP_DoServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testServer, _, obj := testServerEnv(t, 200)
	defer testServer.Close()

	c, err := restClient(testServer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rr := rest.NewRequest(c)

	cliFacMock := kube.NewMockRESTClientFactory(ctrl)
	cliFacMock.
		EXPECT().
		Request(gomock.Any()).
		Return(rr, nil)

	hp := proxy.NewHTTP(
		cliFacMock,
		[]proxy.ResponseBodyTransformer{
			proxy.NewJqResponseBodyTransformer(),
		},
	)

	r, err := http.NewRequest("GET", "https://api.kube-apiserver-proxy.test/api/v1/pods?jq=.kind", nil)
	if err != nil {
		t.Fatalf("cannot create http request: %v", err)
	}

	w := httptest.NewRecorder()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	if err := hp.DoServeHTTP(ctx, w, *r); err != nil {
		t.Errorf("did not expect an error, %v given", err)
	}

	if got, want := strings.TrimSpace(w.Body.String()), `"`+obj.Kind+`"`; got != want {
		t.Errorf("got = %s, want %s", got, want)
	}
}

func testServerEnv(t *testing.T, statusCode int) (*httptest.Server, *utiltesting.FakeHandler, *corev1.PodList) {
	podList := &corev1.PodList{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "PodList"},
		Items:    []corev1.Pod{},
	}

	body, _ := runtime.Encode(scheme.Codecs.LegacyCodec(corev1.SchemeGroupVersion), podList)

	fakeHandler := utiltesting.FakeHandler{
		StatusCode:   statusCode,
		ResponseBody: string(body),
		T:            t,
	}

	testServer := httptest.NewServer(&fakeHandler)

	return testServer, &fakeHandler, podList
}

func restClient(testServer *httptest.Server) (*rest.RESTClient, error) {
	return rest.RESTClientFor(&rest.Config{
		Host: testServer.URL,
		ContentConfig: rest.ContentConfig{
			GroupVersion:         &v1.SchemeGroupVersion,
			NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
		},
		Username: "user",
		Password: "pass",
	})
}
