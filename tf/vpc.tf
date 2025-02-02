resource "yandex_vpc_network" "coi" {}

resource "yandex_vpc_subnet" "coi-a" {
  zone       = var.zone
  network_id = "${yandex_vpc_network.coi.id}"
  v4_cidr_blocks = ["10.5.0.0/24"]

  route_table_id = yandex_vpc_route_table.coi-egress.id
}

resource "yandex_vpc_gateway" "egress" {
  name = "traefik-egress"
  shared_egress_gateway {}
}

resource "yandex_vpc_route_table" "coi-egress" {
  network_id = yandex_vpc_network.coi.id
  static_route {
    destination_prefix = "0.0.0.0/0"
    gateway_id         = yandex_vpc_gateway.egress.id
  }
}