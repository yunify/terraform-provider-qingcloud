# ------------------------------------------------------------------
# EIP
resource "qingcloud_eip" "init"{
	name = "连接第一个主机的地址"
	description = "主机-1"
	billing_mode = "traffic"
	bandwidth = 1
	need_icp = 0
}

# ------------------------------------------------------------------
# 						Security Group
resource "qingcloud_securitygroup" "basic"{
	name = "防火墙"
	description = "这是第一个防火墙"
}

# ------------------------------------------------------------------
# 						SSH
resource "qingcloud_keypair" "arthur"{
	keypair_name = "arthur"
	description = "sdfafd"
	public_key = "${file("~/.ssh/id_rsa.pub")}"
}

resource "qingcloud_instance" "init"{
	count = 1
	name = "master-${count.index}"
	image_id = "centos7x64d"
	instance_type = "c1m1"
	instance_class = "0"
	vxnet_id="vxnet-0"
	keypair_ids = ["${qingcloud_keypair.arthur.id}"]
	security_group_id ="${qingcloud_securitygroup.basic.id}"
}


resource "qingcloud_eip_associate" "init"{
	resource_type = "instance"
	resource_id = "${qingcloud_instance.init.id}"
	eip = "${qingcloud_eip.init.id}"
}