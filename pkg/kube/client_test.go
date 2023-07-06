package kube_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/omissis/kube-apiserver-proxy/pkg/kube"
	gomock "go.uber.org/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	utiltesting "k8s.io/client-go/util/testing"
	"k8s.io/kubectl/pkg/scheme"
)

func testServerEnv(t *testing.T, groupVersion schema.GroupVersion) (*httptest.Server, *utiltesting.FakeHandler, *metav1.Status) {
	status := &metav1.Status{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Status"},
		Status:   "Success",
	}

	expectedBody, _ := runtime.Encode(scheme.Codecs.LegacyCodec(groupVersion), status)

	fakeHandler := utiltesting.FakeHandler{
		StatusCode:   200,
		ResponseBody: string(expectedBody),
		T:            t,
	}

	return httptest.NewServer(&fakeHandler), &fakeHandler, status
}

func TestNewK8sRESTClientFactory(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testServer, _, _ := testServerEnv(t, schema.GroupVersion{})
	defer testServer.Close()

	cfMock := kube.NewMockRESTConfigFactory(ctrl)

	f := kube.NewDefaultK8sRESTClientFactory(cfMock, nil, "")

	if f == nil {
		t.Error("expected non-nil K8sRESTClientFactory")
	}
}

func TestNewK8sRESTClientFactory_Client(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		desc       string
		group      string
		version    string
		wantErr    bool
		wantErrMsg string
	}{
		{
			desc:       "apps group",
			group:      "apps",
			version:    "v1",
			wantErr:    false,
			wantErrMsg: "",
		},
	}

	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			gv := schema.GroupVersion{
				Group:   tC.group,
				Version: tC.version,
			}

			testServer, _, _ := testServerEnv(t, gv)
			defer testServer.Close()

			cfMock := kube.NewMockRESTConfigFactory(ctrl)
			cfMock.
				EXPECT().
				New(gomock.Any()).
				Return(&rest.Config{
					Host: testServer.URL,
				}, nil)

			f := kube.NewDefaultK8sRESTClientFactory(cfMock, nil, "")

			got, err := f.Client(tC.group, tC.version)
			if (err != nil) != tC.wantErr {
				t.Errorf("expected error: %v, got: %v", tC.wantErr, err)
			}

			if err != nil && err.Error() != tC.wantErrMsg {
				t.Errorf("expected error message: %v, got: %v", tC.wantErrMsg, err.Error())
			}

			if err == nil && got == nil {
				t.Errorf("expected a pointer to a client, got nil")
			}

			if cmp.Equal(got.APIVersion(), gv) == false {
				t.Errorf("expected api version: %v, got: %v", gv, got.APIVersion())
			}
		})
	}
}

func TestNewK8sRESTClientFactory_Request(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		desc       string
		group      string
		version    string
		want       string
		wantErr    bool
		wantErrMsg string
	}{
		{
			desc:       "apps group",
			group:      "apps",
			version:    "v1",
			want:       "/api/v1/pods",
			wantErr:    false,
			wantErrMsg: "",
		},
	}
	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			gv := schema.GroupVersion{
				Group:   tC.group,
				Version: tC.version,
			}

			testServer, _, _ := testServerEnv(t, gv)
			defer testServer.Close()

			cfMock := kube.NewMockRESTConfigFactory(ctrl)
			cfMock.
				EXPECT().
				New(gomock.Any()).
				Return(&rest.Config{
					Host: testServer.URL,
				}, nil)

			f := kube.NewDefaultK8sRESTClientFactory(cfMock, nil, "")

			req, err := http.NewRequest("GET", "https://api.kube-apiserver-proxy.dev/api/v1/pods", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			got, err := f.Request(*req)
			if (err != nil) != tC.wantErr {
				t.Errorf("expected error: %v, got: %v", tC.wantErr, err)
			}

			if err != nil && err.Error() != tC.wantErrMsg {
				t.Errorf("expected error message: %v, got: %v", tC.wantErrMsg, err.Error())
			}

			if err == nil && got == nil {
				t.Errorf("expected a pointer to a request, got nil")
			}

			want := testServer.URL + tC.want

			if got.URL().String() != want {
				t.Errorf("expected request url: %v, got: %v", want, got.URL())
			}
		})
	}
}
