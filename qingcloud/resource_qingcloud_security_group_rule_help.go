package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/lowstz/qingcloud-sdk-go/service"
)

func ModifySecurityGroupRuleAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.ModifySecurityGroupRuleAttributesInput)
	if create {
		return nil
	}
	if !d.HasChange("protocol") && !d.HasChange("priority") && !d.HasChange("action") && !d.HasChange("from_port") &&
		!d.HasChange("to_port") && !d.HasChange("cidr_block") {
		return nil
	}
	if d.HasChange("protocol") {
		input.Protocol = qc.String(d.Get("protocol").(string))
	}
	if d.HasChange("priority") {
		input.Priority = qc.Int(d.Get("priority").(int))
	}
	if d.HasChange("action") {
		input.RuleAction = qc.String(d.Get("action").(string))
	}
	if d.HasChange("from_port") {
		input.Val1 = qc.Int(d.Get("from_port").(int))
	}
	if d.HasChange("to_port") {
		input.Val2 = qc.Int(d.Get("to_port").(int))
	}
	_, err := clt.ModifySecurityGroupRuleAttributes(input)
	return err
}
