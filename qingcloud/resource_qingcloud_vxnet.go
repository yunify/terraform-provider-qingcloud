package qingcloud

import (
	"errors"
	"fmt"
	// "log"

	// "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingcloudVxnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudVxnetCreate,
		Read:   resourceQingcloudVxnetRead,
		Update: resourceQingcloudVxnetUpdate,
		Delete: resourceQingcloudVxnetDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Description: "私有网络类型，1 - 受管私有网络，0 - 自管私有网络。	",
				ValidateFunc: withinArrayInt(0, 1),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// 当第一次创建一个私有网络以后，会首先加入到自己定制的router中，不是 vpc
			"router_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_network": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateVxnetsIPNetworkCIDR,
			},
		},
	}
}

func resourceQingcloudVxnetCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet

	routerID := d.Get("router_id").(string)
	ipNetwork := d.Get("ip_network").(string)
	if (routerID != "" && ipNetwork == "") || (routerID == "" && ipNetwork != "") {
		return errors.New("router and ip_network must both be empty or no empty at the same time")
	}

	input := new(qc.CreateVxNetsInput)
	input.Count = qc.Int(1)
	input.VxNetName = qc.String(d.Get("name").(string))
	input.VxNetType = qc.Int(d.Get("type").(int))
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Errorf create vxnet input validate: %s", err)
	}
	output, err := clt.CreateVxNets(input)
	if err != nil {
		return fmt.Errorf("Error create vxnet: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error create vxnet: %s", *output.Message)
	}
	d.SetId(qc.StringValue(output.VxNets[0]))
	if err := modifyVxnetAttributes(d, meta, true); err != nil {
		return err
	}
	if routerID != "" {
		if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, routerID); err != nil {
			return err
		}
		// join the router
		routerClt := meta.(*QingCloudClient).router
		joinRouterInput := new(qc.JoinRouterInput)
		joinRouterInput.VxNet = output.VxNets[0]
		joinRouterInput.Router = qc.String(routerID)
		joinRouterInput.IPNetwork = qc.String(ipNetwork)

		joinRouterOutput, err := routerClt.JoinRouter(joinRouterInput)
		if err != nil {
			return fmt.Errorf("Error create vxnet join router: %s", err)
		}
		if joinRouterOutput.RetCode != nil && qc.IntValue(joinRouterOutput.RetCode) != 0 {
			return fmt.Errorf("Error create vxnet join router: %s", *joinRouterOutput.Message)
		}
		if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, routerID); err != nil {
			return err
		}
	}
	return resourceQingcloudVxnetRead(d, meta)
}

func resourceQingcloudVxnetRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet
	input := new(qc.DescribeVxNetsInput)
	input.VxNets = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error describe vxnet input validate: %s", err)
	}
	output, err := clt.DescribeVxNets(input)
	if err != nil {
		return fmt.Errorf("Error describe vxnet: %s", err)
	}
	if output.RetCode == nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error describe vxnet: %s", *output.Message)
	}
	vxnet := output.VxNetSet[0]
	d.Set("name", qc.StringValue(vxnet.VxNetName))
	d.Set("type", qc.IntValue(vxnet.VxNetType))
	d.Set("description", qc.StringValue(vxnet.Description))
	return nil
}

func resourceQingcloudVxnetDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet

	describeVxnetInstanceInput := new(qc.DescribeVxNetInstancesInput)
	describeVxnetInstanceInput.VxNet = qc.String(d.Id())
	err := describeVxnetInstanceInput.Validate()
	if err != nil {
		return fmt.Errorf("Error describe vxnet instances input validate: %s", err)
	}
	describeVxnetInstanceOutput, err := clt.DescribeVxNetInstances(describeVxnetInstanceInput)
	if err != nil {
		return fmt.Errorf("Error describe vxnet instances: %s", err)
	}
	if describeVxnetInstanceOutput.RetCode != nil && qc.IntValue(describeVxnetInstanceOutput.RetCode) != 0 {
		return fmt.Errorf("Error describe vxnet instances: %s", *describeVxnetInstanceOutput.Message)
	}
	if qc.IntValue(describeVxnetInstanceOutput.TotalCount) > 0 {
		return fmt.Errorf("Error vxnet is using, can't delete")
	}

	routerID := d.Get("router_id").(string)
	// vxnet leave router
	if routerID != "" {
		if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, routerID); err != nil {
			return err
		}
		routerCtl := meta.(*QingCloudClient).router
		leaveRouterInput := new(qc.LeaveRouterInput)
		leaveRouterInput.Router = qc.String(routerID)
		leaveRouterInput.VxNets = []*string{qc.String(d.Id())}
		err := leaveRouterInput.Validate()
		if err != nil {
			return fmt.Errorf("Error leave router input validate: %s", err)
		}
		leaveRouterOutput, err := routerCtl.LeaveRouter(leaveRouterInput)
		if err != nil {
			return fmt.Errorf("Error leave router: %s", err)
		}
		if leaveRouterOutput.RetCode != nil && qc.IntValue(leaveRouterOutput.RetCode) != 0 {
			return fmt.Errorf("Error leave router: %s", *leaveRouterOutput.Message)
		}
		if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, routerID); err != nil {
			return err
		}
	}
	input := new(qc.DeleteVxNetsInput)
	input.VxNets = []*string{qc.String(d.Id())}
	err = input.Validate()
	if err != nil {
		return fmt.Errorf("Error delete vxnet input validate: %s", err)
	}
	output, err := clt.DeleteVxNets(input)
	if err != nil {
		return fmt.Errorf("Error delete vxnet: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error delete vxnet: %s", *output.Message)
	}
	d.SetId("")
	return nil
}

func resourceQingcloudVxnetUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("router_id") || d.HasChange("ip_network") {
		routerClt := meta.(*QingCloudClient).router
		routerID := d.Get("router_id").(string)
		IPNetwork := d.Get("ip_network").(string)
		if (routerID != "" && IPNetwork == "") || (routerID == "" && IPNetwork != "") {
			return errors.New("router and ip_network must both be empty or no empty at the same time")
		}
		oldV, newV := d.GetChange("router_id")
		oldRouterID := oldV.(string)
		newRouterID := newV.(string)
		if oldRouterID == "" {
			// do a join router action
			if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, newRouterID); err != nil {
				return err
			}
			joinRouterInput := new(qc.JoinRouterInput)
			joinRouterInput.VxNet = qc.String(d.Id())
			joinRouterInput.Router = qc.String(newRouterID)
			joinRouterInput.IPNetwork = qc.String(IPNetwork)
			joinRouterOutput, err := routerClt.JoinRouter(joinRouterInput)
			if err != nil {
				return fmt.Errorf("Error create vxnet join router: %s", err)
			}
			if joinRouterOutput.RetCode != nil && qc.IntValue(joinRouterOutput.RetCode) != 0 {
				return fmt.Errorf("Error create vxnet join router: %s", *joinRouterOutput.Message)
			}
			if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, oldRouterID); err != nil {
				return err
			}
		} else if newRouterID == "" {
			// do a leave router action
			if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, oldRouterID); err != nil {
				return err
			}
			leaveRouterInput := new(qc.LeaveRouterInput)
			leaveRouterInput.Router = qc.String(oldRouterID)
			leaveRouterInput.VxNets = []*string{qc.String(d.Id())}
			err := leaveRouterInput.Validate()
			if err != nil {
				return fmt.Errorf("Error leave router input validate: %s", err)
			}
			leaveRouterOutput, err := routerClt.LeaveRouter(leaveRouterInput)
			if err != nil {
				return fmt.Errorf("Error leave router: %s", err)
			}
			if leaveRouterOutput.RetCode != nil && qc.IntValue(leaveRouterOutput.RetCode) != 0 {
				return fmt.Errorf("Error leave router: %s", *leaveRouterOutput.Message)
			}
			if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, oldRouterID); err != nil {
				return err
			}
		} else {
			// do a leave router then do a  join router action
			// leave router
			if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, oldRouterID); err != nil {
				return err
			}
			leaveRouterInput := new(qc.LeaveRouterInput)
			leaveRouterInput.Router = qc.String(oldRouterID)
			leaveRouterInput.VxNets = []*string{qc.String(d.Id())}
			err := leaveRouterInput.Validate()
			if err != nil {
				return fmt.Errorf("Error leave router input validate: %s", err)
			}
			leaveRouterOutput, err := routerClt.LeaveRouter(leaveRouterInput)
			if err != nil {
				return fmt.Errorf("Error leave router: %s", err)
			}
			if leaveRouterOutput.RetCode != nil && qc.IntValue(leaveRouterOutput.RetCode) != 0 {
				return fmt.Errorf("Error leave router: %s", *leaveRouterOutput.Message)
			}
			if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, oldRouterID); err != nil {
				return err
			}
			// join router
			if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, newRouterID); err != nil {
				return err
			}
			joinRouterInput := new(qc.JoinRouterInput)
			joinRouterInput.VxNet = qc.String(d.Id())
			joinRouterInput.Router = qc.String(newRouterID)
			joinRouterInput.IPNetwork = qc.String(IPNetwork)
			joinRouterOutput, err := routerClt.JoinRouter(joinRouterInput)
			if err != nil {
				return fmt.Errorf("Error create vxnet join router: %s", err)
			}
			if joinRouterOutput.RetCode != nil && qc.IntValue(joinRouterOutput.RetCode) != 0 {
				return fmt.Errorf("Error create vxnet join router: %s", *joinRouterOutput.Message)
			}
			if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, newRouterID); err != nil {
				return err
			}
		}
		d.Set("router_id", routerID)
		d.Set("ip_network", IPNetwork)
	}
	err := modifyVxnetAttributes(d, meta, false)
	if err != nil {
		return err
	}
	return resourceQingcloudVxnetRead(d, meta)
}
