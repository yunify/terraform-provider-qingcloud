
resource "qingcloud_eip" "init"{
  name = "tf_eip"
  description = ""
  billing_mode = "traffic"
  bandwidth = 20
  need_icp = 0
}

resource "qingcloud_security_group" "basic"{
  name = "防火墙"
  description = "这是第一个防火墙"
}

resource "qingcloud_security_group_rule" "openport" {
  security_group_id = "${qingcloud_security_group.basic.id}"
  protocol = "tcp"
  priority = 0
  action = "accept"
  direction = 0
  from_port = 22
  to_port = 40000
}

# qingcloud_keypair upload an SSH public key
# In this example, upload ~/.ssh/id_rsa.pub content.
# You may not have this file in your system, you will need to create your own SSH key.
resource "qingcloud_keypair" "arthur"{
  name = "arthur"
  public_key = "${file("~/.ssh/id_rsa.pub")}"
}

resource "qingcloud_instance" "init"{
  count = 1
  name = "master-${count.index}"
  image_id = "centos76x64a"
  cpu = "16"
  memory = "32768"
  instance_class = "0"
  managed_vxnet_id="vxnet-0"
  keypair_ids = ["${qingcloud_keypair.arthur.id}"]
  security_group_id ="${qingcloud_security_group.basic.id}"
  eip_id = "${qingcloud_eip.init.id}"
}

resource "null_resource" "install_kubesphere" {
  provisioner "file" {
    destination = "./install.sh"
    source      = "./install.sh"

    connection {
      type        = "ssh"
      user        = "root"
      host        = "${qingcloud_eip.init.addr}"
      private_key = "${file("~/.ssh/id_rsa")}"
      port        = "22"
    }
  }

  provisioner "remote-exec" {
    inline = [
      "sh install.sh"
    ]

    connection {
      type        = "ssh"
      user        = "root"
      host        = "${qingcloud_eip.init.addr}"
      private_key = "${file("~/.ssh/id_rsa")}"
      port        = "22"
    }
  }
}

