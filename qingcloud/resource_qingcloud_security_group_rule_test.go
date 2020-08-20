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

func TestAccQingcloudSecurityGroupRule_basic(t *testing.T) {
	var sgr qc.DescribeSecurityGroupRulesOutput
	testTag := "terraform-test-sgr-basic" + os.Getenv("CIRCLE_BUILD_NUM")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "qingcloud_security_group_rule.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccSecurityGroupRuleConfig, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("qingcloud_security_group_rule.foo", &sgr),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group_rule.foo", "protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group_rule.foo", "priority", "0"),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group_rule.foo", "action", "accept"),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group_rule.foo", "direction", "0"),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group_rule.foo", "from_port", "0"),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group_rule.foo", "to_port", "0"),
				),
			},
			{
				Config: fmt.Sprintf(testAccSecurityGroupRuleConfigTwo, testTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupRuleExists("qingcloud_security_group_rule.foo", &sgr),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group_rule.foo", resourceName, "first_sgr"),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group_rule.foo", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group_rule.foo", "priority", "1"),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group_rule.foo", "action", "drop"),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group_rule.foo", "direction", "1"),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group_rule.foo", "from_port", "10"),
					resource.TestCheckResourceAttr(
						"qingcloud_security_group_rule.foo", "to_port", "20"),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupRuleExists(n string, sg *qc.DescribeSecurityGroupRulesOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s ", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No SecurityGroupRule ID is set ")
		}
		client := testAccProvider.Meta().(*QingCloudClient)
		input := new(qc.DescribeSecurityGroupRulesInput)
		input.SecurityGroupRules = []*string{qc.String(rs.Primary.ID)}
		d, err := client.securitygroup.DescribeSecurityGroupRules(input)
		log.Printf("[WARN] SecurityGroupRule id %#v", rs.Primary.ID)
		if err != nil {
			return err
		}
		if d == nil || len(d.SecurityGroupRuleSet) == 0 {
			return fmt.Errorf("SecurityGroupRule not found")
		}
		*sg = *d
		return nil
	}
}

func testAccCheckSecurityGroupRuleDestroy(s *terraform.State) error {
	return testAccCheckSecurityGroupRuleDestroyWithProvider(s, testAccProvider)
}

func testAccCheckSecurityGroupRuleDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*QingCloudClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qingcloud_security_group_rule" {
			continue
		}
		// Try to find the resource
		input := new(qc.DescribeSecurityGroupRulesInput)
		input.SecurityGroupRules = []*string{qc.String(rs.Primary.ID)}
		output, err := client.securitygroup.DescribeSecurityGroupRules(input)
		if err == nil {
			if len(output.SecurityGroupRuleSet) != 0 {
				return fmt.Errorf("Found  SecurityGroupRule: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

const testAccSecurityGroupRuleConfig = `
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}

resource "qingcloud_security_group_rule" "foo"{
    security_group_id= "${qingcloud_security_group.foo.id}"
    protocol = "tcp"
    priority = 0
    action = "accept"
    direction = 0
    from_port = 0
    to_port = 0
}
`
const testAccSecurityGroupRuleConfigTwo = `
resource "qingcloud_security_group" "foo" {
    name = "first_sg"
	tag_ids = ["${qingcloud_tag.test.id}"]
}
resource "qingcloud_tag" "test"{
	name="%v"
}

resource "qingcloud_security_group_rule" "foo"{
    security_group_id= "${qingcloud_security_group.foo.id}"
    name = "first_sgr"
    protocol = "udp"
    priority = 1
    action = "drop"
    direction = 1
    from_port = 10
    to_port = 20
}
`
