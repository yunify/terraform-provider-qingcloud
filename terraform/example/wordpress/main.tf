resource "qingcloud_eip" "foo" {
  bandwidth = 2
}

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

resource "null_resource" "run_docker_wordpress" {
  depends_on = ["qingcloud_vpc_static.ssh-wordpress",
    "qingcloud_security_group_rule.ssh-wordpress-in",
    "qingcloud_instance.wordpress",
    "null_resource.run_docker_mysql"
  ]

  provisioner "remote-exec" {
    inline = [
      "curl -fsSL get.docker.com -o get-docker.sh",
      "sh get-docker.sh",
      "systemctl start docker",
      "docker run --name wordpress -d -p 80:80 -e WORDPRESS_DB_HOST=${qingcloud_instance.mysql.private_ip} -e WORDPRESS_DB_PASSWORD=wordpress wordpress",
    ]

    connection {
      type        = "ssh"
      user        = "root"
      host        = "${qingcloud_vpc.foo.public_ip}"
      private_key = "${file("~/.ssh/id_rsa")}"
    }
  }
}

resource "null_resource" "run_docker_mysql" {
  depends_on = ["qingcloud_vpc_static.ssh-mysql",
    "qingcloud_security_group_rule.ssh-mysql-in",
    "qingcloud_instance.wordpress",
  ]

  provisioner "remote-exec" {
    inline = [
      "curl -fsSL get.docker.com -o get-docker.sh",
      "sh get-docker.sh",
      "systemctl start docker",
      "docker run --name wordpress-mysql -v /datadir:/var/lib/mysql  -p 3306:3306 -e MYSQL_ROOT_PASSWORD=wordpress -d  mysql:5.7",
    ]

    connection {
      type        = "ssh"
      user        = "root"
      host        = "${qingcloud_vpc.foo.public_ip}"
      private_key = "${file("~/.ssh/id_rsa")}"
      port        = 2222
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
