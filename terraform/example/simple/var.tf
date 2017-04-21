variable "id" {
	default = "yourID"
}

variable "secret" {
	default = "yourSecret"
}

variable "zone" {
	default = "pek3a"
}

provider "qingcloud" {
	id = "${var.id}"
	secret = "${var.secret}"
	zone = "${var.zone}"
}