package main

import (
	"log"

	"github.com/omissis/kube-apiserver-proxy/internal/app"
	"github.com/omissis/kube-apiserver-proxy/internal/cmd"
)

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

	if err := cmd.NewRootCommand(app.Config{}, versions).Execute(); err != nil {
		log.Fatal(err)
	}
}
