package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyVxnetAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet
	input := new(qc.ModifyVxNetAttributesInput)
	input.VxNet = qc.String(d.Id())
	attributeUpdate := false
	if d.HasChange("description") {
		if d.Get("description").(string) != "" {
			input.Description = qc.String(d.Get("description").(string))
		} else {
			input.Description = qc.String(" ")
		}
		attributeUpdate = true
	}
	if d.HasChange("name") && !d.IsNewResource() {
		if d.Get("name").(string) != "" {
			input.VxNetName = qc.String(d.Get("description").(string))
		} else {
			input.VxNetName = qc.String(" ")
		}
		attributeUpdate = true
	}
	if attributeUpdate {
		_, err := clt.ModifyVxNetAttributes(input)
		return err
	}
	return nil
}

func vxnetJoinRouter(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.JoinRouterInput)
	input.VxNet = qc.String(d.Id())
	input.Router = qc.String(d.Get("vpc_id").(string))
	input.IPNetwork = qc.String(d.Get("ip_network").(string))
	output, err := clt.JoinRouter(input)
	if err != nil {
		return fmt.Errorf("Error vxnet join router: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error vxnet join router: %s", *output.Message)
	}
	if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, d.Get("vpc_id").(string)); err != nil {
		return err
	}
	return nil
}

func vxnetLeaverRouter(d *schema.ResourceData, meta interface{}) error {
	oldVPC, _ := d.GetChange("vpc_id")
	clt := meta.(*QingCloudClient).router
	input := new(qc.LeaveRouterInput)
	input.VxNets = []*string{qc.String(d.Id())}
	input.Router = qc.String(oldVPC.(string))
	output, err := clt.LeaveRouter(input)
	if err != nil {
		return fmt.Errorf("Error vxnet leave router: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error vxnet leave router: %s", *output.Message)
	}
	if _, err := VxnetLeaveRouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, d.Get("vpc_id").(string)); err != nil {
		return err
	}
	return nil
}
