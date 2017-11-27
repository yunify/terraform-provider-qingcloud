---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_vpc"
sidebar_current: "docs-qingcloud-resource-vpc"
description: |-
  Provides a  Vpc resource.
---

# Qingcloud\_vpc

Provides a  Vpc resource.  
Use this resource in SDN 2.0 environment.

## Example Usage

```
# Create a new VPC.
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
}
```
```
# Create a new VPC with tags.
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
	tag_ids = ["${qingcloud_tag.test.id}",
				"${qingcloud_tag.test2.id}"]
}
resource "qingcloud_tag" "test"{
	name="tag1"
}
resource "qingcloud_tag" "test2"{
	name="tag2"
}
```
```
# Create a new VPC with eip.
resource "qingcloud_eip" "foo" {
    bandwidth = 2
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	eip_id = "${qingcloud_eip.foo.id}"
	vpc_network = "192.168.0.0/16"
}
```
## Argument Reference

The following arguments are supported:

* `type` - (Optional, Forces new resource) Type of Vpc: 0 - medium(150kpps), 1 - small(100kpps), 2 - large(200kpps), 3 - ultra-large(250kpps), default 1.
* `name` - (Optional) The name of vpc.
* `description`- (Optional) The description of vpc.
* `vpc_network` - (Required, Forces new resource) Network address range of vpc.192.168.0.0/16", "172.16.0.0/16", "172.17.0.0/16",
                                                                               					"172.18.0.0/16", "172.19.0.0/16", "172.20.0.0/16", "172.21.0.0/16", "172.22.0.0/16",
                                                                               					"172.23.0.0/16", "172.24.0.0/16", "172.25.0.0/16"
* `security_group_id` - (Required) security group id , vpc wants to use.   
* `eip_id` - (Optional) eip id , vpc wants to use.                                                         
* `tag_ids` - (Optional) tag ids , vpc wants to use.
## Attributes Reference

The following attributes are exported:

* `type` - Type of Vpc: 0 - medium(150kpps), 1 - small(100kpps), 2 - large(200kpps), 3 - ultra-large(250kpps), default 1.
* `name` - The name of vpc.
* `description`- The description vpc eip.
* `vpc_network` - Network address range of vpc.192.168.0.0/16", "172.16.0.0/16", "172.17.0.0/16",
                                                                               					"172.18.0.0/16", "172.19.0.0/16", "172.20.0.0/16", "172.21.0.0/16", "172.22.0.0/16",
                                                                               					"172.23.0.0/16", "172.24.0.0/16", "172.25.0.0/16"
* `security_group_id` - security group id , vpc wants to use.
* `eip_id` - eip id , vpc wants to use.                                                         
* `tag_ids` - tag ids , vpc wants to use.
* `tag_names` - tag names , computed by tag_ids.
* `private_ip` - The private ip of vpc.
* `public_ip` - The public ip of vpc.
