package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyVxnetAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).vxnet
	input := new(qc.ModifyVxNetAttributesInput)
	input.VxNet = []string(qc.String(d.Id()))

	if create {
		if description := d.Get("description").(string); description != "" {
			params.Description.Set(description)
			input.Description = qc.String(d.Get("description").(string))
		}
	} else {
		if d.HasChange("description") {
			input.Description = qc.String(d.Get("description").(string))
		}
		if d.HasChange("name") {
			input.VxNetName = qc.String(d.Get("name").(string))
		}
	}
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error modify vxnet attributes input validate: %s", err)
	}
	output, err := clt.ModifyVxNetAttributes(input)
	if err != nil {
		return fmt.Errorf("Error modify vxnet attributes: %s", err)
	}
	if output.RetCode != 0 {
		return fmt.Errorf("Error modify vxnet attributes: %s", output.Message)
	}
	return nil
}
