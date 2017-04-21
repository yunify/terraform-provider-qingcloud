
#  ____    ______   ____    
# /\  _`\ /\__  _\ /\  _`\  
# \ \ \L\_\/_/\ \/ \ \ \L\ \
#  \ \  _\L  \ \ \  \ \ ,__/
#   \ \ \L\ \ \_\ \__\ \ \/ 
#    \ \____/ /\_____\\ \_\ 
#     \/___/  \/_____/ \/_/ 
resource "qingcloud_eip" "init"{
	name = "连接第一个主机的地址"
	description = "主机-1"
	billing_mode = "traffic"
	bandwidth = 1
	need_icp = 0
}

# /\  _`\                                 __/\ \__            /\  _`\                                 
# \ \,\L\_\     __    ___   __  __  _ __ /\_\ \ ,_\  __  __   \ \ \L\_\  _ __   ___   __  __  _____   
#  \/_\__ \   /'__`\ /'___\/\ \/\ \/\`'__\/\ \ \ \/ /\ \/\ \   \ \ \L_L /\`'__\/ __`\/\ \/\ \/\ '__`\ 
#    /\ \L\ \/\  __//\ \__/\ \ \_\ \ \ \/ \ \ \ \ \_\ \ \_\ \   \ \ \/, \ \ \//\ \L\ \ \ \_\ \ \ \L\ \
#    \ `\____\ \____\ \____\\ \____/\ \_\  \ \_\ \__\\/`____ \   \ \____/\ \_\\ \____/\ \____/\ \ ,__/
#     \/_____/\/____/\/____/ \/___/  \/_/   \/_/\/__/ `/___/> \   \/___/  \/_/ \/___/  \/___/  \ \ \/ 
#                                                        /\___/                                 \ \_\ 
#                                                        \/__/                                   \/_/ 
resource "qingcloud_securitygroup" "basic"{
	name = "防火墙"
	description = "这是第一个防火墙"
}

# /\  _`\ /\  _`\ /\ \/\ \    
# \ \,\L\_\ \,\L\_\ \ \_\ \   
#  \/_\__ \\/_\__ \\ \  _  \  
#    /\ \L\ \/\ \L\ \ \ \ \ \ 
#    \ `\____\ `\____\ \_\ \_\
#     \/_____/\/_____/\/_/\/_/
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