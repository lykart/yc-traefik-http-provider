resource "yandex_container_registry" "coi-traefik" {
  name = "coi-traefik"
}

data "external" "git_hash" {
  program = [
    "git",
    "log",
    "--pretty={\"sha\":\"%H\"}",
    "-1",
    "HEAD"
  ]
}

locals {
  sidecar_image = "cr.yandex/${yandex_container_registry.coi-traefik.id}/yc-traefik-http-provider:${data.external.git_hash.result.sha}"
}

resource "null_resource" "build" {
  triggers = {
    git_sha = data.external.git_hash.result.sha
  }

  provisioner "local-exec" {
    command = "cd .. && docker build --platform linux/amd64 -t ${local.sidecar_image} . && docker push ${local.sidecar_image}"
  }
}
