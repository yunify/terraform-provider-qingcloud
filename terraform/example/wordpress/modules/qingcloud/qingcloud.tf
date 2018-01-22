resource "qingcloud_eip" "foo" {
  bandwidth = 2
}

# qingcloud_keypair upload an SSH public key
# In this example, upload ~/.ssh/id_rsa.pub content.
# You may not have this file in your system, you will need to create your own SSH key.
resource "qingcloud_keypair" "foo" {
  public_key = "${file("~/.ssh/id_rsa.pub")}"
}

resource "qingcloud_instance" "wordpress" {
  image_id         = "centos73x64"
  keypair_ids      = ["${qingcloud_keypair.foo.id}"]
  managed_vxnet_id = "${qingcloud_vxnet.foo.id}"
}

resource "qingcloud_instance" "mysql" {
  image_id         = "centos73x64"
  keypair_ids      = ["${qingcloud_keypair.foo.id}"]
  managed_vxnet_id = "${qingcloud_vxnet.foo.id}"
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
  val2   = "${qingcloud_instance.wordpress.private_ip}"
  val3   = "80"
}

resource "qingcloud_vpc_static" "ssh-wordpress" {
  vpc_id = "${qingcloud_vpc.foo.id}"
  type   = 1
  val1   = "22"
  val2   = "${qingcloud_instance.wordpress.private_ip}"
  val3   = "22"
}

resource "qingcloud_vpc_static" "ssh-mysql" {
  vpc_id = "${qingcloud_vpc.foo.id}"
  type   = 1
  val1   = "2222"
  val2   = "${qingcloud_instance.mysql.private_ip}"
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

resource "qingcloud_security_group_rule" "ssh-wordpress-in" {
  security_group_id = "${qingcloud_security_group.foo.id}"
  protocol          = "tcp"
  priority          = 0
  action            = "accept"
  direction         = 0
  from_port         = 22
  to_port           = 22
}

resource "qingcloud_security_group_rule" "ssh-mysql-in" {
  security_group_id = "${qingcloud_security_group.foo.id}"
  protocol          = "tcp"
  priority          = 0
  action            = "accept"
  direction         = 0
  from_port         = 2222
  to_port           = 2222
}

output "ip" {
  value = "${qingcloud_vpc.foo.public_ip}"
}
