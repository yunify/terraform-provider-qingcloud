---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_loadbalancer"
sidebar_current: "docs-qingcloud-resource-loadbalancer"
description: |-
  Provides a  LoadBalancer resource.
---

# Qingcloud\_loadbalancer

Provides a  LoadBalancer resource.  
External loadbalancer use vxnet who's id is vxnet-0.  
Internal loadbalancer use vxnet who under vpc's control.  

## Example Usage

```
# Create a new external loadbalancer.
resource "qingcloud_eip" "foo" {
    bandwidth = 2
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo.id}"]
}
```
```
# Create a new loadbalancer with tags.
resource "qingcloud_eip" "foo" {
    bandwidth = 2
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo.id}"]
	tag_ids = ["${qingcloud_tag.test.id}",
				"${qingcloud_tag.test2.id}"]
}
resource "qingcloud_tag" "test"{
	name="11"
}
resource "qingcloud_tag" "test2"{
	name="22"
}
```

```
# Create a new loadbalancer with muti eips.
resource "qingcloud_eip" "foo1" {
    bandwidth = 2
}
resource "qingcloud_eip" "foo2" {
    bandwidth = 2
}
resource "qingcloud_eip" "foo3" {
    bandwidth = 2
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo1.id}","${qingcloud_eip.foo2.id}","${qingcloud_eip.foo3.id}"]
}
```

```
# Create a new internal loadbalancer.
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
}

resource "qingcloud_vxnet" "foo" {
    type = 1
	vpc_id = "${qingcloud_vpc.foo.id}"
	ip_network = "192.168.0.0/24"
}

resource "qingcloud_loadbalancer" "foo" {
	vxnet_id = "${qingcloud_vxnet.foo.id}"
	private_ips = ["192.168.0.3"]
}
```
## Argument Reference

The following arguments are supported:

* `type` - (Optional, Default 0)Max connections of loadbalancer, 0 - 5000,1 - 20000, 2 - 40000, 3 - 100000, 4 - 200000, 5-500000. 
* `private_ips` - (Optional) Internal loadbalancer can specify private_ip by this field . 
* `node_count` - (Optional) Only external loadbalancer can this field. Number of nodes per public network IP. When the number of nodes is greater than 1, the load balancer is deployed in a multi-node hot standby mode . 
* `eips` - (Optional) Eips loadbalancer want to use .If create external loadbalancer ,at least one eip ,at most four . If create internal loadbalancer ,at most one. 
* `name` - (Optional) The name of loadbalancer . 
* `description`- (Optional) The description of loadbalancer .
* `vxnet` - (Optional, Default vxnet-0 , ForceNew)If the value is vxnet-0 then this loadbalancer is external, else is internal. 
* `tag_ids` - (Optional) tag ids , loadbalancer wants to use.
* `security_group_id` - (Optional) security_group_id , loadbalancer wants to use.
* `http_header_size` - (Optional, Default 15kbytes) max length of http_header_size, 1-127kbytes, This parameter affects the maximum number of connections.

## Attributes Reference

The following attributes are exported:

* `type` - Max connections of loadbalancer, 0 - 5000,1 - 20000, 2 - 40000, 3 - 100000, 4 - 200000, 5-500000. 
* `private_ips` - Internal loadbalancer can specify private_ip by this field . 
* `eips` - Eips loadbalancer want to use .If create external loadbalancer ,at least one eip ,at most four . If create internal loadbalancer ,at most one. 
* `name` - The name of loadbalancer . 
* `description`- The description of loadbalancer .
* `vxnet` - If the value is vxnet-0 then this loadbalancer is external, else is internal. 
* `tag_ids` - tag ids , loadbalancer wants to use.
* `tag_names` - compute by tag ids.
* `security_group_id` - security_group_id , loadbalancer wants to use.
* `http_header_size` - max length of http_header_size, 1-127kbytes, This parameter affects the maximum number of connections.
* `node_count` -  Only external loadbalancer use this field. Number of nodes per public network IP. When the number of nodes is greater than 1, the load balancer is deployed in a multi-node hot standby mode . 


