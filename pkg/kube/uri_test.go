//go:build unit

package kube

import "testing"

func TestGetGroupVersionFromURI(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc        string
		uri         string
		wantGroup   string
		wantVersion string
		wantErr     bool
		wantErrStr  string
	}{
		{
			desc:        "empty url",
			uri:         "",
			wantGroup:   "",
			wantVersion: "",
			wantErr:     true,
			wantErrStr:  "uri '' is not supported",
		},
		{
			desc:        "Core API: one part",
			uri:         "/api",
			wantGroup:   "core",
			wantVersion: "",
			wantErr:     false,
			wantErrStr:  "",
		},
		{
			desc:        "Core v1 API: two parts",
			uri:         "/api/v1",
			wantGroup:   "core",
			wantVersion: "v1",
			wantErr:     false,
			wantErrStr:  "",
		},
		{
			desc:        "Core v1 API: three parts",
			uri:         "/api/v1/pods",
			wantGroup:   "core",
			wantVersion: "v1",
			wantErr:     false,
			wantErrStr:  "",
		},
		{
			desc:        "Core v1 API: five parts",
			uri:         "/api/v1/namespaces/default/pods",
			wantGroup:   "core",
			wantVersion: "v1",
			wantErr:     false,
			wantErrStr:  "",
		},
		{
			desc:        "Version API: one part",
			uri:         "/apis",
			wantGroup:   "apis",
			wantVersion: "",
			wantErr:     false,
			wantErrStr:  "",
		},
		{
			desc:        "Version API: two parts",
			uri:         "/apis/apps",
			wantGroup:   "",
			wantVersion: "",
			wantErr:     true,
			wantErrStr:  "uri '/apis/apps' has less than 3 parts in it",
		},
		{
			desc:        "Apps v1 API: three parts",
			uri:         "/apis/apps/v1",
			wantGroup:   "apps",
			wantVersion: "v1",
			wantErr:     false,
			wantErrStr:  "",
		},
		{
			desc:        "Apps v1 API: four parts",
			uri:         "/apis/apps/v1/deployments",
			wantGroup:   "apps",
			wantVersion: "v1",
			wantErr:     false,
			wantErrStr:  "",
		},
	}
	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			gotFirst, gotSecond, err := GetGroupVersionFromURI(tC.uri)

			if (err != nil) != tC.wantErr {
				t.Errorf("GetGroupVersionFromURI() error = %v, wantErr %v", err, tC.wantErr)
			}

			if gotFirst != tC.wantGroup {
				t.Errorf("GetGroupVersionFromURI() gotFirst = %v, want %v", gotFirst, tC.wantGroup)
			}

			if gotSecond != tC.wantVersion {
				t.Errorf("GetGroupVersionFromURI() gotSecond = %v, want %v", gotSecond, tC.wantVersion)
			}
		})
	}
}
