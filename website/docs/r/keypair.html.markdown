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

For use this resource, user needs to create SSH key in local first, then put the public key's content into public_key to create qingcloud_keypair.  

## Creating an SSH key on Windows

The simplest way to create SSH key on Windows is to use [PuTTYgen](https://www.chiark.greenend.org.uk/~sgtatham/putty/latest.html).

* Download and run PuTTYgen.
* Click the "Generate" button.
* For additional security, you can enter a key passphrase. This will be required to use the SSH key, and will prevent someone with access to your key file from using the key.
* Once the key has been generated, click "Save Private Key". Make sure you save this somewhere safe, as it is not possible to recover this file if it gets lost
* Select all of the text in the "Public key for pasting into OpenSSH `authorized_keys file`, and put it into public_key

## Creating an SSH key on Linux or Mac

* Run: `ssh-keygen -t rsa`. 
* Press enter when asked where you want to save the key (this will use the default location).
* Enter a passphrase for your key.
* Run `cat ~/.ssh/id_rsa.pub` - this will give you the key in the proper format to paste into public_key.
* Make sure you backup the `~/.ssh/id_rsa file`. This cannot be recovered if it is lost.

## Connecting to a server using an SSH key from a Windows client

* Download and run the PuTTY SSH client.
* Type the IP address or Username + IP address ( `user@x.x.x.x` ) of the destination server under the "Host Name" field on the "Session" category.
* Navigate to the "Connection -> SSH -> Auth" category (left-hand side).
* Click "Browse..." near "Private key file for authentication". Choose the private key file (ending in `.ppk`) that you generated earlier with PuTTYgen.
* Click "Open" to initiate the connection.
* When finished, end your session by pressing `Ctrl+d`.

## Connecting to a server using an SSH key from a Linux or Mac client

* Check that your Linux operating system has an SSH client installed ( which ssh ). If a client is not installed, you will need to install one.
* nitiate a connection: `ssh -i /path/to/id_rsa user@x.x.x.x`
* When finished, end your session by pressing `Ctrl+d`.


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

