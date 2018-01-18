---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_server_certificate"
sidebar_current: "docs-qingcloud-resource-server-certificate"
description: |-
  Provides a Server Certificate resource.
---

# Qingcloud\_server\_certificate

Provides a Server Certificate resource.

Resource can be imported.

## Example Usage

```
# Create a new Server Certificate.
resource "qingcloud_server_certificate" "foo"{
  certificate_content = <<EOF
-----BEGIN CERTIFICATE-----
MIICqzCCAhSgAwIBAgIBATANBgkqhkiG9w0BAQsFADBYMQswCQYDVQQGEwJDTjEO
MAwGA1UECAwFSHViZWkxDjAMBgNVBAcMBVd1aGFuMQ8wDQYDVQQKDAZ5dW5pZnkx
GDAWBgNVBAsMD3d3dy5qdW53dWh1aS5jbjAeFw0xNzEyMTkwMzM0NDZaFw0xODEy
MTkwMzM0NDZaMGIxCzAJBgNVBAYTAkNOMQ4wDAYDVQQIDAVIdWJlaTEPMA0GA1UE
CgwGeXVuaWZ5MRgwFgYDVQQLDA93d3cuanVud3VodWkuY24xGDAWBgNVBAMMD3d3
dy5qdW53dWh1aS5jbjCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAy3Ecw8hi
042R4XPEsBS+C4APC1rexFgg9EQfWyVybO/vRz9cHN4mEiwXpPIMzi1JrH9gsY2W
FJzBhpoOXugTgARMpnCU6iBqUmAPaNj61LDLv+n+6kdA59wlb2pAiWbx91ErsUtz
.........................
-----END CERTIFICATE-----
EOF
  private_key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDLcRzDyGLTjZHhc8SwFL4LgA8LWt7EWCD0RB9bJXJs7+9HP1wc
3iYSLBek8gzOLUmsf2CxjZYUnMGGmg5e6BOABEymcJTqIGpSYA9o2PrUsMu/6f7q
R0Dn3CVvakCJZvH3USuxS3Okt6VolxSkZmCvnwdrRsXdq1h522bVgR8f+QIDAQAB
AoGBAIjYyXSY8oFDlYGGEiQvj7bEqVoGAhso/OHSgRUal2HX86iFYjy44fsPVchK
WXrG0+wIss48Y1vyJeuY7VnB2nrrSkg4+a/je9GNUlow/eiUDeVv5gEwkeccgcsN
eQRxUQLQfsj2a/P6NX3pEodTo2o4plJCj5Vr30Pw3yJk4VRBAkEA6hHHX6GMpqiY
.........................
-----END RSA PRIVATE KEY-----
EOF
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of server certificate.
* `description`- (Optional) The description of server certificate.
* `certificate_content` - (Required if create, Force New) Content of the certificate
* `private_key` - (Required if create, Force New) Private key of the certificate

## Attributes Reference

The following attributes are exported:
* `name` - The name of security group.
* `description`- The description of security group.
* `certificate_content` - (Required if create, Force New) Content of the certificate
* `private_key` - (Required if create, Force New) Private key of the certificate

