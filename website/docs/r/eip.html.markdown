---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_eip"
sidebar_current: "docs-qingcloud-resource-eip"
description: |-
  Provides a  EIP resource.
---

# Qingcloud\_eip

Provides a  EIP resource.

Resource can be imported.

## Example Usage

```
# Create a new EIP.
resource "qingcloud_eip" "init"{
        bandwidth = 2
}
```
```
# Create a new EIP with tags.
resource "qingcloud_eip" "foo" {
    name = "eip"
    description = "eip"
    billing_mode = "bandwidth"
    bandwidth = 4
    need_icp = 0
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

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) Maximum bandwidth to the elastic public network, measured in Mbps (Mega bit per second). 
* `billing_mode` - (Optional, Forces new resource) Internet charge type of the EIP : bandwidth , traffic ,default bandwidth.
* `name` - (Optional) The name of eip
* `description`- (Optional) The description of eip
* `need_icp` - (Optional) need icp , 1 need , 0 no need ,default 0
* `tag_ids` - (Optional) tag ids , eip wants to use
## Attributes Reference

The following attributes are exported:

* `id` - The EIP ID.
* `bandwidth` - Maximum bandwidth to the elastic public network, measured in Mbps (Mega bit per second). 
* `billing_mode` - Internet charge type of the EIP : bandwidth , traffic ,default bandwidth.
* `name` - The name of eip.
* `description`- The description of eip.
* `need_icp` - need icp , 1 need , 0 no need ,default 0.
* `tag_ids` - tag ids , eip wants to use.
* `tag_names` - compute by tag ids.
* `status` - The EIP current status.
* `addr` - The elastic ip address.
* `resource` - The resource who use this eip
