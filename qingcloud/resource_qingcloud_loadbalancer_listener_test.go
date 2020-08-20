package qingcloud

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func TestAccQingcloudLoadBalancerListener_basic(t *testing.T) {
	var lbl qc.DescribeLoadBalancerListenersOutput
	testTag := "terraform-test-lb-listener-basic" + os.Getenv("CIRCLE_BUILD_NUM")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "qingcloud_loadbalancer_listener.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckLoadBalancerListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccLBLConfigBasic, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerListenerExists(
						"qingcloud_loadbalancer_listener.foo", &lbl),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "listener_port", "80"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "listener_protocol", "http"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "timeout", "50"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "balance_mode", "roundrobin"),
				),
			},
			{
				Config: fmt.Sprintf(testAccLBLConfigBasicTwo, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerListenerExists(
						"qingcloud_loadbalancer_listener.foo", &lbl),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "listener_port", "443"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "listener_protocol", "https"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "timeout", "50"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "balance_mode", "roundrobin"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "listener_option", "256"),
				),
			},
			{
				Config: fmt.Sprintf(testAccLBLConfigBasicThree, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerListenerExists(
						"qingcloud_loadbalancer_listener.foo", &lbl),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "listener_port", "443"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "listener_protocol", "https"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "timeout", "40"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "balance_mode", "source"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "listener_option", "207"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "healthy_check_method", "http|/index.html|vhost.example.com"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "healthy_check_option", "10|5|5|5"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "forwardfor", "7"),
					resource.TestCheckResourceAttr("qingcloud_loadbalancer_listener.foo", "session_sticky", "prefix|sk"),
				),
			},
		},
	})
}

func testAccCheckLoadBalancerListenerDestroy(s *terraform.State) error {
	return testAccCheckLoadBalancerListenerDestroyWithProvider(s, testAccProvider)
}

func testAccCheckLoadBalancerListenerDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_loadbalancer_listener" {
			continue
		}
		input := new(qc.DescribeLoadBalancerListenersInput)
		input.LoadBalancerListeners = []*string{qc.String(rs.Primary.ID)}
		output, err := client.loadbalancer.DescribeLoadBalancerListeners(input)
		if err == nil {
			if len(output.LoadBalancerListenerSet) != 0 {
				return fmt.Errorf("fount  loadbalancer listener: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccCheckLoadBalancerListenerExists(n string, i *qc.DescribeLoadBalancerListenersOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no loadbalancer listener ID is set")
		}

		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeLoadBalancerListenersInput)
		input.Verbose = qc.Int(1)
		input.LoadBalancerListeners = []*string{qc.String(rs.Primary.ID)}
		d, err := client.loadbalancer.DescribeLoadBalancerListeners(input)

		log.Printf("[WARN] loadbalancer listener id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil || len(d.LoadBalancerListenerSet) == 0 {
			return fmt.Errorf("Lb listener not found ")
		}

		*i = *d
		return nil
	}
}

