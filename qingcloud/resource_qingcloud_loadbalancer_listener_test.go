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
			resource.TestStep{
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
			resource.TestStep{
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
			resource.TestStep{
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
