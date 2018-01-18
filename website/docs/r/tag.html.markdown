---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_tag"
sidebar_current: "docs-qingcloud-resource-tag"
description: |-
  Provides a  Tag resource.
---

# Qingcloud\_tag

Provides a  TAG resource.

Resource can be imported.

## Example Usage

```
# Create a new Tag.
resource "qingcloud_tag" "foo"{
	name="tag1"
}
```
```
# Create a new tag with color.
resource "qingcloud_tag" "foo"{
	name="tag1"
	description="test"
	color = "#ffffff"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of tag.
* `color` - (Optional) The color of tag.
* `description`- (Optional) The description of tag.

## Attributes Reference

The following attributes are exported:
* `name` - The name of tag.
* `color` - The color of tag.
* `description`- The description of tag.
