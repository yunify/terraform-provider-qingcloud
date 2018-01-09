---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_loadbalancer_backend"
sidebar_current: "docs-qingcloud_loadbalancer_backend"
description: |-
  Provides a  Loadbalancer Backend resource.
---

# Qingcloud\_loadbalancer\_backend

Provides a  Loadbalancer Backend resource.  

## Example Usage

```
# Create a new Loadbalancer backend.
resource "qingcloud_loadbalancer" "foo" {
  eip_ids =["${qingcloud_eip.foo.id}"]
}

resource "qingcloud_loadbalancer_listener" "foo"{
  load_balancer_id = "${qingcloud_loadbalancer.foo.id}"
  listener_port = "80"
  listener_protocol = "http"
}
resource "qingcloud_instance" "foo" {
  image_id = "centos7x64d"
  keypair_ids = ["${qingcloud_keypair.foo.id}"]
}
resource "qingcloud_loadbalancer_backend" "foo" {
  loadbalancer_listener_id = "${qingcloud_loadbalancer_listener.foo.id}"
  port = 80
  resource_id = "${qingcloud_instance.foo.id}"
}
resource "qingcloud_keypair" "foo"{
  public_key = "    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCyLSPqVIdXGH0QlGeWcPwa1fjTRKl6WtMiaSsP8/GnwjakDSKILUCoNe1yIpiK8F0/gmL71xaDQyfl7k6aE+gn6lSLUjpDmucAF1luGg6l7CIN+6hCqY3YqlAI05Tqwu0PdLAwCbGwdHcaWfECcbROJk5D0zpCTHmissrrAxdOv72g9Ple8KJ6C7F1tz6wmG0zUeineguGjW/PvfZiBDWZ/CyXGPeMDJxv3lrIiLa/ShgnQOxFTdHJPCw+F0/XlSzlIzP3gfni1vXxJWvYjdE9ULo7Z1DLWgZ73FCbeAvX/0e9C9jwT21Qa5RUy4pSP8m4WXSJgw2f9IpY1vIJFSZP root@centos1    "
}
```

## Argument Reference

The following arguments are supported:  

* `loadbalancer_listener_id` - (Required, ForceNew) load balancer listener id .  
* `name` - (Optional) The name of load balancer backend.  
* `port` - (Required) Backend server port.  
* `resource_id` - (Required, ForceNew)  Backend instance id
* `weight` - (Optional, Default 1)Range 1-100. Backend service weight.   

## Attributes Reference

The following attributes are exported:

* `loadbalancer_listener_id` - load balancer listener id .  
* `name` - The name of load balancer backend.  
* `port` - Backend server port.  
* `resource_id` - Backend instance id
* `weight` - Backend service weight.   
