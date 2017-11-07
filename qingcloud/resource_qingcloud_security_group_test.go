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

func TestAccQingcloudSecurityGroup_basic(t *testing.T) {
	var sg qc.DescribeSecurityGroupsOutput
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "qingcloud_security_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("qingcloud_security_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group.foo", "name", "first_eip"),
				),
			},
			resource.TestStep{
				Config: testAccSecurityGroupConfigTwo,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("qingcloud_security_group.foo", &sg),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group.foo", "name", "test"),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group.foo", "description", "test"),
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
		if err == nil && qc.IntValue(output.RetCode) == 0 {
			if len(output.SecurityGroupSet) != 0 {
				return fmt.Errorf("Found  SecurityGroup: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAccSecurityGroupConfig = `
resource "qingcloud_security_group" "foo" {
    name = "first_eip"
} `
const testAccSecurityGroupConfigTwo = `
resource "qingcloud_security_group" "foo" {
    name = "test"
	description = "test"
}`
