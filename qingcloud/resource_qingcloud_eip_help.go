package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyEipAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip
	input := new(qc.ModifyEIPAttributesInput)
	input.EIP = qc.String(d.Id())
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
		if d.Get("name") == "" {
			input.EIPName = qc.String(" ")
		} else {
			input.EIPName = qc.String(d.Get("name").(string))
		}
		attributeUpdate = true
	}
	if attributeUpdate {
		var output *qc.ModifyEIPAttributesOutput
		var err error
		simpleRetry(func() error {
			output, err = clt.ModifyEIPAttributes(input)
			return isServerBusy(err)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func getEIPResourceMap(data *qc.EIP) map[string]interface{} {
	var a = make(map[string]interface{}, 3)
	a["resource_id"] = qc.StringValue(data.Resource.ResourceID)
	a["resource_name"] = qc.StringValue(data.Resource.ResourceName)
	a["resource_type"] = qc.StringValue(data.Resource.ResourceType)
	return a
}
func waitEipLease(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).eip
	input := new(qc.DescribeEIPsInput)
	input.EIPs = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	var output *qc.DescribeEIPsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeEIPs(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	//wait for lease info
	WaitForLease(output.EIPSet[0].CreateTime)
	return nil
}
