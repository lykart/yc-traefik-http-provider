resource "yandex_iam_service_account" "coi" {
  name      = "coi-logger"
  folder_id = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_binding" "logger" {
  for_each = toset([
    "logging.writer",
    "container-registry.images.puller"
  ])
  role      = each.value
  folder_id = var.folder_id
  members = [
    "serviceAccount:${yandex_iam_service_account.coi.id}",
  ]
  sleep_after = 5
}


resource "yandex_resourcemanager_folder_iam_binding" "certs" {
  folder_id = data.yandex_cm_certificate.cert.folder_id
  members = [
    "serviceAccount:${yandex_iam_service_account.coi.id}",
  ]
  role        = "certificate-manager.certificates.downloader"
  sleep_after = 5
}

resource "yandex_iam_service_account" "ig-sa" {
  name        = "ig-sa"
  description = "Сервисный аккаунт для управления группой ВМ."
}

resource "yandex_resourcemanager_folder_iam_member" "compute_editor" {
  folder_id  = var.folder_id
  for_each = toset([
    "compute.editor",
    "load-balancer.editor",
    "iam.serviceAccounts.user",
  ])
  role = each.value
  member     = "serviceAccount:${yandex_iam_service_account.ig-sa.id}"
  depends_on = [
    yandex_iam_service_account.ig-sa,
  ]
}