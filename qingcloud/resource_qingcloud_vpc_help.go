package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyRouterAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.ModifyRouterAttributesInput)
	input.Router = qc.String(d.Id())
	attributeUpdate := false
	descriptionUpdate := false
	input.RouterName, attributeUpdate = getNamePointer(d)
	input.Description, descriptionUpdate = getDescriptionPointer(d)
	if d.HasChange("eip_id") {
		input.EIP = getUpdateStringPointer(d, "eip_id")
		attributeUpdate = true
	}
	if d.HasChange("security_group_id") && !d.IsNewResource() {
		input.SecurityGroup = getUpdateStringPointer(d, "security_group_id")
		attributeUpdate = true
	}

	if attributeUpdate || descriptionUpdate {
		var output *qc.ModifyRouterAttributesOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.ModifyRouterAttributes(input)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func applyRouterUpdate(routerId *string, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.UpdateRoutersInput)
	input.Routers = []*string{routerId}
	var output *qc.UpdateRoutersOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.UpdateRouters(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if _, err = RouterTransitionStateRefresh(clt, *routerId); err != nil {
		return err
	}
	return nil
}

func waitRouterLease(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.DescribeRoutersInput)
	input.Routers = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	var output *qc.DescribeRoutersOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeRouters(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	//wait for lease info
	WaitForLease(output.RouterSet[0].CreateTime)
	return nil
}
