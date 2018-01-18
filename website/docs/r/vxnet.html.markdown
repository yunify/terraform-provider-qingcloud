---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_vxnet"
sidebar_current: "docs-qingcloud-resource-vxnet"
description: |-
  Provides a  VxNet resource.
---

# Qingcloud\_vxnet

Provides a  VxNet resource.  

Resource can be imported.

## Example Usage

```
# Create a new VxNet.
resource "qingcloud_vxnet" "foo" {
    type = 1
} 
```
```
# Create a new VxNet with tags.
resource "qingcloud_vxnet" "foo" {
    type = 1
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
# Create a new VxNet with Vpc.
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
```
## Argument Reference

The following arguments are supported:

* `type` - (Required, Forces new resource) type of vxnet,1 - Managed vxnet,0 - Self-managed vxnet..
* `name` - (Optional) The name of vxnet.
* `description`- (Optional) The description vpc vxnet.
* `ip_network` - (Optional) Network address range of vxnet.192.168.x.0/24", "172.16.x.0/24", "172.17.x.0/24",
                                                                               					"172.18.x.0/24", "172.19.x.0/24", "172.20.x.0/24", "172.21.x.0/24", "172.22.x.0/24",
                                                                               					"172.23.x.0/24", "172.24.x.0/24", "172.25.x.0/24"
* `vpc_id` - (Optional) The vpc id , vxnet want to join.                                                  
* `tag_ids` - (Optional) tag ids , vxnet wants to use.
## Attributes Reference

The following attributes are exported:

* `type` - type of vxnet,1 - Managed vxnet,0 - Self-managed vxnet.
* `name` - The name of vxnet.
* `description`- The description vpc vxnet.
* `ip_network` - Network address range of vxnet.192.168.x.0/24", "172.16.x.0/24", "172.17.x.0/24",
                                                                               					"172.18.x.0/24", "172.19.x.0/24", "172.20.x.0/24", "172.21.x.0/24", "172.22.x.0/24",
                                                                               					"172.23.x.0/24", "172.24.x.0/24", "172.25.x.0/24"
* `vpc_id` - The vpc id , vxnet want to join.                                                  
* `tag_ids` - tag ids , vxnet wants to use.
* `tag_names` - tag names , computed by tag_ids.
