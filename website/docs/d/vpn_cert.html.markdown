---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_vpn_cert"
sidebar_current: "docs-qingcloud-datasource-vpn-cert"
description: |-
    Provides vpn cert of vpc/router..
---

# qingcloud\_vpn\_cert

The vpn cert data source get the certs info of vpc/router

## Example Usage

```
data "qingcloud_vpn_cert" "test"{
  router_id = "xxx"
  platform = "linux"
}


```

## Argument Reference

The following arguments are supported:

* `router_id` - (Required) vpc/router id. 
* `platform` - (Optional, Default:linux)Validate value:"linux","mac","windows"  VPN's conf sample's platform.

## Attributes Reference


* `router_id` - vpc/router id. 
* `platform` - Validate value:"linux","mac","windows"  VPN's conf sample's platform.
* `client_crt` - Client certificate content.
* `client_key` - The client certificate private key content.
* `static_key` - vpn encrypt private key content.
* `ca_cert` - vpn trust certificate content.
* `conf_sample` - Sample configuration file content.
