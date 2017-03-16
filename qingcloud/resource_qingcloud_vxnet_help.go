package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func modifyVxnetAttributes(d *schema.ResourceData, meta interface{}, create bool) error {
	clt := meta.(*QingCloudClient).vxnet
	input := new(qc.ModifyVxNetAttributesInput)
	input.VxNet = qc.String(d.Id())
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
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error modify vxnet attributes: %s", *output.Message)
	}
	return nil
}

func computedVxnetVPCID(d *schema.ResourceData, meta interface{}, vxnetID string) (string, error) {
	clt := meta.(*QingCloudClient).router
	input := new(qc.DescribeRoutersInput)
	input.Verbose = qc.Int(1)
	input.Status = []*string{qc.String("active")}
	err := input.Validate()
	if err != nil {
		return "", fmt.Errorf("Error describe router input validate: %s", err)
	}
	output, err := clt.DescribeRouters(input)
	if err != nil {
		return "", fmt.Errorf("Error describe router: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return "", fmt.Errorf("Error describe router: %s", err)
	}
	for _, router := range output.RouterSet {
		if router.RouterID != nil {
			DescribeRouterVxnetsInput := new(qc.DescribeRouterVxNetsInput)
			DescribeRouterVxnetsInput.Router = router.RouterID
			err := DescribeRouterVxnetsInput.Validate()
			if err != nil {
				return "", fmt.Errorf("Error describe router vxnet: %s", err)
			}
			o, err := clt.DescribeRouterVxNets(DescribeRouterVxnetsInput)
			if err != nil {
				return "", fmt.Errorf("Error describe router vxnet: %s", err)
			}
			for _, vxnet := range o.RouterVxNetSet {
				if qc.StringValue(vxnet.VxNetID) == vxnetID {
					return qc.StringValue(router.RouterID), nil
				}
			}
		}
	}
	return "", nil
}
