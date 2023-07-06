package main

import (
	"log"

	"github.com/omissis/kube-apiserver-proxy/internal/cmd"
)

//go:generate mockgen -source pkg/kube/config.go -destination pkg/kube/config_mock.gen.go -package kube
//go:generate mockgen -source pkg/kube/client.go -destination pkg/kube/client_mock.gen.go -package kube

var (
	version   = "unknown"
	gitCommit = "unknown"
	buildTime = "unknown"
	goVersion = "unknown"
	osArch    = "unknown"
)

func main() {
	versions := map[string]string{
		"version":   version,
		"gitCommit": gitCommit,
		"buildTime": buildTime,
		"goVersion": goVersion,
		"osArch":    osArch,
	}

	if err := cmd.NewRootCommand(versions).Execute(); err != nil {
		log.Fatal(err)
	}
}
