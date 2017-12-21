variable "access_key" {
  default = "QINGCLOUD_ACCESS_KEY"
}

variable "secret_key" {
  default = "QINGCLOUD_SECRET_KEY"
}

variable "zone" {
  default = "pek3a"
}

provider "qingcloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  zone       = "${var.zone}"
}
