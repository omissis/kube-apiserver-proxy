# Kubernetes Api Server Proxy

This project is a proxy for the Kubernetes API server.

It is designed to be used in a Kubernetes cluster to allow access to parts of the API server from outside the cluster
in a convenient manner for other projects to use.

## Development

In order to start developing, it is recommended to install [asdf][asdf] and [direnv][direnv].
Once those two tools are in place, you should copy the `.envrc.dist` file to `.envrc` and the `tilt_config.json.dist`
file to `tilt_config.json`, and edit the latter's "workdir" value to point to the root of this project.
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

# Step 2: Edit the workdir value in tilt_config.json to point to the root of this project

# Step 3: Install asdf and direnv

# Step 4: Install the dependencies
asdf install

# Step 5: Load the dependencies
direnv allow

# Step 6: start the development environment
make dev-up CLUSTER_VERSION=1.27.1
```

[asdf]: https://asdf-vm.com/
[direnv]: https://direnv.net/
[kind]: https://kind.sigs.k8s.io/
[ingress-nginx]: https://kubernetes.github.io/ingress-nginx/
[metrics-server]: https://github.com/kubernetes-sigs/metrics-server
