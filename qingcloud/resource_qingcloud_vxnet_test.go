/**
 * Copyright (c) 2016 Magicshui
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
*/

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

func TestAccQingcloudVxNet_basic(t *testing.T) {
	var vxnet qc.DescribeVxNetsOutput

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "qingcloud_vxnet.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVxNetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVxNetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVxNetExists(
						"qingcloud_vxnet.foo", &vxnet),
					resource.TestCheckResourceAttr(
						"qingcloud_vxnet.foo", "type", "1"),
				),
			},
			resource.TestStep{
				Config: testAccVxNetConfigTwo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVxNetExists(
						"qingcloud_vxnet.foo", &vxnet),
					resource.TestCheckResourceAttr(
						"qingcloud_vxnet.foo", "type", "1"),
					resource.TestCheckResourceAttr(
						"qingcloud_vxnet.foo", resourceDescription, "vxnet"),
					resource.TestCheckResourceAttr(
						"qingcloud_vxnet.foo", resourceName, "vxnet"),
				),
			},
			resource.TestStep{
				Config: testAccVxNetConfigThree,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVxNetExists(
						"qingcloud_vxnet.foo", &vxnet),
					resource.TestCheckResourceAttr(
						"qingcloud_vxnet.foo", "type", "0"),
					resource.TestCheckResourceAttr(
						"qingcloud_vxnet.foo", resourceDescription, "vxnet"),
					resource.TestCheckResourceAttr(
						"qingcloud_vxnet.foo", resourceName, "vxnet"),
				),
			},
		},
	})

}

func TestAccQingcloudVxNet_tag(t *testing.T) {
	var vxnet qc.DescribeVxNetsOutput
	vxnetTag1Name := os.Getenv("TRAVIS_BUILD_ID") + "-" + os.Getenv("TRAVIS_JOB_NUMBER") + "-vxnet-tag1"
	vxnetTag2Name := os.Getenv("TRAVIS_BUILD_ID") + "-" + os.Getenv("TRAVIS_JOB_NUMBER") + "-vxnet-tag2"

	testTagNameValue := func(names ...string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			tags := vxnet.VxNetSet[0].Tags
			same_count := 0
			for _, tag := range tags {
				for _, name := range names {
					if qc.StringValue(tag.TagName) == name {
						same_count++
					}
					if same_count == len(vxnet.VxNetSet[0].Tags) {
						return nil
					}
				}
			}
			return fmt.Errorf("tag name error %#v", names)
		}
	}
	testTagDetach := func() resource.TestCheckFunc {
		return func(state *terraform.State) error {
			if len(vxnet.VxNetSet[0].Tags) != 0 {
				return fmt.Errorf("tag not detach ")
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "qingcloud_vxnet.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVxNetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccVxNetConfigTagTemplate, vxnetTag1Name, vxnetTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVxNetExists(
						"qingcloud_vxnet.foo", &vxnet),
					testTagNameValue(vxnetTag1Name, vxnetTag2Name),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testAccVxNetConfigTagTwoTemplate, vxnetTag1Name, vxnetTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVxNetExists(
						"qingcloud_vxnet.foo", &vxnet),
					testTagDetach(),
				),
			},
		},
	})

}

