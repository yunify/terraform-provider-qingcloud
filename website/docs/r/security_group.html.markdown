---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_security_group"
sidebar_current: "docs-qingcloud-resource-security-group"
description: |-
  Provides a Security Group resource.
---

# Qingcloud\_security\_group

Provides a Security Group resource.

## Example Usage

```
# Create a new security group.
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
}
```
```
# Create a new security group with tags.
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
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

* `name` - (Required) The name of security group.
* `description`- (Optional) The description of security group.
* `tag_ids` - (Optional) tag ids , security group wants to use

## Attributes Reference

The following attributes are exported:
* `name` - The name of security group.
* `description`- The description of security group.
* `tag_ids` - tag ids , security group wants to use
* `tag_names` - compute by tag ids.
