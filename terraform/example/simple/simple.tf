
resource "qingcloud_eip" "init"{
  name = "连接第一个主机的地址"
  description = "主机-1"
  billing_mode = "traffic"
  bandwidth = 1
  need_icp = 0
}

resource "qingcloud_security_group" "basic"{
  name = "防火墙"
  description = "这是第一个防火墙"
}

resource "qingcloud_security_group_rule" "ssh-wordpress-in" {
  security_group_id = "${qingcloud_security_group.basic.id}"
  protocol = "tcp"
  priority = 0
  action = "accept"
  direction = 0
  from_port = 22
  to_port = 22
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
  image_id = "centos7x64d"
  instance_class = "0"
  managed_vxnet_id="vxnet-0"
  keypair_ids = ["${qingcloud_keypair.arthur.id}"]
  security_group_id ="${qingcloud_security_group.basic.id}"
  eip_id = "${qingcloud_eip.init.id}"
}

