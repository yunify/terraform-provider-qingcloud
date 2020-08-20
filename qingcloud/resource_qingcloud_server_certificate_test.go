package qingcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func TestAccQingcloudServerCertificate_basic(t *testing.T) {
	var cert qc.DescribeServerCertificatesOutput

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "qingcloud_server_certificate.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckServerCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerCertificateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerCertificateExists("qingcloud_server_certificate.foo", &cert),
				),
			},
			{
				Config: testAccServerCertificateConfigTwo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerCertificateExists("qingcloud_server_certificate.foo", &cert),
					resource.TestCheckResourceAttr(
						"qingcloud_server_certificate.foo", resourceName, "test"),
					resource.TestCheckResourceAttr(
						"qingcloud_server_certificate.foo", resourceDescription, "test"),
				),
			},
		},
	})
}

func testAccCheckServerCertificateExists(n string, cert *qc.DescribeServerCertificatesOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Server Certificate ID is set")
		}

		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeServerCertificatesInput)
		input.ServerCertificates = []*string{qc.String(rs.Primary.ID)}
		d, err := client.loadbalancer.DescribeServerCertificates(input)

		log.Printf("[WARN] server certificate id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil || qc.StringValue(d.ServerCertificateSet[0].ServerCertificateID) == "" {
			return fmt.Errorf("server certificate not found")
		}

		*cert = *d
		return nil
	}
}

func testAccCheckServerCertificateDestroy(s *terraform.State) error {
	return testAccCheckServerCertificateDestroyWithProvider(s, testAccProvider)
}

func testAccCheckServerCertificateDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_server_certificate" {
			continue
		}

		// Try to find the resource
		input := new(qc.DescribeServerCertificatesInput)
		input.ServerCertificates = []*string{qc.String(rs.Primary.ID)}
		output, err := client.loadbalancer.DescribeServerCertificates(input)
		if err == nil {
			if len(output.ServerCertificateSet) != 0 {
				return fmt.Errorf("Found  ServerCertificate: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAccServerCertificateConfig = `
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
const testAccServerCertificateConfigTwo = `
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
  name = "test"
  description = "test"
}
`
