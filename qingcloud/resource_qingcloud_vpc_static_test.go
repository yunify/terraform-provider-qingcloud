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

func TestAccQingcloudVpcStatic_basic(t *testing.T) {
	var vpcStatic qc.DescribeRouterStaticsOutput
	testTag := "terraform-test-vpc-static-basic" + os.Getenv("CIRCLE_BUILD_NUM")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "qingcloud_vpc_static.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpcStaticDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccVpcStaticConfig, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcStaticExists(
						"qingcloud_vpc_static.foo", &vpcStatic),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc_static.foo", "type", "1"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc_static.foo", "val1", "80"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc_static.foo", "val2", "192.168.0.3"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc_static.foo", "val3", "81"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc_static.foo", "val4", "tcp"),
				),
			},
			{
				Config: fmt.Sprintf(testAccVpcStaticConfigTwo, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcStaticExists(
						"qingcloud_vpc_static.foo", &vpcStatic),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc_static.foo", resourceName, "test"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc_static.foo", "type", "1"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc_static.foo", "val1", "81"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc_static.foo", "val2", "192.168.0.4"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc_static.foo", "val3", "82"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc_static.foo", "val4", "udp"),
				),
			},
		},
	})

}
func testAccCheckVpcStaticExists(n string, routerStatic *qc.DescribeRouterStaticsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VpcStatic ID is set ")
		}

		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeRouterStaticsInput)
		input.RouterStatics = []*string{qc.String(rs.Primary.ID)}
		input.Router = qc.String(rs.Primary.Attributes["vpc_id"])
		d, err := client.router.DescribeRouterStatics(input)

		log.Printf("[WARN] vpc static  id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil || len(d.RouterStaticSet) == 0 {
			return fmt.Errorf("VpcStatic not found")
		}

		*routerStatic = *d
		return nil
	}
}

func testAccCheckVpcStaticDestroy(s *terraform.State) error {
	return testAccCheckVpcStaticDestroyWithProvider(s, testAccProvider)
}

func testAccCheckVpcStaticDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_vpc_static" {
			continue
		}
		// Try to find the resource
		input := new(qc.DescribeRouterStaticsInput)
		input.RouterStatics = []*string{qc.String(rs.Primary.ID)}
		output, err := client.router.DescribeRouterStatics(input)
		if err == nil {
			if len(output.RouterStaticSet) != 0 {
				return fmt.Errorf("Found vpc static : %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAccVpcStaticConfig = `
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
resource "qingcloud_vpc_static" "foo"{
        vpc_id = "${qingcloud_vpc.foo.id}"
        type = 1
        val1 = "80"
        val2 = "192.168.0.3"
        val3 = "81"
}

`

const testAccVpcStaticConfigTwo = `
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
resource "qingcloud_vpc_static" "foo"{
        vpc_id = "${qingcloud_vpc.foo.id}"
		name = "test"
        type = 1
        val1 = "81"
        val2 = "192.168.0.4"
        val3 = "82"
		val4 = "udp"
}
`
