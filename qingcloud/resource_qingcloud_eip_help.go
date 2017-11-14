package qingcloud

import (
	"fmt"

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
			if err == nil {
				if output.RetCode != nil && *output.RetCode == 5100 {
					return fmt.Errorf("allocate EIP Server Busy")
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error modify eip attributes: %s", err)
		}
		if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
			return fmt.Errorf("Error modify eip attributes: %s", *output.Message)
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
	describeinput := new(qc.DescribeEIPsInput)
	describeinput.EIPs = []*string{qc.String(d.Id())}
	describeinput.Verbose = qc.Int(1)
	var output *qc.DescribeEIPsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeEIPs(describeinput)
		if err == nil {
			if output.RetCode != nil && *output.RetCode == 5100 {
				return fmt.Errorf("allocate EIP Server Busy")
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error describe eip: %s", err)
	}
	if *output.RetCode != 0 {
		return fmt.Errorf("Error describe eip: %s", *output.Message)
	}
	//wait for lease info
	WaitForLease(output.EIPSet[0].CreateTime)
	return nil
}
