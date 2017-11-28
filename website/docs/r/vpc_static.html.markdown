---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_vpc_static"
sidebar_current: "docs-qingcloud-resource-vpc-static"
description: |-
  Provides a  VpcStatic resource.
---

# Qingcloud\_vpc\_static

Provides a  VpcStatic resource.  
Use this resource in SDN 2.0 environment.

## Example Usage

```
# Create a new VPC Static for port forward.
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
}
resource "qingcloud_vpc_static" "foo"{
        vpc_id = "${qingcloud_vpc.foo.id}"
        type = 1
        val1 = "80"
        val2 = "192.168.0.3"
        val3 = "81"
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required , ForceNew) Type of VpcStatic:  1 : port_forwarding , 2 : VPN rule , 3 : DHCP , 4 :  Two layers GRE , 6 :  Three layers GRE , 7 :  Three layers IPsec , 8 :  Private DNS .
* `name` - (Optional) The name of vpc static.
* `vpc_id`- (Required , ForceNew) The id of vpc.
* `val1` - (Required) "port_forwarding : source port "  
                "VPN : type of vpn ,'openvpn', 'pptp','l2tp', default 'openvpn' "
                "DHCP : id of DHCP host"   
                "Two layers GRE : remote ip , secret key, example : gre|1.2.3.4|888	"   
                "Three layers GRE: remote ip , secret key, local p2p ip , opposite end p2p ip , example : 6.6.6.6|key|1.2.3.4|4.3.2.1 "   
                "Three layers IPsec : remote ip(support 0.0.0.0 for any) ; encryption method :phase2alg&ike , default aes ; secret key & remote device id"   
                "Private DNS : private domain name"
* `val2` - (Optional) "port_forwarding : destination ip "   
                      	"OpenVPN : VPN Server Port , default 1194"   
                      	"PPTP/L2TP : username & password , format (user:password)"   
          				"DHCP : DHCP  configuration content "   
          				"Three layers GRE: destination network , multiple networks are separated by '|' "   
                      	"Three layers IPsec : local network , multiple networks are separated by '|' "   
                      	"Private DNS : IP address ,192.168.1.2;192.168.1.3"
* `val3` - (Optional) "port_forwarding : destination port "   
           	"OpenVPN : VPN protocol , default udp"   
           	"PPTP VPN : Max Connections , 1-253"   
           	"L2TP VPN :(PSK, preshared secrets) "   
           	"Three layers IPsec : destination network , multiple networks are separated by '|' "                                                     
* `val4` - (Optional) "port_forwarding : protocol , default tcp , support udp & tcp "   
                      	"VPN : client CIDR ,support 10.255.x.0/24 , default auto allocation"   
                      	"Three layers IPsec : tunnel pattern . default main , support main & aggrmode ".
* `val5` - (Optional) "OpenVPN :  verification mode "   
                      	"L2TP VPN : L2TP server port 1701".
## Attributes Reference

The following attributes are exported:

* `type` - Type of VpcStatic:  1 : port_forwarding , 2 : VPN rule , 3 : DHCP , 4 :  Two layers GRE , 6 :  Three layers GRE , 7 :  Three layers IPsec , 8 :  Private DNS .
* `name` - The name of vpc static.
* `vpc_id`- The id of vpc.
* `val1` - "port_forwarding : source port "  
                "VPN : type of vpn ,'openvpn', 'pptp','l2tp', default 'openvpn' "
                "DHCP : id of DHCP host"   
                "Two layers GRE : remote ip , secret key, example : gre|1.2.3.4|888	"   
                "Three layers GRE: remote ip , secret key, local p2p ip , opposite end p2p ip , example : 6.6.6.6|key|1.2.3.4|4.3.2.1 "   
                "Three layers IPsec : remote ip(support 0.0.0.0 for any) ; encryption method :phase2alg&ike , default aes ; secret key & remote device id"   
                "Private DNS : private domain name"
* `val2` - "port_forwarding : destination ip "   
                      	"OpenVPN : VPN Server Port , default 1194"   
                      	"PPTP/L2TP : username & password , format (user:password)"   
          				"DHCP : DHCP  configuration content "   
          				"Three layers GRE: destination network , multiple networks are separated by '|' "   
                      	"Three layers IPsec : local network , multiple networks are separated by '|' "   
                      	"Private DNS : IP address ,192.168.1.2;192.168.1.3"
* `val3` - "port_forwarding : destination port "   
           	"OpenVPN : VPN protocol , default udp"   
           	"PPTP VPN : Max Connections , 1-253"   
           	"L2TP VPN :(PSK, preshared secrets) "   
           	"Three layers IPsec : destination network , multiple networks are separated by '|' "                                                     
* `val4` - "port_forwarding : protocol , default tcp , support udp & tcp "   
                      	"VPN : client CIDR ,support 10.255.x.0/24 , default auto allocation"   
                      	"Three layers IPsec : tunnel pattern . default main , support main & aggrmode ".
* `val5` - "OpenVPN :  verification mode "   
                      	"L2TP VPN : L2TP server port 1701".
