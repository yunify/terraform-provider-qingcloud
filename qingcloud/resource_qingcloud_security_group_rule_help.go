package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

// Warning
// The null character string and null pointer is difference when Go SDK processing parameter.
// For example,if val3 is "",then request has "val3=";if it is nil,and the request doesn't have "val3=".
func ModifySecurityGroupRuleAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.ModifySecurityGroupRuleAttributesInput)
	input.SecurityGroupRule = qc.String(d.Id())
	if d.Get("name").(string) != "" {
		input.SecurityGroupRuleName = qc.String(d.Get("name").(string))
	} else if d.HasChange("name") {
		return fmt.Errorf("name can not be modified to nil")
	} else {
		input.SecurityGroupRuleName = nil
	}
	input.Direction = qc.Int(d.Get("direction").(int))
	input.SecurityGroup = qc.String(d.Get("security_group_id").(string))
	input.Protocol = qc.String(d.Get("protocol").(string))
	input.Priority = qc.Int(d.Get("priority").(int))
	input.RuleAction = qc.String(d.Get("action").(string))
	if d.Get("from_port").(string) != "" {
		input.Val1 = qc.String(d.Get("from_port").(string))
	} else if d.HasChange("from_port") {
		input.Val1 = qc.String(" ")
	} else {
		input.Val1 = nil
	}
	if d.Get("to_port").(string) != "" {
		input.Val2 = qc.String(d.Get("to_port").(string))
	} else if d.HasChange("to_port") {
		input.Val2 = qc.String(" ")
	} else {
		input.Val2 = nil
	}
	if d.Get("cidr_block").(string) != "" {
		input.Val3 = qc.String(d.Get("cidr_block").(string))
	} else if d.HasChange("cidr_block") {
		input.Val3 = qc.String(" ")
	} else {
		input.Val3 = nil
	}
	var output *qc.ModifySecurityGroupRuleAttributesOutput
	var err error
	retryServerBusy(func() (s *int, err error) {
		output, err = clt.ModifySecurityGroupRuleAttributes(input)
		return output.RetCode, err
	})
	if err != nil {
		return err
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error update security group rule: %s", err)
	}
	return nil
}
