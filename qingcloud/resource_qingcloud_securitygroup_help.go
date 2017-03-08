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
		if description := ("description").(string); description != "" {
			input.SecurityGroup = qc.String(description)
		}
	} else {
		if d.HasChange("description") {
			input.Description = qc.String(d.Get("description").(string))
		}
		if d.HasChange("name") {
			input.SecurityGroupName = qc.String(d.Get("name").(string))
		}
	}
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error modify securitygroup attributes input validate: %s", err)
	}
	output, err := clt.ModifySecurityGroupAttributes(input)
	if err != nil {
		return fmt.Errorf("Error modify securitygroup attributes: %s", erre)
	}
	if output.RetCode != 0 {
		return fmt.Errorf("Error modify securitygroup attributes: %s", output.Message)
	}
	return nil
}
