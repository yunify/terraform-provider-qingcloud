package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifySecurityGroupAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.ModifySecurityGroupAttributesInput)
	if create {
		if description := d.Get("description").(string); description == "" {
			return nil
		}
		input.SecurityGroup = qc.String(d.Get("description").(string))
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
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error modify security group attributes input validate: %s", err)
	}
	output, err := clt.ModifySecurityGroupAttributes(input)
	if err != nil {
		return fmt.Errorf("Error modify security group attributes: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error modify security group attributes: %s", *output.Message)
	}
	return nil
}
