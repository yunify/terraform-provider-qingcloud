package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func ModifySecurityGroupRuleAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.ModifySecurityGroupRuleAttributesInput)
	input.SecurityGroupRule = qc.String(d.Id())
	if d.Get(resourceName).(string) != "" {
		input.SecurityGroupRuleName = qc.String(d.Get(resourceName).(string))
	} else if d.HasChange(resourceName) {
		return fmt.Errorf("name can not be modified to nil")
	} else {
		input.SecurityGroupRuleName = nil
	}
	input.Direction = qc.Int(d.Get("direction").(int))
	input.SecurityGroup = qc.String(d.Get("security_group_id").(string))
	input.Protocol = qc.String(d.Get("protocol").(string))
	input.Priority = qc.Int(d.Get("priority").(int))
	input.RuleAction = qc.String(d.Get("action").(string))
	input.Val1 = getUpdateStringPointer(d, "from_port")
	input.Val2 = getUpdateStringPointer(d, "to_port")
	input.Val3 = getUpdateStringPointer(d, "cidr_block")
	var output *qc.ModifySecurityGroupRuleAttributesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.ModifySecurityGroupRuleAttributes(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	return nil
}