const testAccLBLConfigBasic = `
resource "qingcloud_eip" "foo" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo.id}"]
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}

resource "qingcloud_loadbalancer_listener" "foo"{
  load_balancer_id = "${qingcloud_loadbalancer.foo.id}"
  listener_port = "80"
  listener_protocol = "http"
}

resource "qingcloud_server_certificate" "foo"{
  certificate_content = <<EOF
-----BEGIN CERTIFICATE-----
MIIDUjCCAjoCCQDbFSEJheoKwTANBgkqhkiG9w0BAQsFADBqMQswCQYDVQQGEwJD
TjEOMAwGA1UECAwFaHViZWkxDjAMBgNVBAcMBXd1aGFuMQ8wDQYDVQQKDAZ5dW5p
ZnkxEjAQBgNVBAsMCWFwcGNlbnRlcjEWMBQGA1UEAwwNcWluZ2Nsb3VkLmNvbTAg
Fw0xOTAxMjQxMjIxNTRaGA8yMTE4MTIzMTEyMjE1NFowajELMAkGA1UEBhMCQ04x
DjAMBgNVBAgMBWh1YmVpMQ4wDAYDVQQHDAV3dWhhbjEPMA0GA1UECgwGeXVuaWZ5
MRIwEAYDVQQLDAlhcHBjZW50ZXIxFjAUBgNVBAMMDXFpbmdjbG91ZC5jb20wggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDT1NzqInpG9v5hstKexG7M3lmw
scJxsZT4a0b/vPpfbc5gskdYSoIZ+DREnKC6YkPwqvoLfJKNAuRHr+Vm0K5Vkj0R
ffjMb6kfgXfwkqNeg42TCbgeqVjdP5z/z5l61RGk0QgXtIKDGdpQXd1VidyoHBY5
AjB4/orYwE85T68WZZLCV3Em290+C4vSABGl/rPft/8ngVfiRicMtFd8FAhj4ibO
g+uNofhQ1bmD/xqHyhnza4V5s8r3T8cXmjw3UzAuJ6TJMweVsCcQIpdJsxs+4/TG
FleXuC2KQ2EXLiDg22v4W6wThjl5ns9nuU/AfbJT3IpEfM33GcFBJU+7AI+TAgMB
AAEwDQYJKoZIhvcNAQELBQADggEBAJPyCKPiwf6VRupzlHwK0UJ3vi/KFGkKFUin
jdw6uftqzYOiapZ06Q2r5BuDbK0YCcr9Vio+GfbifSG87iflaTAUwnqQvQd4dhlO
P+H+wXM28K9BlNP5Kq/pjeUJQp85fBB+pZuhDxb6s5FnudPL1TgchhUDwDcdY7ih
exr1dAOXlBIMoAdVEqtiaBcjidrKCYM4hAh2FHabKjt4L5pgG4yZ/8S0E6hva969
yDIKBxEgEZ3mnDPi8U4TIsbTsuZ5aBM+YNOsB+cAqtrsUuvukNt0VUlkWCJ8c3mD
uXxANM1fZYc/p8GXEcMBPDEbcBsGcOCofNd7CSOIFwXXBaIWDBE=
-----END CERTIFICATE-----
EOF
  private_key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA09Tc6iJ6Rvb+YbLSnsRuzN5ZsLHCcbGU+GtG/7z6X23OYLJH
WEqCGfg0RJygumJD8Kr6C3ySjQLkR6/lZtCuVZI9EX34zG+pH4F38JKjXoONkwm4
HqlY3T+c/8+ZetURpNEIF7SCgxnaUF3dVYncqBwWOQIweP6K2MBPOU+vFmWSwldx
JtvdPguL0gARpf6z37f/J4FX4kYnDLRXfBQIY+ImzoPrjaH4UNW5g/8ah8oZ82uF
ebPK90/HF5o8N1MwLiekyTMHlbAnECKXSbMbPuP0xhZXl7gtikNhFy4g4Ntr+Fus
E4Y5eZ7PZ7lPwH2yU9yKRHzN9xnBQSVPuwCPkwIDAQABAoIBAQC77qq7sjDnirPu
u3au0rk2WsIZx+spcRIoPwyzUNaUGVgyY5h2VUwNfC3q/UZ/dTSvbRD/Zdqi7gDX
NM+CIvu4AVDalvdHcH0L7ZIaRg5YiL/uxn2p/jZPu+Mu9OBGoIfRwH28gjIT27ja
+humivPP1XNFypJ0ledbG2pt/yrn40Dm4MglaqVM8l7dlQBJbBz0Ow8Q+SSVPI9l
/h7UbT9eEwxn0G2gHAPyN+dB5zHKCRm0evocOSM4U7if+4kVoHA+JOlH+H18O49b
ps9aJlqE0VGFZCvv/rJlRPZpaLIc0As6Z7IMk2bxnxOFV68k6uea8CS1toZ16yv+
Zi0uv5NRAoGBAP0n+hV2wM1Nm33MaYBsZs/Z1iThW51k/VjdfBbCv1wJurB3jzXt
3Hx+KwsGWyj+UmRFkOcjIWNH2M4kw5570N2EaTmWhKGAOGVgp6iAj6uHjrbbAAp0
lsoJy2wnZnyUwizcFkTBaam4vPAOSToq28FDdnEt0oG31K/0KIU+zQT5AoGBANY2
C45NnykSrklmaerslY0UwDsBpCJlpPcDzfL3NZSi81MxDl9MOQ0/G3p2TZlONX67
eCGJOOHuUSJCHM625HppA/u4c8oMo8QIA4QqEgUThJtwH4vJMFbetR1wCFT0PwGF
KaHca0eEzLu7buAyPG38Ihm/5mAKcelxXSEK97frAoGAP9coAueqoVtz3dqBtIgh
uULW9P/7yYphNVrNYzQDa+NsN/o+nDv6wU5T7njQ3lqcTnsYmqFKVy4UJ5Av4LSa
rHIq0wH5On0KO86PGTgqgvgxbj12GizipdqoeQLKnpopCYUK/JXF3q4ev27q2oda
Wbd4k/wZPOst8J1i5o86xokCgYAusCDcpzZlcVjjTpsPRPljgn5TXgw0IwtNe1rL
9e2Ls+hs3WhQhQB6TqLikh5fp5gpQxrv1ES9mX+9g8Nbmqk5tOHVX2J9Szv2YfjC
OZkr2hEw/8MgaH1MscWv4NcwDPwejLOCP9RyBhPJZxTHcKuTHT0hd53ymNQzGS4/
IXDUeQKBgAThpKXGclvXVNhP16YR/4j/z31gJA66xg8IoB9ekE09AwHusIpOfKM3
UQHZE8E9gb14oq1779/AUuH5qZo7rs0iDKDoLpG0/Myr+uhHa9tWj1CNZtD17vbG
X5M3++yP6sh1MDbPxcI9Ui/pl/W6/jKoAYeDP7NSJFTcRUp5fzvB
-----END RSA PRIVATE KEY-----
EOF
}

`

