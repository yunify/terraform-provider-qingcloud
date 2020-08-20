package qingcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	qc "github.com/yunify/qingcloud-sdk-go/service"
	"os"
)

func TestAccQingcloudSecurityGroup_basic(t *testing.T) {
	var sg qc.DescribeSecurityGroupsOutput
	testTag := "terraform-test-sg-basic" + os.Getenv("TRAVIS_BUILD_ID") + "-" + os.Getenv("TRAVIS_JOB_NUMBER")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "qingcloud_security_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccSecurityGroupConfig, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("qingcloud_security_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group.foo", resourceName, "first_sg"),
				),
			},
			{
				Config: fmt.Sprintf(testAccSecurityGroupConfigTwo, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("qingcloud_security_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group.foo", resourceName, "test"),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group.foo", resourceDescription, "test"),
				),
			},
		},
	})
}

func TestAccQingcloudSecurityGroup_tag(t *testing.T) {
	var sg qc.DescribeSecurityGroupsOutput
	sgTag1Name := "terraform-" + os.Getenv("CIRCLE_BUILD_NUM") + "-sg-tag1"
	sgTag2Name := "terraform-" + os.Getenv("CIRCLE_BUILD_NUM") + "-sg-tag2"
	testTagNameValue := func(names ...string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			tags := sg.SecurityGroupSet[0].Tags
			same_count := 0
			for _, tag := range tags {
				for _, name := range names {
					if qc.StringValue(tag.TagName) == name {
						same_count++
					}
					if same_count == len(sg.SecurityGroupSet[0].Tags) {
						return nil
					}
				}
			}
			return fmt.Errorf("tag name error %#v", names)
		}
	}
	testTagDetach := func() resource.TestCheckFunc {
		return func(state *terraform.State) error {
			if len(sg.SecurityGroupSet[0].Tags) != 0 {
				return fmt.Errorf("tag not detach ")
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "qingcloud_security_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccSecurityGroupConfigTagTemplate, sgTag1Name, sgTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(
						"qingcloud_security_group.foo", &sg),
					testTagNameValue(sgTag1Name, sgTag2Name),
				),
			},
			{
				Config: fmt.Sprintf(testAccSecurityGroupConfigTagTwoTemplate, sgTag1Name, sgTag2Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists(
						"qingcloud_security_group.foo", &sg),
					testTagDetach(),
				),
			},
		},
	})

}

func testAccCheckSecurityGroupExists(n string, sg *qc.DescribeSecurityGroupsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s ", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No SecurityGroup ID is set ")
		}
		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeSecurityGroupsInput)
		input.SecurityGroups = []*string{qc.String(rs.Primary.ID)}
		d, err := client.securitygroup.DescribeSecurityGroups(input)
		log.Printf("[WARN] SecurityGroup id %#v", rs.Primary.ID)
		if err != nil {
			return err
		}
		if d == nil || len(d.SecurityGroupSet) == 0 {
			return fmt.Errorf("SecurityGroup not found")
		}
		*sg = *d
		return nil
	}
}

func testAccCheckSecurityGroupDestroy(s *terraform.State) error {
	return testAccCheckSecurityGroupDestroyWithProvider(s, testAccProvider)
}

func testAccCheckSecurityGroupDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_security_group" {
			continue
		}
		// Try to find the resource
		input := new(qc.DescribeSecurityGroupsInput)
		input.SecurityGroups = []*string{qc.String(rs.Primary.ID)}
		output, err := client.securitygroup.DescribeSecurityGroups(input)
		if err == nil {
			if len(output.SecurityGroupSet) != 0 {
				return fmt.Errorf("Found  SecurityGroup: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAccSecurityGroupConfig = `
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`
const testAccSecurityGroupConfigTwo = `
resource "qingcloud_security_group" "foo" {
    name = "test"
	description = "test"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}
`

const testAccSecurityGroupConfigTagTemplate = `

resource "qingcloud_security_group" "foo" {
    name = "first_sg"
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
const testAccSecurityGroupConfigTagTwoTemplate = `

resource "qingcloud_security_group" "foo" {
    name = "first_sg"
}
resource "qingcloud_tag" "test"{
	name="%v"
}
resource "qingcloud_tag" "test2"{
	name="%v"
}
`
