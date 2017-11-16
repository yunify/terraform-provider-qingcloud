package qingcloud

import (
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
		var output *qc.ModifyVxNetAttributesOutput
		var err error
		retryServerBusy(func() (*int, error) {
			output, err = clt.ModifyVxNetAttributes(input)
			return output.RetCode, err
		})
		if err := getQingCloudErr("modify vxnet attributes", output.RetCode, output.Message, err); err != nil {
			return err
		}
	}
	return nil
}

func vxnetJoinRouter(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.JoinRouterInput)
	input.VxNet = qc.String(d.Id())
	input.Router = qc.String(d.Get("vpc_id").(string))
	input.IPNetwork = qc.String(d.Get("ip_network").(string))
	var output *qc.JoinRouterOutput
	var err error
	retryServerBusy(func() (*int, error) {
		output, err = clt.JoinRouter(input)
		return output.RetCode, err
	})
	if err := getQingCloudErr("join router", output.RetCode, output.Message, err); err != nil {
		return err
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
	var output *qc.LeaveRouterOutput
	var err error
	retryServerBusy(func() (*int, error) {
		output, err = clt.LeaveRouter(input)
		return output.RetCode, err
	})
	if err := getQingCloudErr("leave router", output.RetCode, output.Message, err); err != nil {
		return err
	}
	if _, err := VxnetLeaveRouterTransitionStateRefresh(meta.(*QingCloudClient).vxnet, d.Id()); err != nil {
		return err
	}
	if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, d.Get("vpc_id").(string)); err != nil {
		return err
	}
	return nil
}
