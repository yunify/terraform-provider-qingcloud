package qingcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/vxnet"
)

func modifyVxnetAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).vxnet
	params := vxnet.ModifyVxnetAttributesRequest{}
	params.Vxnet.Set(d.Id())
	if create {
		if description := d.Get("description").(string); description != "" {
			params.Description.Set(description)
		}
	} else {
		if d.HasChange("description") {
			params.Description.Set(d.Get("description").(string))
		}
		if d.HasChange("name") {
			params.VxnetName.Set(d.Get("name").(string))
		}
	}
	_, err := clt.ModifyVxnetAttributes(params)
	if err != nil {
		return fmt.Errorf("Error modify vxnet description: %s", err)
	}
	return nil
}
