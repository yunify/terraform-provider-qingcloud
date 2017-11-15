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
						"qingcloud_vxnet.foo", "description", "vxnet"),
					resource.TestCheckResourceAttr(
						"qingcloud_vxnet.foo", "name", "vxnet"),
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
						"qingcloud_vxnet.foo", "description", "vxnet"),
					resource.TestCheckResourceAttr(
						"qingcloud_vxnet.foo", "name", "vxnet"),
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

func testAccCheckVxNetExists(n string, eip *qc.DescribeVxNetsOutput) resource.TestCheckFunc {
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

		log.Printf("[WARN] eip id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		if d == nil || len(d.VxNetSet) == 0 {
			return fmt.Errorf("VxNet not found")
		}

		*eip = *d
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
		if err == nil && qc.IntValue(output.RetCode) == 0 {
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