const testAccLBLConfigBasicTwo = `
resource "qingcloud_eip" "foo" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo.id}"]
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}

resource "qingcloud_loadbalancer_listener" "foo"{
  load_balancer_id = "${qingcloud_loadbalancer.foo.id}"
  listener_port = "443"
  listener_protocol = "https"
  server_certificate_id = ["${qingcloud_server_certificate.foo.id}"]
  listener_option = 256
}

resource "qingcloud_server_certificate" "foo"{
  certificate_content = <<EOF
-----BEGIN CERTIFICATE-----
MIIDUjCCAjoCCQDbFSEJheoKwTANBgkqhkiG9w0BAQsFADBqMQswCQYDVQQGEwJD
TjEOMAwGA1UECAwFaHViZWkxDjAMBgNVBAcMBXd1aGFuMQ8wDQYDVQQKDAZ5dW5p
ZnkxEjAQBgNVBAsMCWFwcGNlbnRlcjEWMBQGA1UEAwwNcWluZ2Nsb3VkLmNvbTAg
Fw0xOTAxMjQxMjIxNTRaGA8yMTE4MTIzMTEyMjE1NFowajELMAkGA1UEBhMCQ04x
DjAMBgNVBAgMBWh1YmVpMQ4wDAYDVQQHDAV3dWhhbjEPMA0GA1UECgwGeXVuaWZ5
MRIwEAYDVQQLDAlhcHBjZW50ZXIxFjAUBgNVBAMMDXFpbmdjbG91ZC5jb20wggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDT1NzqInpG9v5hstKexG7M3lmw
scJxsZT4a0b/vPpfbc5gskdYSoIZ+DREnKC6YkPwqvoLfJKNAuRHr+Vm0K5Vkj0R
ffjMb6kfgXfwkqNeg42TCbgeqVjdP5z/z5l61RGk0QgXtIKDGdpQXd1VidyoHBY5
AjB4/orYwE85T68WZZLCV3Em290+C4vSABGl/rPft/8ngVfiRicMtFd8FAhj4ibO
g+uNofhQ1bmD/xqHyhnza4V5s8r3T8cXmjw3UzAuJ6TJMweVsCcQIpdJsxs+4/TG
FleXuC2KQ2EXLiDg22v4W6wThjl5ns9nuU/AfbJT3IpEfM33GcFBJU+7AI+TAgMB
AAEwDQYJKoZIhvcNAQELBQADggEBAJPyCKPiwf6VRupzlHwK0UJ3vi/KFGkKFUin
jdw6uftqzYOiapZ06Q2r5BuDbK0YCcr9Vio+GfbifSG87iflaTAUwnqQvQd4dhlO
P+H+wXM28K9BlNP5Kq/pjeUJQp85fBB+pZuhDxb6s5FnudPL1TgchhUDwDcdY7ih
exr1dAOXlBIMoAdVEqtiaBcjidrKCYM4hAh2FHabKjt4L5pgG4yZ/8S0E6hva969
yDIKBxEgEZ3mnDPi8U4TIsbTsuZ5aBM+YNOsB+cAqtrsUuvukNt0VUlkWCJ8c3mD
uXxANM1fZYc/p8GXEcMBPDEbcBsGcOCofNd7CSOIFwXXBaIWDBE=
-----END CERTIFICATE-----
EOF
  private_key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA09Tc6iJ6Rvb+YbLSnsRuzN5ZsLHCcbGU+GtG/7z6X23OYLJH
WEqCGfg0RJygumJD8Kr6C3ySjQLkR6/lZtCuVZI9EX34zG+pH4F38JKjXoONkwm4
HqlY3T+c/8+ZetURpNEIF7SCgxnaUF3dVYncqBwWOQIweP6K2MBPOU+vFmWSwldx
JtvdPguL0gARpf6z37f/J4FX4kYnDLRXfBQIY+ImzoPrjaH4UNW5g/8ah8oZ82uF
ebPK90/HF5o8N1MwLiekyTMHlbAnECKXSbMbPuP0xhZXl7gtikNhFy4g4Ntr+Fus
E4Y5eZ7PZ7lPwH2yU9yKRHzN9xnBQSVPuwCPkwIDAQABAoIBAQC77qq7sjDnirPu
u3au0rk2WsIZx+spcRIoPwyzUNaUGVgyY5h2VUwNfC3q/UZ/dTSvbRD/Zdqi7gDX
NM+CIvu4AVDalvdHcH0L7ZIaRg5YiL/uxn2p/jZPu+Mu9OBGoIfRwH28gjIT27ja
+humivPP1XNFypJ0ledbG2pt/yrn40Dm4MglaqVM8l7dlQBJbBz0Ow8Q+SSVPI9l
/h7UbT9eEwxn0G2gHAPyN+dB5zHKCRm0evocOSM4U7if+4kVoHA+JOlH+H18O49b
ps9aJlqE0VGFZCvv/rJlRPZpaLIc0As6Z7IMk2bxnxOFV68k6uea8CS1toZ16yv+
Zi0uv5NRAoGBAP0n+hV2wM1Nm33MaYBsZs/Z1iThW51k/VjdfBbCv1wJurB3jzXt
3Hx+KwsGWyj+UmRFkOcjIWNH2M4kw5570N2EaTmWhKGAOGVgp6iAj6uHjrbbAAp0
lsoJy2wnZnyUwizcFkTBaam4vPAOSToq28FDdnEt0oG31K/0KIU+zQT5AoGBANY2
C45NnykSrklmaerslY0UwDsBpCJlpPcDzfL3NZSi81MxDl9MOQ0/G3p2TZlONX67
eCGJOOHuUSJCHM625HppA/u4c8oMo8QIA4QqEgUThJtwH4vJMFbetR1wCFT0PwGF
KaHca0eEzLu7buAyPG38Ihm/5mAKcelxXSEK97frAoGAP9coAueqoVtz3dqBtIgh
uULW9P/7yYphNVrNYzQDa+NsN/o+nDv6wU5T7njQ3lqcTnsYmqFKVy4UJ5Av4LSa
rHIq0wH5On0KO86PGTgqgvgxbj12GizipdqoeQLKnpopCYUK/JXF3q4ev27q2oda
Wbd4k/wZPOst8J1i5o86xokCgYAusCDcpzZlcVjjTpsPRPljgn5TXgw0IwtNe1rL
9e2Ls+hs3WhQhQB6TqLikh5fp5gpQxrv1ES9mX+9g8Nbmqk5tOHVX2J9Szv2YfjC
OZkr2hEw/8MgaH1MscWv4NcwDPwejLOCP9RyBhPJZxTHcKuTHT0hd53ymNQzGS4/
IXDUeQKBgAThpKXGclvXVNhP16YR/4j/z31gJA66xg8IoB9ekE09AwHusIpOfKM3
UQHZE8E9gb14oq1779/AUuH5qZo7rs0iDKDoLpG0/Myr+uhHa9tWj1CNZtD17vbG
X5M3++yP6sh1MDbPxcI9Ui/pl/W6/jKoAYeDP7NSJFTcRUp5fzvB
-----END RSA PRIVATE KEY-----
EOF
}

`

