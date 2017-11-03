package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyEipAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).eip
	input := new(qc.ModifyEIPAttributesInput)
	input.EIP = qc.String(d.Id())
	if create {
		if description := d.Get("description").(string); description == "" {
			return nil
		}
		input.Description = qc.String(d.Get("description").(string))
	} else {
		if !d.HasChange("description") && !d.HasChange("name") {
			return nil
		}
		if d.HasChange("description") {
			if d.Get("description") == "" {
				input.Description = qc.String(" ")
			} else {
				input.Description = qc.String(d.Get("description").(string))
			}
		}
		if d.HasChange("name") {
			if d.Get("name") == "" {
				input.EIPName = qc.String(" ")
			} else {
				input.EIPName = qc.String(d.Get("name").(string))
			}
		}
	}
	output, err := clt.ModifyEIPAttributes(input)
	if err != nil {
		return fmt.Errorf("Error modify eip attributes: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error modify eip attributes: %s", *output.Message)
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
	describeoutput, err := clt.DescribeEIPs(describeinput)
	if err != nil {
		return fmt.Errorf("Error describe eip: %s", err)
	}
	if *describeoutput.RetCode != 0 {
		return fmt.Errorf("Error describe eip: %s", *describeoutput.Message)
	}
	//wait for lease info
	WaitForLease(describeoutput.EIPSet[0].CreateTime)
	return nil
}
