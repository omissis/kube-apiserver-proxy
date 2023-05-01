# Kubernetes Api Server Proxy

> THIS PROJECT IS IN DEVELOPMENT AND IT IS NOT READY FOR PRODUCTION USE
>
> DO NOT CONSIDER IT UNLESS YOU ARE WILLING TO CONTRIBUTE TO IT

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/omissis/kube-apiserver-proxy?style=flat)](https://github.com/omissis/kube-apiserver-proxy/releases/latest)
[![GitHub Workflow Status (event)](https://img.shields.io/github/actions/workflow/status/omissis/kube-apiserver-proxy/development.yaml?style=flat)](https://github.com/omissis/kube-apiserver-proxy/actions?workflow=development)
[![License](https://img.shields.io/github/license/omissis/kube-apiserver-proxy?style=flat)](/LICENSE)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/omissis/kube-apiserver-proxy?style=flat)](https://tip.golang.org/doc/go1.20)
[![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/omissis/kube-apiserver-proxy?style=flat)](https://github.com/omissis/kube-apiserver-proxy)
[![GitHub repo file count (file type)](https://img.shields.io/github/directory-file-count/omissis/kube-apiserver-proxy?style=flat)](https://github.com/omissis/kube-apiserver-proxy)
[![GitHub all releases](https://img.shields.io/github/downloads/omissis/kube-apiserver-proxy/total?style=flat)](https://github.com/omissis/kube-apiserver-proxy)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/omissis/kube-apiserver-proxy?style=flat)](https://github.com/omissis/kube-apiserver-proxy/commits)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=flat)](https://conventionalcommits.org)
[![Codecov](https://img.shields.io/codecov/c/gh/omissis/kube-apiserver-proxy?style=flat&token=lPWlXd3MVK)](https://codecov.io/gh/omissis/kube-apiserver-proxy)
[![Code Climate maintainability](https://img.shields.io/codeclimate/maintainability/omissis/kube-apiserver-proxy?style=flat)](https://codeclimate.com/github/omissis/kube-apiserver-proxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/omissis/kube-apiserver-proxy)](https://goreportcard.com/report/github.com/omissis/kube-apiserver-proxy)

This project is a proxy for the Kubernetes API server.

It is designed to be used in a Kubernetes cluster to allow access to parts of the API server from outside the cluster
in a convenient manner for other projects to use.

## Contributing

### Setting up the environment

In order to start developing, it is recommended to install [asdf][asdf] and [direnv][direnv].
Once those two tools are in place, you should copy the `.envrc.dist` file to `.envrc` and the `tilt_config.json.dist`
file to `tilt_config.json`, they already contain sensible defaults for the development environment.
Then, you can run `asdf install` to install the correct version of the required tools and
`direnv allow` to load them within the context of this folder.
Once that is done, you can run `make dev-up CLUSTER_VERSION=1.27.1` to startup a Kubernetes cluster for the project:
this will spawn a Tilt.dev process, which in turn will start a Kubernetes cluster using [kind][kind] and
deploy [ingress-nginx][ingress-nginx] and [metrics-server][metrics-server].

In a nutshell:

```shell
# Step 1: copy the .envrc.dist and tilt_config.json.dist files
cp .envrc.dist .envrc
cp tilt_config.json.dist tilt_config.json

# Step 2: install asdf and direnv using your favorite package manager

# Step 3: install the asdf dependencies
asdf install

# Step 4: load the asdf dependencies
direnv allow

# Step 5: install the brew dependencies (MacOS only, for other OSes, please install the dependencies manually)
make tools-brew

# Step 6: install the golang dependencies
make tools-go

# Step 7: install npm dependencies
npm install

# Step 8: start the development environment
make dev-up CLUSTER_VERSION=1.27.1

# Step 9: stop the development environment
make dev-down
```

### Development

The project offers a `Makefile` containing most of the commands you'll need for development.
In there, you'll find targets for running several linters and formatters, for building and releasing the project,
for starting and stopping the dev environment, for running tests and for generating the code and the graphql schemas.
Feel free to explore it to find out more, and don't forget to have a look at the `scripts/` folder
for more details on the implementation.

### Architecture

For more information on the architectural decisions that have been made, refer to the [docs/arch/](./docs/arch/) folder.

[asdf]: https://asdf-vm.com/
[direnv]: https://direnv.net/
[kind]: https://kind.sigs.k8s.io/
[ingress-nginx]: https://kubernetes.github.io/ingress-nginx/
[metrics-server]: https://github.com/kubernetes-sigs/metrics-server
