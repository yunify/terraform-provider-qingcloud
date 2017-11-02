package qingcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func TestAccQingcloudEIP_basic(t *testing.T) {
	var eip qc.DescribeEIPsOutput

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "qingcloud_eip.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(
						"qingcloud_eip.foo", &eip),
					resource.TestCheckResourceAttr(
						"qingcloud_eip.foo", "bandwidth", "2"),
					resource.TestCheckResourceAttr(
						"qingcloud_eip.foo", "billing_mode", "traffic"),
					resource.TestCheckResourceAttr(
						"qingcloud_eip.foo", "description", "first"),
					resource.TestCheckResourceAttr(
						"qingcloud_eip.foo", "name", "first_eip"),
					resource.TestCheckResourceAttr(
						"qingcloud_eip.foo", "need_icp", "0"),
				),
			},
			resource.TestStep{
				Config: testAccEIPConfigTwo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(
						"qingcloud_eip.foo", &eip),
					resource.TestCheckResourceAttr(
						"qingcloud_eip.foo", "bandwidth", "4"),
					resource.TestCheckResourceAttr(
						"qingcloud_eip.foo", "billing_mode", "bandwidth"),
					resource.TestCheckResourceAttr(
						"qingcloud_eip.foo", "description", "eip"),
					resource.TestCheckResourceAttr(
						"qingcloud_eip.foo", "name", "eip"),
					resource.TestCheckResourceAttr(
						"qingcloud_eip.foo", "need_icp", "0"),
				),
			},
		},
	})

}

func testAccCheckEIPExists(n string, eip *qc.DescribeEIPsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EIP ID is set")
		}

		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeEIPsInput)
		input.EIPs = []*string{qc.String(rs.Primary.ID)}
		d, err := client.eip.DescribeEIPs(input)

		log.Printf("[WARN] eip id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil || qc.StringValue(d.EIPSet[0].EIPAddr) == "" {
			return fmt.Errorf("EIP not found")
		}

		*eip = *d
		return nil
	}
}
func testAccCheckEIPDestroy(s *terraform.State) error {
	return nil
}

const testAccEIPConfig = `
resource "qingcloud_eip" "foo" {
    name = "first_eip"
    description = "first"
    billing_mode = "traffic"
    bandwidth = 2
    need_icp = 0
} `
const testAccEIPConfigTwo = `
resource "qingcloud_eip" "foo" {
    name = "eip"
    description = "eip"
    billing_mode = "bandwidth"
    bandwidth = 4
    need_icp = 0
} `
