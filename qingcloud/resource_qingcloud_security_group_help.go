package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifySecurityGroupAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.ModifySecurityGroupAttributesInput)
	input.SecurityGroup = qc.String(d.Id())
	attributeUpdate := false
	if d.HasChange("description") {
		if d.Get("description") == "" {
			input.Description = qc.String(" ")
		} else {
			input.Description = qc.String(d.Get("description").(string))
		}
		attributeUpdate = true
	}
	if d.HasChange("name") && !d.IsNewResource() {
		input.SecurityGroupName = qc.String(d.Get("name").(string))
		attributeUpdate = true
	}
	if attributeUpdate {
		_, err := clt.ModifySecurityGroupAttributes(input)
		if err != nil {
			return fmt.Errorf("Error modify security group attributes: %s ", err)
		}
	}
	return nil
}

func applySecurityGroupRule(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	sgID := d.Get("security_group_id").(string)
	input := new(qc.ApplySecurityGroupInput)
	input.SecurityGroup = qc.String(sgID)
	_, err := clt.ApplySecurityGroup(input)
	if err != nil {
		return err
	}
	if _, err := SecurityGroupApplyTransitionStateRefresh(clt, sgID); err != nil {
		return err
	}
	return nil
}
