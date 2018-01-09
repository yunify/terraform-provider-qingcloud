---
layout: "qingcloud"
page_title: "Qingcloud: qingcloud_loadbalancer_listener"
sidebar_current: "docs-qingcloud_loadbalancer_listener"
description: |-
  Provides a  Loadbalancer Listener resource.
---

# Qingcloud\_loadbalancer\_listener

Provides a  Loadbalancer Listener resource.  

## Example Usage

```
# Create a new Loadbalancer Listener for port 80 's http requests.
resource "qingcloud_eip" "foo" {
    bandwidth = 2
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo.id}"]
}

resource "qingcloud_loadbalancer_listener" "foo"{
  load_balancer_id = "${qingcloud_loadbalancer.foo.id}"
  listener_port = "80"
  listener_protocol = "http"
}
```

## Argument Reference

The following arguments are supported:  

* `load_balancer_id` - (Required, ForceNew) load balancer id .  
* `name` - (Optional) The name of load balancer listener.  
* `listener_port` - (Required, ForceNew) Listening port.  
* `listener_protocol` - (Required, ForceNew) Listening protocol, support `http`, `tcp`, `ssl`, `https`. ssl and https need server_certificate_id.  
* `server_certificate_id` - (Optional) CA Certificate id.   
* `balance_mode` - (Optional, Default:roundrobin) The listener load balancing method , support `roundrobin`, `leastconn`, `source`.   
* `session_sticky` - (Optional) Keep the session.  
                               "Insert: The cookie is inserted by the load balancer, which is specified by the load balancer, and the user only needs to provide the timeout for the cookie."  
                               "Insert: insert|cookie_timeout"  
                               "Example: insert|3600" cookie_timeout can be 0, which means never timeout   
                               "Rewrite: Users to specify and maintain cookies, users need to take the initiative to insert a cookie client, and provide expiration time. The load balancer implements the session hold by overwriting the cookie (prefixing server name with cookie name). When the request is forwarded to the back-end server, the load balancer will take the initiative to delete the server title to achieve cookie transparent to the back-end server."  
                               "Rewrite: prefix|cookie_name"  
                               "Example: prefix|sk"  
* `forwardfor` - (Optional) HTTP headers are required for forwarding requests,The value is a decimal number that is obtained by "bitwise ANDing" with the 3 additional header fields currently supported.  
                            "X-Forwarded-For: bit is 1(Binary 1), Indicates whether the real client IP is passed to the backend. The additional option "Get Client IP" is turned off and the client IP obtained by the backend server is the IP address of the load balancer itself. After turning this feature on, the backend server can get the real user IP through the X-Forwarded-For field in the request."  
                            "QC-LBID :bit is 2(Binary 10), Indicates whether the Header contains the ID of LoadBalancer"   
                            "QC-LBIP: bit is 3(Binary 100), Indicates whether the Header contains the public network IP of LoadBalancer "  
                             For example, if Header contains X-Forwarded-For and QC-LBIP, the forwarfor value is:
                            "X-Forwarded-For | QC-LBIP" with a binary result of 101 and a final conversion of 5 to decimal.  
* `timeout` - (Optional,Default:50) Listener timeout, in seconds.  
* `healthy_check_method` - (Optional,Default:tcp) Listener health check. There are two ways to check HTTP and TCP.  
                                                  TCP: `tcp`,  
                                                  HTTP: `http|url|host` ,example :`http|/index.html` or `http|/index.html|vhost.example.com`  
* `healthy_check_option` - (Optional,Default:10|5|2|5) The listener health check parameter configuration is valid only if the health check is enabled. Format for:   
                                                        inter | timeout | fall | rise  
                                                        Check Interval (2-60s) | Timeout (5-300s) | Unhealthy Threshold (2-10) | Health Threshold (2-10)                                                  
* `listener_option` - (Optional) Additional options. This value is derived from the currently supported 2 additional options in "bitwise AND".  
                                 Cancel URL check: The bit is 1 (binary 1), indicating whether the load balancer can accept URLs that do not conform to encoding specifications, such as URLs that contain unencoderated Chinese characters  
                                 Get client IP: bit is 2 (binary 10), that the client's IP directly to the backend. With this feature on, the load balancer is completely transparent to the backend. The source address of the back-end host TCP connection is the client's IP, not the load balancer's IP. Note: Only backends in managed networks are supported. This feature is not available when using the underlying network backend.  
                                 Data Compression: The bit is 4 (binary 100), indicating whether the gzip algorithm is used to compress the text data to reduce network traffic.  
                                 Disable insecure encryption: bit is 8 (binary 1000), to disable the existence of a security risk of encryption, may not be compatible with the lower version of the client.  
                                 Enabling IE-compatible encryption: bit is 16 (10000 in binary), using the algorithm: AES128 + EECDH: AES128 + EDH no-tls-tickets no-sslv3.  
                                 Enable HTTPS redirection: The bit is 32, redirecting this listener's HTTP request to the HTTPS listener.  
                                 Refuse TCP connection request: bit bit is 64, when all back-end is not available, no longer accept the connection, but directly reject the TCP connection request. Only tcp listeners are supported.  
                                 Disabling the HTTP Header Field Proxy: The bit bit is 128. For security reasons, some web applications with a back-end CGI need to disable this field in user requests to prevent injection attacks.  
                                 Enabling Secure Encryption: bit is 256 using algorithm: ECDHE-ECDSA-AES128-GCM-SHA256: ECDHE-ECDSA-AES256-GCM-SHA384 ECDHE-ECDSA-AES128-SHA: ECDHE-ECDSA-AES256-SHA: ECDHE -ECDSA-AES128-SHA256: ECDHE-ECDSA-AES256-SHA384 ECDHE-RSA-AES128-GCM-SHA256: ECDHE-RSA-AES256-GCM-SHA384: ECDHE-RSA-AES128-SHA: ECDHE-RSA-AES256-SHA ECDHE -RSA-AES128-SHA256: ECDHE-RSA-AES256-SHA384: AES128-GCM-SHA256: AES256-GCM-SHA384: AES128-SHA256 AES256-SHA256: AES128-SHA: AES256-SHA: DES-CBC3-SHA:! DSS no-tls-tickets force-tlsv12 no-sslv3  .
                                 

