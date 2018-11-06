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
			resource.TestStep{
				Config: testAccServerCertificateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServerCertificateExists("qingcloud_server_certificate.foo", &cert),
				),
			},
			resource.TestStep{
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
MIICqzCCAhSgAwIBAgIBATANBgkqhkiG9w0BAQsFADBYMQswCQYDVQQGEwJDTjEO
MAwGA1UECAwFSHViZWkxDjAMBgNVBAcMBVd1aGFuMQ8wDQYDVQQKDAZ5dW5pZnkx
GDAWBgNVBAsMD3d3dy5qdW53dWh1aS5jbjAeFw0xNzEyMTkwMzM0NDZaFw0xODEy
MTkwMzM0NDZaMGIxCzAJBgNVBAYTAkNOMQ4wDAYDVQQIDAVIdWJlaTEPMA0GA1UE
CgwGeXVuaWZ5MRgwFgYDVQQLDA93d3cuanVud3VodWkuY24xGDAWBgNVBAMMD3d3
dy5qdW53dWh1aS5jbjCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEAy3Ecw8hi
042R4XPEsBS+C4APC1rexFgg9EQfWyVybO/vRz9cHN4mEiwXpPIMzi1JrH9gsY2W
FJzBhpoOXugTgARMpnCU6iBqUmAPaNj61LDLv+n+6kdA59wlb2pAiWbx91ErsUtz
pLelaJcUpGZgr58Ha0bF3atYedtm1YEfH/kCAwEAAaN7MHkwCQYDVR0TBAIwADAs
BglghkgBhvhCAQ0EHxYdT3BlblNTTCBHZW5lcmF0ZWQgQ2VydGlmaWNhdGUwHQYD
VR0OBBYEFGO4AhAHP9vucMikI+xeh6/OGiI5MB8GA1UdIwQYMBaAFPwY6mMxV/Dv
1tJCusIXdZUkOdf2MA0GCSqGSIb3DQEBCwUAA4GBAHL14KnNjk6ZOsdhRK6jABDr
sYjdTxZDA8y0nPEw5ULyia+ZKOCUz/JB5hwHjC1UVqFshLEL1Y8c+VEMCBeOnXjq
xqGqKh4VdwqgTCoYCp3/MWle1zNXyqLxe3k3BJPKI1ljLj40ncBBPMVOxx0wJG4B
7CpIUVzL+JxemTQk9c4U
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
HBQAZghJs/gE4zi2wJu6L6w7z1p22zCMsFZvI9t4NLdspOAKEFYDMQvInf+v+CCv
eCfmasqbrQJBAN6AuNA5rCSyNHyXFmMPAwU+nnH0BkkMgAcPEtt5JLheapiluYSw
P24V0fUhNHGe/8vAthinkEjT9MFlVa/ZHv0CQQCeeUkF+ydyEoVhxTz717Km0U3l
1RkOUKD+89pOqg38muM15F885LN+5Yz+F91YcBObGkJKjrlCAkcqz8DWHrTFAkBS
0qF4yO7+HeORuP/ZUb4zFpMOIeKxEFkbx42iap6zjlmpho7fCGgkBzVHRNvrq17W
Ll7aII2Bvnwt/RV/RpfVAkEAyZBS3D61oqF15ZZ3BdUPEuMkgiDkcZ843nP+8tMF
tfgT7foVUnEqBbKSe71tJG15ZVvUDQ+5yb7AroGT2131ZA==
-----END RSA PRIVATE KEY-----
EOF
}
`
const testAccServerCertificateConfigTwo = `
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
pLelaJcUpGZgr58Ha0bF3atYedtm1YEfH/kCAwEAAaN7MHkwCQYDVR0TBAIwADAs
BglghkgBhvhCAQ0EHxYdT3BlblNTTCBHZW5lcmF0ZWQgQ2VydGlmaWNhdGUwHQYD
VR0OBBYEFGO4AhAHP9vucMikI+xeh6/OGiI5MB8GA1UdIwQYMBaAFPwY6mMxV/Dv
1tJCusIXdZUkOdf2MA0GCSqGSIb3DQEBCwUAA4GBAHL14KnNjk6ZOsdhRK6jABDr
sYjdTxZDA8y0nPEw5ULyia+ZKOCUz/JB5hwHjC1UVqFshLEL1Y8c+VEMCBeOnXjq
xqGqKh4VdwqgTCoYCp3/MWle1zNXyqLxe3k3BJPKI1ljLj40ncBBPMVOxx0wJG4B
7CpIUVzL+JxemTQk9c4U
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
HBQAZghJs/gE4zi2wJu6L6w7z1p22zCMsFZvI9t4NLdspOAKEFYDMQvInf+v+CCv
eCfmasqbrQJBAN6AuNA5rCSyNHyXFmMPAwU+nnH0BkkMgAcPEtt5JLheapiluYSw
P24V0fUhNHGe/8vAthinkEjT9MFlVa/ZHv0CQQCeeUkF+ydyEoVhxTz717Km0U3l
1RkOUKD+89pOqg38muM15F885LN+5Yz+F91YcBObGkJKjrlCAkcqz8DWHrTFAkBS
0qF4yO7+HeORuP/ZUb4zFpMOIeKxEFkbx42iap6zjlmpho7fCGgkBzVHRNvrq17W
Ll7aII2Bvnwt/RV/RpfVAkEAyZBS3D61oqF15ZZ3BdUPEuMkgiDkcZ843nP+8tMF
tfgT7foVUnEqBbKSe71tJG15ZVvUDQ+5yb7AroGT2131ZA==
-----END RSA PRIVATE KEY-----
EOF
  name = "test"
  description = "test"
}
`
