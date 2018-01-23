# An example for create sn instance that can be ssh in a public network.

For user security and ease of use, Terraform-Qingcloud only support keypair login mode.  
And the resource keypair is upload public_key.  
User need to generate SSH keys locally.

## How to generate SSH keys

Reference to this [doc](https://github.com/yunify/terraform-provider-qingcloud/blob/master/website/docs/r/keypair.html.markdown)

## Create resource keypair

After generating the SSH key, we need use the public_key to create resource keypair, in this example, config is

```hcl
resource "qingcloud_keypair" "arthur"{
  name = "arthur"
  public_key = "${file("~/.ssh/id_rsa.pub")}"
}
```
And in this config ,we use Terraform's build-in func file,to get `~/.ssh/id_rsa.pub`'s content, or you can just put the content into it just like:
```hcl

resource "qingcloud_keypair" "foo"{
  name = "arthur"
  public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
```

To get more info about Terraform build-in func. reference to [official document](https://www.terraform.io/docs/configuration/interpolation.html#built-in-functions
)
## Associate the instance and keypair

We could associate the instance and keypair in .tf config file.  
Each resource maybe have different ways to associate,It depends on the provider's implementation.  
In this config file, ` keypair_ids = ["${qingcloud_keypair.arthur.id}"]` is the point.  
```hcl
resource "qingcloud_keypair" "arthur"{
  name = "arthur"
  public_key = "${file("~/.ssh/id_rsa.pub")}"
}

resource "qingcloud_instance" "init"{
  count = 1
  name = "master-${count.index}"
  image_id = "centos7x64d"
  instance_class = "0"
  managed_vxnet_id="vxnet-0"
  keypair_ids = ["${qingcloud_keypair.arthur.id}"]
  security_group_id ="${qingcloud_security_group.basic.id}"
  eip_id = "${qingcloud_eip.init.id}"
}
```

## Use SSH key to login

Reference to this [doc](https://github.com/yunify/terraform-provider-qingcloud/blob/master/website/docs/r/keypair.html.markdown)
