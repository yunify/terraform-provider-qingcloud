package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

// func applyRouterUpdates(meta interface{}, routerID string) error {
// 	clt := meta.(*QingCloudClient).router
// 	params := router.UpdateRoutersRequest{}
// 	params.RoutersN.Add(routerID)
// 	if _, err := clt.UpdateRouters(params); err != nil {
// 		return err
// 	}
// 	if _, err := RouterTransitionStateRefresh(clt, routerID); err != nil {
// 		return err
// 	}
// 	return nil
// }

func modifyRouterAttributes(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.ModifyRouterAttributesInput)
	input.Router = qc.String(d.Id())
	attributeUpdate := false
	if d.HasChange("name") && !d.IsNewResource() {
		if d.Get("name") != "" {
			input.RouterName = qc.String(d.Get("name").(string))
		} else {
			input.RouterName = qc.String(" ")
		}
		attributeUpdate = true
	}
	if d.HasChange("description") {
		if d.Get("description") != "" {
			input.Description = qc.String(d.Get("description").(string))
		} else {
			input.Description = qc.String(" ")
		}
		attributeUpdate = true
	}
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

	if attributeUpdate {
		output, err := clt.ModifyRouterAttributes(input)
		if err != nil {
			return fmt.Errorf("Error modify router attributes: %s", err)
		}
		if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
			return fmt.Errorf("Error modify router attrubites: %s", *output.Message)
		}
		return nil
	}
	return nil
}

func applyRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.UpdateRoutersInput)
	input.Routers = []*string{qc.String(d.Id())}
	output, err := clt.UpdateRouters(input)
	if err != nil {
		return fmt.Errorf("Error update router: %s", err.Error())
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error update router: %s", *output.Message)
	}
	_, err = RouterTransitionStateRefresh(clt, d.Id())
	if err != nil {
		return fmt.Errorf("Error waiting for router (%s) to start: %s", d.Id(), err.Error())
	}

	return nil
}

func waitRouterLease(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	describeinput := new(qc.DescribeRoutersInput)
	describeinput.Routers = []*string{qc.String(d.Id())}
	describeinput.Verbose = qc.Int(1)
	describeoutput, err := clt.DescribeRouters(describeinput)
	if err != nil {
		return fmt.Errorf("Error describe router: %s", err)
	}
	if *describeoutput.RetCode != 0 {
		return fmt.Errorf("Error describe router: %s", *describeoutput.Message)
	}
	//wait for lease info
	WaitForLease(describeoutput.RouterSet[0].CreateTime)
	return nil
}

// func modifyRouterVxnets(d *schema.ResourceData, meta interface{}, create bool) error {
// 	clt := meta.(*QingCloudClient).router
// 	if create {
// 		map
// 	} else {
// 		if
// 	}
// }

// func getEIPInfoMap(data *qc.EIP) map[string]interface{} {
// 	var a = make(map[string]interface{}, 3)
// 	a["eip_id"] = qc.EIP.EIPID
// 	a["eip_name"] = qc.EIP.EIPName
// 	a["eip_addr"] = qc.EIP.EIPAddr
// 	return a
// }

// func getVxnetsMap(data []*qc.VxNet) map[string]interface{} {
// 	length := len(data)
// 	if data > 0 {
// 		var a = make(map[string]interface{}, length)
// 		for _, vxnet := range data {
// 			a[vxnet.VxNetID] = vxnet.NICID
// 		}
// 		return a
// 	}
// 	return nil
// }
