---
layout: "qingcloud"
page_title: "Provider: qingcloud"
sidebar_current: "docs-qingcloud-index"
description: |-
  The Qingcloud provider is used to interact with many resources supported by Qingcloud. The provider needs to be configured with the proper credentials before it can be used..
---

# Qingcloud Provider

The Qingcloud provider is used to interact with the
many resources supported by [Qingcloud](https://www.qingcloud.com). The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Qingcloud Provider
provider "Qingcloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.zone}"
}
# Create a eip
resource "qingcloud_eip" "init"{
	name = "连接第一个主机的地址"
	description = "主机-1"
	billing_mode = "traffic"
	bandwidth = 2
	need_icp = 0
}

# Create a keypair
resource "qingcloud_keypair" "arthur"{
	name = "arthur"
	description = "sdfafd"
	public_key = "${file("~/.ssh/id_rsa.pub")}"
}
# Create a web server
resource "qingcloud_instance" "init"{
	count = 1
	name = "master-${count.index}"
	image_id = "centos7x64d"
	instance_class = "0"
	keypair_ids = ["${qingcloud_keypair.arthur.id}"]
	security_group_id ="${qingcloud_security_group.test.id}"
	eip_id = "${qingcloud_eip.init.id}"
    vxnet_id="vxnet-0"
}
# Create security group
resource "qingcloud_security_group" "test"{
        name = "testsg"
}
```
## Authentication

The Qingcloud provider offers a flexible means of providing credentials for authentication.
The following methods are supported, in this order, and explained below:

- Static credentials
- Environment variables

### Static credentials ###

Static credentials can be provided by adding an `access_key` `secret_key` `zone` and `endpoint` in-line in the
qingcloud provider block:

Usage:

```hcl
provider "qingcloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
  endpoint   = "${var.endpoint}"
}
```


###Environment variables

You can provide your credentials via `QINGCLOUD_ACCESS_KEY` and `QINGCLOUD_SECRET_KEY`,
environment variables, representing your Qingcloud Access Key and Secret Key, respectively.
`QINGCLOUD_ZONE` and `QINGCLOUD_ENDPOINT` is also used, if applicable:

```hcl
provider "qingcloud" {}
```

Usage:

```shell
$ export QINGCLOUD_ACCESS_KEY="anaccesskey"
$ export QINGCLOUD_SECRET_KEY="asecretkey"
$ export QINGCLOUD_ZONE="pek3a"
$ export QINGCLOUD_ENDPOINT="https://api.qingcloud.com:443/iaas"
$ terraform plan
```


## Argument Reference

The following arguments are supported:

* `access_key` - (Optional) This is the Qingcloud access key. It must be provided,
but it can also be sourced from the `QINGCLOUD_ACCESS_KEY` environment variable.
In old version access_key name is id , and it was deprecated.

* `secret_key` - (Optional) This is the Qingcloud secret key. It must be provided, but
it can also be sourced from the `QINGCLOUD_SECRET_KEY` environment variable.
In old version secret_key name is secret , and it was deprecated.

* `endpoint` - (Optional) This is the Qingcloud API address , default `https://api.qingcloud.com:443/iaas`.  
It can also be sourced from the  `QINGCLOUD_ENDPOINT` . This parameter is often used in private clouds.

* `zone` - (Required) This is the Qingcloud zone. It must be provided, but
it can also be sourced from the `QINGCLOUD_ZONE` environment variables.


## Testing

Credentials must be provided via the `QINGCLOUD_ACCESS_KEY`, and `QINGCLOUD_SECRET_KEY` environment variables in order to run acceptance tests.
