package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyRouterAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.ModifyRouterAttributesInput)
	input.Router = qc.String(d.Id())
	attributeUpdate := false
	attributeUpdate2 := false
	input.RouterName, attributeUpdate = getNamePointer(d)
	input.Description, attributeUpdate2 = getDescriptionPointer(d)
	if d.HasChange("eip_id") {
		if d.Get("eip_id") != "" {
			input.EIP = qc.String(d.Get("eip_id").(string))
		} else {
			input.EIP = qc.String(" ")
		}
		attributeUpdate = true
	}
	if d.HasChange("security_group_id") && !d.IsNewResource() {
		if d.Get("security_group_id") != "" {
			input.SecurityGroup = qc.String(d.Get("security_group_id").(string))
		} else {
			input.SecurityGroup = qc.String(" ")
		}
		attributeUpdate = true
	}

	if attributeUpdate || attributeUpdate2 {
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

func applyRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.UpdateRoutersInput)
	input.Routers = []*string{qc.String(d.Id())}
	var output *qc.UpdateRoutersOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.UpdateRouters(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if _, err = RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return fmt.Errorf("Error waiting for router (%s) to start: %s", d.Id(), err.Error())
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