## Attributes Reference

The following attributes are exported:

* `load_balancer_id` - load balancer id .  
* `name` - The name of load balancer listener.  
* `listener_port` -  Listening port.  
* `listener_protocol` - Listening protocol, support `http`, `tcp`, `ssl`, `https`. ssl and https need server_certificate_id.  
* `server_certificate_id` - CA Certificate id.   
* `balance_mode` - The listener load balancing method , support `roundrobin`, `leastconn`, `source`.   
* `session_sticky` - Keep the session.  
                               "Insert: The cookie is inserted by the load balancer, which is specified by the load balancer, and the user only needs to provide the timeout for the cookie."  
                               "Insert: insert|cookie_timeout"  
                               "Example: insert|3600" cookie_timeout can be 0, which means never timeout   
                               "Rewrite: Users to specify and maintain cookies, users need to take the initiative to insert a cookie client, and provide expiration time. The load balancer implements the session hold by overwriting the cookie (prefixing server name with cookie name). When the request is forwarded to the back-end server, the load balancer will take the initiative to delete the server title to achieve cookie transparent to the back-end server."  
                               "Rewrite: prefix|cookie_name"  
                               "Example: prefix|sk"  
* `forwardfor` - HTTP headers are required for forwarding requests,The value is a decimal number that is obtained by "bitwise ANDing" with the 3 additional header fields currently supported.  
                            "X-Forwarded-For: bit is 1(Binary 1), Indicates whether the real client IP is passed to the backend. The additional option "Get Client IP" is turned off and the client IP obtained by the backend server is the IP address of the load balancer itself. After turning this feature on, the backend server can get the real user IP through the X-Forwarded-For field in the request."  
                            "QC-LBID :bit is 2(Binary 10), Indicates whether the Header contains the ID of LoadBalancer"   
                            "QC-LBIP: bit is 3(Binary 100), Indicates whether the Header contains the public network IP of LoadBalancer "  
                             For example, if Header contains X-Forwarded-For and QC-LBIP, the forwarfor value is:
                            "X-Forwarded-For | QC-LBIP" with a binary result of 101 and a final conversion of 5 to decimal.  
* `timeout` - Listener timeout, in seconds.  
* `healthy_check_method` - Listener health check. There are two ways to check HTTP and TCP.  
                                                  TCP: `tcp`,  
                                                  HTTP: `http|url|host` ,example :`http|/index.html` or `http|/index.html|vhost.example.com`  
* `healthy_check_option` - The listener health check parameter configuration is valid only if the health check is enabled. Format for:   
                                                        inter | timeout | fall | rise  
                                                        Check Interval (2-60s) | Timeout (5-300s) | Unhealthy Threshold (2-10) | Health Threshold (2-10)                                                  
* `listener_option` - Additional options. This value is derived from the currently supported 2 additional options in "bitwise AND".  
                                 Cancel URL check: The bit is 1 (binary 1), indicating whether the load balancer can accept URLs that do not conform to encoding specifications, such as URLs that contain unencoderated Chinese characters  
                                 Get client IP: bit is 2 (binary 10), that the client's IP directly to the backend. With this feature on, the load balancer is completely transparent to the backend. The source address of the back-end host TCP connection is the client's IP, not the load balancer's IP. Note: Only backends in managed networks are supported. This feature is not available when using the underlying network backend.  
                                 Data Compression: The bit is 4 (binary 100), indicating whether the gzip algorithm is used to compress the text data to reduce network traffic.  
                                 Disable insecure encryption: bit is 8 (binary 1000), to disable the existence of a security risk of encryption, may not be compatible with the lower version of the client.  
                                 Enabling IE-compatible encryption: bit is 16 (10000 in binary), using the algorithm: AES128 + EECDH: AES128 + EDH no-tls-tickets no-sslv3.  
                                 Enable HTTPS redirection: The bit is 32, redirecting this listener's HTTP request to the HTTPS listener.  
                                 Refuse TCP connection request: bit bit is 64, when all back-end is not available, no longer accept the connection, but directly reject the TCP connection request. Only tcp listeners are supported.  
                                 Disabling the HTTP Header Field Proxy: The bit bit is 128. For security reasons, some web applications with a back-end CGI need to disable this field in user requests to prevent injection attacks.  
                                 Enabling Secure Encryption: bit is 256 using algorithm: ECDHE-ECDSA-AES128-GCM-SHA256: ECDHE-ECDSA-AES256-GCM-SHA384 ECDHE-ECDSA-AES128-SHA: ECDHE-ECDSA-AES256-SHA: ECDHE -ECDSA-AES128-SHA256: ECDHE-ECDSA-AES256-SHA384 ECDHE-RSA-AES128-GCM-SHA256: ECDHE-RSA-AES256-GCM-SHA384: ECDHE-RSA-AES128-SHA: ECDHE-RSA-AES256-SHA ECDHE -RSA-AES128-SHA256: ECDHE-RSA-AES256-SHA384: AES128-GCM-SHA256: AES256-GCM-SHA384: AES128-SHA256 AES256-SHA256: AES128-SHA: AES256-SHA: DES-CBC3-SHA:! DSS no-tls-tickets force-tlsv12 no-sslv3  .
                                 

