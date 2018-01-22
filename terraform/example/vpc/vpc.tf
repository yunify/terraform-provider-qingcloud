
#  ____    ______   ____    
# /\  _`\ /\__  _\ /\  _`\  
# \ \ \L\_\/_/\ \/ \ \ \L\ \
#  \ \  _\L  \ \ \  \ \ ,__/
#   \ \ \L\ \ \_\ \__\ \ \/ 
#    \ \____/ /\_____\\ \_\ 
#     \/___/  \/_____/ \/_/
                          
resource "qingcloud_eip" "vpc"{
  name = "first"
  description = "first one"
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
resource "qingcloud_security_group" "basic"{
  name = "防火墙"
  description = "这是第一个防火墙"
}

resource "qingcloud_security_group_rule" "allow-in-80"{
  name = "允许使用80"
  security_group_id  = "${qingcloud_security_group.basic.id}"
  protocol = "tcp"
  priority = 1
  action = "accept"
  direction = 0
  from_port = "80"
  to_port = "80"
}
resource "qingcloud_security_group_rule" "allow-in-81"{
  name = "允许使用81"
  security_group_id = "${qingcloud_security_group.basic.id}"
  protocol = "tcp"
  priority = 1
  action = "accept"
  direction = 0
  from_port = "81"
  to_port = "81"
}

#  ____    ____    __  __     
# /\  _`\ /\  _`\ /\ \/\ \    
# \ \,\L\_\ \,\L\_\ \ \_\ \   
#  \/_\__ \\/_\__ \\ \  _  \  
#    /\ \L\ \/\ \L\ \ \ \ \ \ 
#    \ `\____\ `\____\ \_\ \_\
#     \/_____/\/_____/\/_/\/_/

# qingcloud_keypair upload an SSH public key
# In this example, upload ~/.ssh/id_rsa.pub content.
# You may not have this file in your system, you will need to create your own SSH key.
resource "qingcloud_keypair" "arthur"{
  name = "arthur"
  public_key = "${file("~/.ssh/id_rsa.pub")}"
}

#  ____                     __
# /\  _`\                  /\ \__
# \ \ \L\ \    ___   __  __\ \ ,_\    __   _ __
#  \ \ ,  /   / __`\/\ \/\ \\ \ \/  /'__`\/\`'__\
#   \ \ \\ \ /\ \L\ \ \ \_\ \\ \ \_/\  __/\ \ \/
#    \ \_\ \_\ \____/\ \____/ \ \__\ \____\\ \_\
#     \/_/\/ /\/___/  \/___/   \/__/\/____/ \/_/
                                               
resource "qingcloud_vpc" "vpc"{
  name = "vpc-network"
  type = 1
  vpc_network = "172.16.0.0/16"
  security_group_id = "${qingcloud_security_group.basic.id}"
  description = "测试的网络"
  eip_id = "${qingcloud_eip.vpc.id}"
}

resource "qingcloud_vpc_static" "foo1"{
  vpc_id = "${qingcloud_vpc.vpc.id}"
  type = 1
  val1 = "80"
  val2 = "${qingcloud_instance.master.private_ip}"
  val3 = "80"
}

resource "qingcloud_vpc_static" "foo2"{
  vpc_id = "${qingcloud_vpc.vpc.id}"
  type = 1
  val1 = "81"
  val2 = "${qingcloud_instance.master.private_ip}"
  val3 = "81"
}

#  __  __                        __      
# /\ \/\ \                      /\ \__   
# \ \ \ \ \  __  _   ___      __\ \ ,_\  
#  \ \ \ \ \/\ \/'\/' _ `\  /'__`\ \ \/  
#   \ \ \_/ \/>  <//\ \/\ \/\  __/\ \ \_ 
#    \ `\___//\_/\_\ \_\ \_\ \____\\ \__\
#     `\/__/ \//\/_/\/_/\/_/\/____/ \/__/
                                       
resource "qingcloud_vxnet" "vx"{
  name = "app vxnet"
  type = 1
  description = "应用的网络"
  vpc_id = "${qingcloud_vpc.vpc.id}"
  ip_network = "172.16.1.0/24"
}


#  ______                   __                                     
# /\__  _\                 /\ \__                                  
# \/_/\ \/     ___     ____\ \ ,_\    __      ___     ___     __   
#    \ \ \   /' _ `\  /',__\\ \ \/  /'__`\  /' _ `\  /'___\ /'__`\ 
#     \_\ \__/\ \/\ \/\__, `\\ \ \_/\ \L\.\_/\ \/\ \/\ \__//\  __/ 
#     /\_____\ \_\ \_\/\____/ \ \__\ \__/.\_\ \_\ \_\ \____\ \____\
#     \/_____/\/_/\/_/\/___/   \/__/\/__/\/_/\/_/\/_/\/____/\/____/
resource "qingcloud_instance" "master"{
  image_id = "trustysrvx64f"
  instance_class = "0"
  managed_vxnet_id = "${qingcloud_vxnet.vx.id}"
  keypair_ids = ["${qingcloud_keypair.arthur.id}"]
  security_group_id ="${qingcloud_security_group.basic.id}"
}

resource "qingcloud_instance" "slave"{
  count = 3

  name = "slave-${count.index}"
  image_id = "trustysrvx64f"
  instance_class = "0"
  managed_vxnet_id = "${qingcloud_vxnet.vx.id}"
  keypair_ids = ["${qingcloud_keypair.arthur.id}"]
  security_group_id ="${qingcloud_security_group.basic.id}"
}




