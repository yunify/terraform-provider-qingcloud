package qingcloud

import (
	"fmt"
	"log"

	qc "github.com/yunify/qingcloud-sdk-go/service"
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

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
		if err == nil && qc.IntValue(output.RetCode) == 0 {
			if len(output.SecurityGroupRuleSet) != 0 {
				return fmt.Errorf("Found  SecurityGroupRule: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

