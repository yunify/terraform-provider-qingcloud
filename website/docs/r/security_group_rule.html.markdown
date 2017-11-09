---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_security_group_rule"
sidebar_current: "docs-qingcloud-resource-security-group-rule"
description: |-
  Provides a Security Group Rule resource.
---

# Qingcloud\_security\_group\_rule

Provides a Security Group Rule resource.

## Example Usage

```
# Create a new security group.
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
}
# Add a security group rule
resource "qingcloud_security_group_rule" "foo"{
    security_group_id= "${qingcloud_security_group.foo.id}"
    protocol = "tcp"
    priority = 0
    action = "accept"
    direction = 0
    from_port = 0
    to_port = 0
}
# Add another security group rule
resource "qingcloud_security_group_rule" "foo1"{
    security_group_id= "${qingcloud_security_group.foo.id}"
    protocol = "udp"
    priority = 1
    action = "drop"
    direction = 1
    from_port = 10
    to_port = 20
}
```
## Argument Reference

The following arguments are supported:

* `security_group_id` - (Required) The id of security group.
* `name`- (Optional) The name of security group rule.
* `protocol` - (Optional) "tcp", "udp", "icmp", "gre", "esp", "ah", "ipip".
* `priority`- (Optional) From high to low 0 - 100 , default 0.
* `action`- (Required) accept/drop.
* `direction`- (Optional) 0 express down ,1 express up.default 0 .
* `from_port`- (Optional) If protocol is tcp or udp,this value is start port. else if protocol is icmp,this value is the type of ICMP. The others protocol don't need this value.
* `to_port`- (Optional) If protocol is tcp or udp,this value is end port. else if protocol is icmp,this value is the code of ICMP. the others protocol don't need this value.
* `cidr_block`- (Optional) target IP,the Security Group Rule only affect to those IPs .

## Attributes Reference

The following attributes are exported:
* `security_group_id` - The id of security group.
* `name`- The name of security group rule.
* `protocol` - "tcp", "udp", "icmp", "gre", "esp", "ah", "ipip".
* `priority`- From high to low 0 - 100 , default 0.
* `action`- accept/drop.
* `direction`- 0 express down ,1 express up.default 0 .
* `from_port`- If protocol is tcp or udp,this value is start port. else if protocol is icmp,this value is the type of ICMP. The others protocol don't need this value.
* `to_port`- If protocol is tcp or udp,this value is end port. else if protocol is icmp,this value is the code of ICMP. the others protocol don't need this value.
* `cidr_block`- target IP,the Security Group Rule only affect to those IPs .
