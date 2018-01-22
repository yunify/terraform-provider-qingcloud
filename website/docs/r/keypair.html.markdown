---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_keypair"
sidebar_current: "docs-qingcloud-resource-eip"
description: |-
  Provides a  Keypair resource.
---

# Qingcloud\_keypair

Uploads an SSH public key.

Resource can be imported.

## Example Usage

```
# Create a new Keypair.
resource "qingcloud_keypair" "foo"{
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
```
```
# Create a new Keypair with tags.
resource "qingcloud_keypair" "foo" {
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
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

* `public_key` - (Required) The SSH public key.
* `name` - (Optional) The name of keypair.
* `description`- (Optional) The description of keypair.
* `tag_ids` - (Optional) tag ids , keypair wants to use.

## Attributes Reference

The following attributes are exported:

* `id` - The Keypair ID.
* `public_key` - The SSH public key.
* `name` - The name of keypair.
* `description`- The description of keypair.
* `tag_ids` - tag ids , keypair wants to use.
* `tag_names` - compute by tag ids.

