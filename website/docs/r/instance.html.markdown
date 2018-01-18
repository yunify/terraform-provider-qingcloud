---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_instance"
sidebar_current: "docs-qingcloud-resource-instance"
description: |-
  Provides a  Instance resource.
---

# Qingcloud\_instance

Provides a  Instance resource.

Resource can be imported.

## Example Usage

```hcl
# Create a new Instance.
resource "qingcloud_keypair" "foo"{
	public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
resource "qingcloud_instance" "foo" {
	image_id = "centos7x64d"
	keypair_ids = ["${qingcloud_keypair.foo.id}"]
}
```
```hcl
# Create a new Instance with tags.
resource "qingcloud_keypair" "foo"{
	public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
resource "qingcloud_instance" "foo" {
	image_id = "centos7x64d"
	keypair_ids = ["${qingcloud_keypair.foo.id}"]
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
```hcl
# Create a new Instance with keypairs.
resource "qingcloud_keypair" "foo1"{
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
resource "qingcloud_keypair" "foo2"{
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
resource "qingcloud_keypair" "foo3"{
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
resource "qingcloud_instance" "foo" {
	image_id = "centos7x64d"
	keypair_ids = ["${qingcloud_keypair.foo1.id}","${qingcloud_keypair.foo2.id}","${qingcloud_keypair.foo3.id}"]
}
```
```hcl
# Create a new Instance with volumes
resource "qingcloud_keypair" "foo"{
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
resource "qingcloud_volume" "foo1"{
	size = 10
}
resource "qingcloud_volume" "foo2"{
	size = 10
}
resource "qingcloud_volume" "foo3"{
	size = 10
}
resource "qingcloud_instance" "foo" {
	image_id = "centos7x64d"
	volume_ids = ["${qingcloud_volume.foo1.id}","${qingcloud_volume.foo2.id}","${qingcloud_volume.foo3.id}"]
	keypair_ids = ["${qingcloud_keypair.foo.id}"]
}
```
```hcl
# Create a new Instance with eip
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
}
resource "qingcloud_eip" "foo" {
    bandwidth = 2
}
resource "qingcloud_keypair" "foo"{
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
resource "qingcloud_instance" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	image_id = "centos7x64d"
	eip_id = "${qingcloud_eip.foo.id}"
	keypair_ids = ["${qingcloud_keypair.foo.id}"]
}
```
```hcl
# Create a new Instance with userdata
resource "qingcloud_instance" "foo" {
  image_id         = "centos73x64"
  keypair_ids      = ["${qingcloud_keypair.foo.id}"]
  userdata = "${base64encode(file("./hello.zip"))}"
}
resource "qingcloud_keypair" "foo"{
	public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required) The Image to use for the instance . 
* `name` - (Optional) The name of instance.
* `description`- (Optional) The description of instance.
* `cpu` - (Optional) cpu of instance , effective value(core) :1, 2, 4, 8, 16.
* `memory` - (Optional) memory of instance , effective value(GB) :1024, 2048, 4096, 6144, 8192, 12288, 16384, 24576, 32768.
* `instance_class` - (Optional) Type of instance , 0 - Performance type , 1 - Ultra high performance type.
* `managed_vxnet_id` - (Optional) managed vxnet id , instance wants to use.
* `private_ip` - (Optional) the ip in managed_vxnet_id.
* `keypair_ids` - (Required) List of keypair ids , instance wants to use.
* `security_group_id` - (Optional) security group id , instance wants to use.
* `eip_id` - (Optional) eip id , instance wants to use.
* `volume_ids` - (Optional) List of volume ids , instance wants to use.
* `userdata` - (Optional, ForceNew)Maximum 2M  Upload an archive(zip,tar,tgz,tbz), it would be extract in `/` , need base64-encode. 

## Attributes Reference

The following attributes are exported:

* `image_id` - The Image to use for the instance . 
* `name` - The name of instance.
* `description`- The description of instance.
* `cpu` - cpu of instance , effective value(core) :1, 2, 4, 8, 16.
* `memory` - memory of instance , effective value(GB) :1024, 2048, 4096, 6144, 8192, 12288, 16384, 24576, 32768.
* `instance_class` - Type of instance , 0 - Performance type , 1 - Ultra high performance type.
* `managed_vxnet_id` - managed vxnet id , instance wants to use.
* `private_ip` - the ip in managed_vxnet_id.
* `keypair_ids` - List of keypair ids , instance wants to use.
* `security_group_id` - security group id , instance wants to use.
* `eip_id` - eip id , instance wants to use.
* `volume_ids` - List of volume ids , instance wants to use.
* `public_ip` - Public ip of instance.
* `userdata` - An archive's base64 code.
