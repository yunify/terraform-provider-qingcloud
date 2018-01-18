---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_volume"
sidebar_current: "docs-qingcloud-resource-volume"
description: |-
  Provides a  Volume resource.
---

# Qingcloud\_volume

Provides a  Volume resource.  

Resource can be imported.

## Example Usage

```
# Create a new Volume.
resource "qingcloud_volume" "foo"{
	size = 10
}
```
```
# Create a new Volume with tags.
resource "qingcloud_volume" "foo"{
	size = 10
	tag_ids = ["${qingcloud_tag.test.id}",
				"${qingcloud_tag.test2.id}"]
}
resource "qingcloud_tag" "test"{
	name="1"
}
resource "qingcloud_tag" "test2"{
	name="2"
}
```
## Argument Reference

The following arguments are supported:

* `type` - (Required, Forces new resource, Default 0) performance type volume is 0;
                                           					Ultra high performance type volume is 3 (only attach to ultra high performance type instance);
                                           					Capacity type volume ,the values vary from region to region , Some region are 1 and some are 2.
* `name` - (Optional) The name of volume.
* `description`- (Optional) The description vpc volume.
* `size` - (Required) size of volume,measured in gb ,min 10 ,max 5000 ,multiples of 10.                                                  
* `tag_ids` - (Optional) tag ids , volume wants to use.
## Attributes Reference

The following attributes are exported:

* `type` - performance type volume is 0;
           Ultra high performance type volume is 3 (only attach to ultra high performance type instance);
           Capacity type volume ,the values vary from region to region , Some region are 1 and some are 2.
* `name` - The name of volume.
* `description`- The description vpc volume.
* `size` - size of volume,measured in gb ,min 10 ,max 5000 ,multiples of 10.                                                  
* `tag_ids` - tag ids , volume wants to use.
* `tag_names` - tag names , computed by tag_ids.
