variable "access_key" {
	default = "yourID"
}

variable "secret_key" {
	default = "yourSecret"
}

variable "zone" {
	default = "pek3a"
}

provider "qingcloud" {
	access_key = "${var.access_key}"
	secret_key = "${var.secret_key}"
	zone = "${var.zone}"
}
