package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifySecurityGroupAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.ModifySecurityGroupAttributesInput)
	input.SecurityGroup = qc.String(d.Id())
	if create {
		if description := d.Get("description").(string); description == "" {
			return nil
		}
		input.Description = qc.String(d.Get("description").(string))
	} else {
		if !d.HasChange("description") && !d.HasChange("name") {
			return nil
		}
		if d.HasChange("description") {
			input.Description = qc.String(d.Get("description").(string))
		}
		if d.HasChange("name") {
			input.SecurityGroupName = qc.String(d.Get("name").(string))
		}
	}
	_, err := clt.ModifySecurityGroupAttributes(input)
	if err != nil {
		return fmt.Errorf("Error modify security group attributes: %s", err)
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
