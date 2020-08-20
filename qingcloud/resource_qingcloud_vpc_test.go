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

func TestAccQingcloudVpc_basic(t *testing.T) {
	var vpc qc.DescribeRoutersOutput
	testTag := "terraform-test-vpc-basic" + os.Getenv("CIRCLE_BUILD_NUM")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "qingcloud_vpc.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccVpcConfig, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(
						"qingcloud_vpc.foo", &vpc),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc.foo", "vpc_network", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc.foo", "type", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccVpcConfigTwo, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(
						"qingcloud_vpc.foo", &vpc),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc.foo", "type", "2"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc.foo", "vpc_network", "172.24.0.0/16"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc.foo", resourceDescription, "test"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc.foo", resourceName, "test"),
				),
			},
		},
	})

}

func TestAccQingcloudVpc_eip(t *testing.T) {
	var vpc qc.DescribeRoutersOutput
	testTag := "terraform-test-vpc-eip" + os.Getenv("CIRCLE_BUILD_NUM")

	testEIP := func() resource.TestCheckFunc {
		return func(state *terraform.State) error {
			if vpc.RouterSet[0].EIP != nil {
				input := new(qc.DescribeEIPsInput)
				input.EIPs = []*string{vpc.RouterSet[0].EIP.EIPID}
				client := testAccProvider.Meta().(*QingCloudClient)
				d, err := client.eip.DescribeEIPs(input)

				if err != nil {
					return err
				}
				if d == nil || len(d.EIPSet) == 0 {
					return fmt.Errorf("EIP not found ")
				}
				if qc.StringValue(d.EIPSet[0].EIPAddr) != qc.StringValue(vpc.RouterSet[0].EIP.EIPAddr) {
					return fmt.Errorf("EIP not match ")
				}
				return nil
			} else {
				return fmt.Errorf("Can not find eip ")
			}
		}
	}
	testDetachEIP := func() resource.TestCheckFunc {
		return func(state *terraform.State) error {
			if vpc.RouterSet[0].EIP != nil && vpc.RouterSet[0].EIP.EIPID != nil && qc.StringValue(vpc.RouterSet[0].EIP.EIPID) != "" {
				return fmt.Errorf("EIP not detach ")
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "qingcloud_vpc.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccVpcConfigEIP, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(
						"qingcloud_vpc.foo", &vpc),
					testEIP(),
				),
			},
			{
				Config: fmt.Sprintf(testAccVpcConfigEIPTwo, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(
						"qingcloud_vpc.foo", &vpc),
					testDetachEIP(),
				),
			},
		},
	})

}

func TestAccQingcloudVpc_tag(t *testing.T) {
	var vpc qc.DescribeRoutersOutput
	vpcTag1Name := "terraform-" + os.Getenv("CIRCLE_BUILD_NUM") + "-vpc-tag1"
	vpcTag2Name := "terraform-" + os.Getenv("CIRCLE_BUILD_NUM") + "-vpc-tag2"
	testTagNameValue := func(names ...string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			tags := vpc.RouterSet[0].Tags
			sameCount := 0
			for _, tag := range tags {
				for _, name := range names {
					if qc.StringValue(tag.TagName) == name {
						sameCount++
					}
					if sameCount == len(vpc.RouterSet[0].Tags) {
						return nil
					}
				}
			}
			return fmt.Errorf("tag name error %#v", names)
		}
	}
	testTagDetach := func() resource.TestCheckFunc {
		return func(state *terraform.State) error {
			if len(vpc.RouterSet[0].Tags) != 0 {
				return fmt.Errorf("tag not detach ")
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "qingcloud_vpc.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccVpcConfigTagTemplate, vpcTag1Name, vpcTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(
						"qingcloud_vpc.foo", &vpc),
					testTagNameValue(vpcTag1Name, vpcTag2Name),
				),
			},
			{
				Config: fmt.Sprintf(testAccVpcConfigTagTwoTemplate, vpcTag1Name, vpcTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(
						"qingcloud_vpc.foo", &vpc),
					testTagDetach(),
				),
			},
		},
	})

}

func testAccCheckVpcExists(n string, router *qc.DescribeRoutersOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Vpc ID is set")
		}

		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeRoutersInput)
		input.Routers = []*string{qc.String(rs.Primary.ID)}
		d, err := client.router.DescribeRouters(input)

		log.Printf("[WARN] router id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil {
			return fmt.Errorf("Router not found ")
		}

		*router = *d
		return nil
	}
}

func testAccCheckVpcDestroy(s *terraform.State) error {
	return testAccCheckVpcDestroyWithProvider(s, testAccProvider)
}

func testAccCheckVpcDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_vpc" {
			continue
		}

		// Try to find the resource
		input := new(qc.DescribeRoutersInput)
		input.Routers = []*string{qc.String(rs.Primary.ID)}
		output, err := client.router.DescribeRouters(input)
		if err == nil {
			if !isRouterDeleted(output.RouterSet) {
				return fmt.Errorf("Found  Router: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAccVpcConfig = `
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
`
const testAccVpcConfigTwo = `
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "172.24.0.0/16"
	name ="test"
	description = "test"
	type = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`

const testAccVpcConfigTagTemplate = `

resource "qingcloud_security_group" "foo" {
    name = "first_sg"
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
	tag_ids = ["${qingcloud_tag.test.id}",
				"${qingcloud_tag.test2.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
resource "qingcloud_tag" "test2"{
	name="%v"
}
`

const testAccVpcConfigTagTwoTemplate = `

resource "qingcloud_security_group" "foo" {
    name = "first_sg"
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
}
resource "qingcloud_tag" "test"{
	name="%v"
}
resource "qingcloud_tag" "test2"{
	name="%v"
}
`
const testAccVpcConfigEIP = `
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_eip" "foo" {
    bandwidth = 2
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	eip_id = "${qingcloud_eip.foo.id}"
	vpc_network = "192.168.0.0/16"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`
const testAccVpcConfigEIPTwo = `
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_eip" "foo" {
    bandwidth = 2
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
`
