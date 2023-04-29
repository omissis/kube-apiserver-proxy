config.define_string("workdir")
config.define_string_list("projects")
config.define_string_list("services")

parsed_config = config.parse()

projects = parsed_config.get("projects", [])
services = parsed_config.get("services", [])

load_dynamic("./configs/tiltfiles/setup.tiltfile")

path = parsed_config.get("workdir")

for project in projects:
  load_dynamic("%s/%s/Tiltfile" % (path, project))

for service in services:
  load_dynamic("./configs/tiltfiles/%s.tiltfile" % (service))
