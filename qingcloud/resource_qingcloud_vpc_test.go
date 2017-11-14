package qingcloud

import (
	"fmt"
	"log"
	"testing"
	"os"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	qc "github.com/yunify/qingcloud-sdk-go/service"

)

func TestAccQingcloudVpc_basic(t *testing.T) {
	var vpc qc.DescribeRoutersOutput

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "qingcloud_vpc.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(
						"qingcloud_vpc.foo", &vpc),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc.foo", "vpc_network", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc.foo", "type", "1"),
				),
			},
			resource.TestStep{
				Config: testAccVpcConfigTwo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(
						"qingcloud_vpc.foo", &vpc),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc.foo", "type", "2"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc.foo", "vpc_network", "172.24.0.0/16"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc.foo", "description", "test"),
					resource.TestCheckResourceAttr(
						"qingcloud_vpc.foo", "name", "test"),
				),
			},
		},
	})

}

func TestAccQingcloudVpc_tag(t *testing.T) {
	var vpc qc.DescribeRoutersOutput
	vpcTag1Name := os.Getenv("TRAVIS_BUILD_ID") + "-" + os.Getenv("TRAVIS_JOB_NUMBER") + "-vpc-tag1"
	vpcTag2Name := os.Getenv("TRAVIS_BUILD_ID") + "-" + os.Getenv("TRAVIS_JOB_NUMBER") + "-vpc-tag2"
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
			resource.TestStep{
				Config: fmt.Sprintf(testAccVpcConfigTagTemplate, vpcTag1Name, vpcTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(
						"qingcloud_vpc.foo", &vpc),
					testTagNameValue(vpcTag1Name, vpcTag2Name),
				),
			},
			resource.TestStep{
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
		if err == nil && qc.IntValue(output.RetCode) == 0 {
			if len(output.RouterSet) != 0 && qc.StringValue(output.RouterSet[0].Status) != "deleted" {
				return fmt.Errorf("Found  Router: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAccVpcConfig = `
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
}
`
const testAccVpcConfigTwo = `
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "172.24.0.0/16"
	name ="test"
	description = "test"
	type = 2
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
