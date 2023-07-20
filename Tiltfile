config.define_string_list("services")

parsed_config = config.parse()

load_dynamic("./configs/tiltfiles/setup.tiltfile")

for service in parsed_config.get("services", []):
  load_dynamic("./configs/tiltfiles/%s.tiltfile" % (service))

load('ext://restart_process', 'docker_build_with_restart')
load('ext://helm_resource', 'helm_resource', 'helm_repo')

docker_build_with_restart(
  "kube-apiserver-proxy",
  ".",
  target="dev",
  live_update=[
    sync("./internal", "/app/internal"),
    sync("./pkg", "/app/pkg"),
    sync("./go.mod", "/app/go.mod"),
    sync("./go.sum", "/app/go.sum"),
    sync("./main.go", "/app/main.go"),
  ],
  build_args={},
  entrypoint=["go", "run", "main.go", "serve"],
)

helm_resource(
  'kube-apiserver-proxy',
  './deployments/helm/kube-apiserver-proxy',
  release_name='kube-apiserver-proxy',
  namespace='default',
  flags=['--values', './configs/helm-values/kube-apiserver-proxy.yaml'],
  image_deps=['kube-apiserver-proxy'],
  image_keys=[('image.repository', 'image.tag')],
)

k8s_resource(
  workload="kube-apiserver-proxy",
  links=[
    link("https://api.kube-apiserver-proxy.dev/api/v1/pods", "kube-apiserver-proxy"),
  ],
  labels="projects",
  trigger_mode=TRIGGER_MODE_AUTO
)