const testAccLBLConfigBasicThree = `
resource "qingcloud_eip" "foo" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_loadbalancer" "foo" {
	eip_ids =["${qingcloud_eip.foo.id}"]
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}

resource "qingcloud_loadbalancer_listener" "foo"{
  load_balancer_id = "${qingcloud_loadbalancer.foo.id}"
  listener_port = "443"
  listener_protocol = "https"
  server_certificate_id = ["${qingcloud_server_certificate.foo.id}"]
  listener_option = 207
  healthy_check_method = "http|/index.html|vhost.example.com"
  healthy_check_option = "10|5|5|5"
  timeout = "40"
  forwardfor = "7"
  session_sticky = "prefix|sk"
  balance_mode = "source"
}

resource "qingcloud_server_certificate" "foo"{
  certificate_content = <<EOF
-----BEGIN CERTIFICATE-----
MIIDUjCCAjoCCQDbFSEJheoKwTANBgkqhkiG9w0BAQsFADBqMQswCQYDVQQGEwJD
TjEOMAwGA1UECAwFaHViZWkxDjAMBgNVBAcMBXd1aGFuMQ8wDQYDVQQKDAZ5dW5p
ZnkxEjAQBgNVBAsMCWFwcGNlbnRlcjEWMBQGA1UEAwwNcWluZ2Nsb3VkLmNvbTAg
Fw0xOTAxMjQxMjIxNTRaGA8yMTE4MTIzMTEyMjE1NFowajELMAkGA1UEBhMCQ04x
DjAMBgNVBAgMBWh1YmVpMQ4wDAYDVQQHDAV3dWhhbjEPMA0GA1UECgwGeXVuaWZ5
MRIwEAYDVQQLDAlhcHBjZW50ZXIxFjAUBgNVBAMMDXFpbmdjbG91ZC5jb20wggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDT1NzqInpG9v5hstKexG7M3lmw
scJxsZT4a0b/vPpfbc5gskdYSoIZ+DREnKC6YkPwqvoLfJKNAuRHr+Vm0K5Vkj0R
ffjMb6kfgXfwkqNeg42TCbgeqVjdP5z/z5l61RGk0QgXtIKDGdpQXd1VidyoHBY5
AjB4/orYwE85T68WZZLCV3Em290+C4vSABGl/rPft/8ngVfiRicMtFd8FAhj4ibO
g+uNofhQ1bmD/xqHyhnza4V5s8r3T8cXmjw3UzAuJ6TJMweVsCcQIpdJsxs+4/TG
FleXuC2KQ2EXLiDg22v4W6wThjl5ns9nuU/AfbJT3IpEfM33GcFBJU+7AI+TAgMB
AAEwDQYJKoZIhvcNAQELBQADggEBAJPyCKPiwf6VRupzlHwK0UJ3vi/KFGkKFUin
jdw6uftqzYOiapZ06Q2r5BuDbK0YCcr9Vio+GfbifSG87iflaTAUwnqQvQd4dhlO
P+H+wXM28K9BlNP5Kq/pjeUJQp85fBB+pZuhDxb6s5FnudPL1TgchhUDwDcdY7ih
exr1dAOXlBIMoAdVEqtiaBcjidrKCYM4hAh2FHabKjt4L5pgG4yZ/8S0E6hva969
yDIKBxEgEZ3mnDPi8U4TIsbTsuZ5aBM+YNOsB+cAqtrsUuvukNt0VUlkWCJ8c3mD
uXxANM1fZYc/p8GXEcMBPDEbcBsGcOCofNd7CSOIFwXXBaIWDBE=
-----END CERTIFICATE-----
EOF
  private_key = <<EOF
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA09Tc6iJ6Rvb+YbLSnsRuzN5ZsLHCcbGU+GtG/7z6X23OYLJH
WEqCGfg0RJygumJD8Kr6C3ySjQLkR6/lZtCuVZI9EX34zG+pH4F38JKjXoONkwm4
HqlY3T+c/8+ZetURpNEIF7SCgxnaUF3dVYncqBwWOQIweP6K2MBPOU+vFmWSwldx
JtvdPguL0gARpf6z37f/J4FX4kYnDLRXfBQIY+ImzoPrjaH4UNW5g/8ah8oZ82uF
ebPK90/HF5o8N1MwLiekyTMHlbAnECKXSbMbPuP0xhZXl7gtikNhFy4g4Ntr+Fus
E4Y5eZ7PZ7lPwH2yU9yKRHzN9xnBQSVPuwCPkwIDAQABAoIBAQC77qq7sjDnirPu
u3au0rk2WsIZx+spcRIoPwyzUNaUGVgyY5h2VUwNfC3q/UZ/dTSvbRD/Zdqi7gDX
NM+CIvu4AVDalvdHcH0L7ZIaRg5YiL/uxn2p/jZPu+Mu9OBGoIfRwH28gjIT27ja
+humivPP1XNFypJ0ledbG2pt/yrn40Dm4MglaqVM8l7dlQBJbBz0Ow8Q+SSVPI9l
/h7UbT9eEwxn0G2gHAPyN+dB5zHKCRm0evocOSM4U7if+4kVoHA+JOlH+H18O49b
ps9aJlqE0VGFZCvv/rJlRPZpaLIc0As6Z7IMk2bxnxOFV68k6uea8CS1toZ16yv+
Zi0uv5NRAoGBAP0n+hV2wM1Nm33MaYBsZs/Z1iThW51k/VjdfBbCv1wJurB3jzXt
3Hx+KwsGWyj+UmRFkOcjIWNH2M4kw5570N2EaTmWhKGAOGVgp6iAj6uHjrbbAAp0
lsoJy2wnZnyUwizcFkTBaam4vPAOSToq28FDdnEt0oG31K/0KIU+zQT5AoGBANY2
C45NnykSrklmaerslY0UwDsBpCJlpPcDzfL3NZSi81MxDl9MOQ0/G3p2TZlONX67
eCGJOOHuUSJCHM625HppA/u4c8oMo8QIA4QqEgUThJtwH4vJMFbetR1wCFT0PwGF
KaHca0eEzLu7buAyPG38Ihm/5mAKcelxXSEK97frAoGAP9coAueqoVtz3dqBtIgh
uULW9P/7yYphNVrNYzQDa+NsN/o+nDv6wU5T7njQ3lqcTnsYmqFKVy4UJ5Av4LSa
rHIq0wH5On0KO86PGTgqgvgxbj12GizipdqoeQLKnpopCYUK/JXF3q4ev27q2oda
Wbd4k/wZPOst8J1i5o86xokCgYAusCDcpzZlcVjjTpsPRPljgn5TXgw0IwtNe1rL
9e2Ls+hs3WhQhQB6TqLikh5fp5gpQxrv1ES9mX+9g8Nbmqk5tOHVX2J9Szv2YfjC
OZkr2hEw/8MgaH1MscWv4NcwDPwejLOCP9RyBhPJZxTHcKuTHT0hd53ymNQzGS4/
IXDUeQKBgAThpKXGclvXVNhP16YR/4j/z31gJA66xg8IoB9ekE09AwHusIpOfKM3
UQHZE8E9gb14oq1779/AUuH5qZo7rs0iDKDoLpG0/Myr+uhHa9tWj1CNZtD17vbG
X5M3++yP6sh1MDbPxcI9Ui/pl/W6/jKoAYeDP7NSJFTcRUp5fzvB
-----END RSA PRIVATE KEY-----
EOF
}

`
