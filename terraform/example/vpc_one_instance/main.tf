resource "qingcloud_eip" "foo" {
  bandwidth = 2
}

resource "qingcloud_keypair" "foo" {
  public_key = "${file("~/.ssh/id_rsa.pub")}"
}

resource "qingcloud_instance" "foo" {
  image_id         = "centos73x64"
  keypair_ids      = ["${qingcloud_keypair.foo.id}"]
  managed_vxnet_id = "${qingcloud_vxnet.foo.id}"
}

resource "null_resource" "run_docker_nginx" {
  depends_on = ["qingcloud_vpc_static.ssh-portforward",
    "qingcloud_vpc_static.http-portforward",
    "qingcloud_security_group_rule.ssh-in",
    "qingcloud_security_group_rule.http-in",
    "qingcloud_instance.foo",
  ]

  provisioner "remote-exec" {
    inline = [
      "curl -fsSL get.docker.com -o get-docker.sh",
      "sh get-docker.sh",
      "systemctl start docker",
      "docker run --name docker-nginx -d -p 80:80 nginx",
    ]

    connection {
      type        = "ssh"
      user        = "root"
      host        = "${qingcloud_vpc.foo.public_ip}"
      private_key = "${file("~/.ssh/id_rsa")}"
    }
  }
}

resource "qingcloud_security_group" "foo" {
  name = "first_sg"
}

resource "qingcloud_vpc" "foo" {
  security_group_id = "${qingcloud_security_group.foo.id}"
  vpc_network       = "192.168.0.0/16"
  eip_id            = "${qingcloud_eip.foo.id}"
}

resource "qingcloud_vxnet" "foo" {
  type       = 1
  vpc_id     = "${qingcloud_vpc.foo.id}"
  ip_network = "192.168.0.0/24"
}

resource "qingcloud_vpc_static" "http-portforward" {
  vpc_id = "${qingcloud_vpc.foo.id}"
  type   = 1
  val1   = "80"
  val2   = "${qingcloud_instance.foo.private_ip}"
  val3   = "80"
}

resource "qingcloud_vpc_static" "ssh-portforward" {
  vpc_id = "${qingcloud_vpc.foo.id}"
  type   = 1
  val1   = "22"
  val2   = "${qingcloud_instance.foo.private_ip}"
  val3   = "22"
}

resource "qingcloud_security_group_rule" "http-in" {
  security_group_id = "${qingcloud_security_group.foo.id}"
  protocol          = "tcp"
  priority          = 0
  action            = "accept"
  direction         = 0
  from_port         = 80
  to_port           = 80
}

resource "qingcloud_security_group_rule" "ssh-in" {
  security_group_id = "${qingcloud_security_group.foo.id}"
  protocol          = "tcp"
  priority          = 0
  action            = "accept"
  direction         = 0
  from_port         = 22
  to_port           = 22
}

output "ip" {
  value = "${qingcloud_vpc.foo.public_ip}"
}
