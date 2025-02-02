data "yandex_compute_image" "coi" {
  family = "container-optimized-image"
}

resource "yandex_compute_instance_group" "coi-traefik-ig" {
  name = "coi-traefik"

  folder_id          = var.folder_id
  service_account_id = yandex_iam_service_account.ig-sa.id
  instance_template {
    platform_id        = "standard-v3"
    service_account_id = yandex_iam_service_account.coi.id
    resources {
      memory = 2
      cores  = 2
    }
    boot_disk {
      mode = "READ_WRITE"
      initialize_params {
        image_id = data.yandex_compute_image.coi.id
        size     = 33
        type     = "network-ssd"
      }
    }
    network_interface {
      network_id = yandex_vpc_subnet.coi-a.network_id
      subnet_ids = [yandex_vpc_subnet.coi-a.id]
    }

    metadata = {
      ssh-keys = "ubuntu:${file("~/.ssh/id_rsa.pub")}"
      docker-compose = templatefile("./files/spec.yaml", {
        YC_GROUP_ID    = yandex_logging_group.coi.id,
        TRAEFIK_IMAGE  = local.traefik_image,
        SIDECAR_IMAGE  = local.sidecar_image,
        HTTPBIN_IMAGE  = local.httpbin_image,
        CERTIFICATE_ID = var.certificate_id,
      })
      user-data = data.cloudinit_config.user-data.rendered
    }
  }

  scale_policy {
    fixed_scale {
      size = 2
    }
  }

  allocation_policy {
    zones = ["ru-central1-a"]
  }

  deploy_policy {
    max_unavailable = 1
    max_creating    = 1
    max_expansion   = 1
    max_deleting    = 1
  }

  load_balancer {
    target_group_name        = "traefik-target-group"
    target_group_description = "Целевая группа Network Load Balancer для Traefik"
  }

  depends_on = [
    yandex_iam_service_account.coi,
    yandex_logging_group.coi,
    yandex_vpc_subnet.coi-a,
    null_resource.build,
    null_resource.traefik,
    null_resource.httpbin,
  ]
}