func TestAccQingcloudVxNet_vpc(t *testing.T) {
	var vxnet qc.DescribeVxNetsOutput

	testVpcAttach := func() resource.TestCheckFunc {
		return func(state *terraform.State) error {
			if vxnet.VxNetSet[0].Router != nil {
				input := new(qc.DescribeRouterVxNetsInput)
				input.Router = vxnet.VxNetSet[0].VpcRouterID
				input.Verbose = qc.Int(1)
				client := testAccProvider.Meta().(*QingCloudClient)
				d, err := client.router.DescribeRouterVxNets(input)
				if err != nil {
					return err
				}
				if d == nil || len(d.RouterVxNetSet) == 0 {
					return fmt.Errorf("Router not found ")
				}
				haveVxnet := false
				for _, oneVxnet := range d.RouterVxNetSet {
					if qc.StringValue(oneVxnet.VxNetID) == qc.StringValue(vxnet.VxNetSet[0].VxNetID) {
						haveVxnet = true
					}
				}
				if !haveVxnet {
					return fmt.Errorf("Router not match ")
				}
				return nil
			} else {
				return fmt.Errorf("Can not find router ")
			}
		}
	}
	testVpcDetach := func() resource.TestCheckFunc {
		return func(state *terraform.State) error {
			if vxnet.VxNetSet[0].Router != nil && qc.StringValue(vxnet.VxNetSet[0].Router.RouterID) != "" {
				return fmt.Errorf("Router not detach ")
			}
			return nil
		}
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "qingcloud_vxnet.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVxNetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVxNetConfigVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVxNetExists(
						"qingcloud_vxnet.foo", &vxnet),
					testVpcAttach(),
				),
			},
			resource.TestStep{
				Config: testAccVxNetConfigVpcTwo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVxNetExists(
						"qingcloud_vxnet.foo", &vxnet),
					testVpcDetach(),
				),
			},
			resource.TestStep{
				Config: testAccVxNetConfigVpcThree,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVxNetExists(
						"qingcloud_vxnet.foo", &vxnet),
				),
			},
		},
	})

}

func testAccCheckVxNetExists(n string, vxnet *qc.DescribeVxNetsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VxNet ID is set")
		}

		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeVxNetsInput)
		input.VxNets = []*string{qc.String(rs.Primary.ID)}
		d, err := client.vxnet.DescribeVxNets(input)

		log.Printf("[WARN] vxnet id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil || len(d.VxNetSet) == 0 {
			return fmt.Errorf("VxNet not found")
		}

		*vxnet = *d
		return nil
	}
}

func testAccCheckVxNetDestroy(s *terraform.State) error {
	return testAccCheckVxNetDestroyWithProvider(s, testAccProvider)
}

func testAccCheckVxNetDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_vxnet" {
			continue
		}
		// Try to find the resource
		input := new(qc.DescribeVxNetsInput)
		input.VxNets = []*string{qc.String(rs.Primary.ID)}
		output, err := client.vxnet.DescribeVxNets(input)
		if err == nil {
			if len(output.VxNetSet) != 0 {
				return fmt.Errorf("Found  VxNet: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAccVxNetConfig = `
resource "qingcloud_vxnet" "foo" {
    type = 1
} `

const testAccVxNetConfigTwo = `
resource "qingcloud_vxnet" "foo" {
    name = "vxnet"
    description = "vxnet"
	type = 1
} `
const testAccVxNetConfigThree = `
resource "qingcloud_vxnet" "foo" {
    name = "vxnet"
    description = "vxnet"
	type = 0
} `

const testAccVxNetConfigTagTemplate = `

resource "qingcloud_vxnet" "foo" {
    type = 1
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

const testAccVxNetConfigTagTwoTemplate = `

resource "qingcloud_vxnet" "foo" {
    type = 1
}
resource "qingcloud_tag" "test"{
	name="%v"
}
resource "qingcloud_tag" "test2"{
	name="%v"
}
`

const testAccVxNetConfigVpc = `

resource "qingcloud_security_group" "foo" {
    name = "first_sg"
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
}

resource "qingcloud_vxnet" "foo" {
    type = 1
	vpc_id = "${qingcloud_vpc.foo.id}"
	ip_network = "192.168.0.0/24"
}
`

const testAccVxNetConfigVpcTwo = `

resource "qingcloud_security_group" "foo" {
    name = "first_sg"
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
}

resource "qingcloud_vxnet" "foo" {
    type = 1
}
`
const testAccVxNetConfigVpcThree = `

resource "qingcloud_security_group" "foo" {
    name = "first_sg"
}
resource "qingcloud_vpc" "foo" {
	security_group_id = "${qingcloud_security_group.foo.id}"
	vpc_network = "192.168.0.0/16"
}

resource "qingcloud_vxnet" "foo" {
    type = 1
	vpc_id = "${qingcloud_vpc.foo.id}"
	ip_network = "192.168.0.0/24"
}
`
