package qingcloud

import (
	"errors"
	"fmt"

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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of vxnet",
			},
			"type": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				Description:  "type of vxnet,1 - Managed vxnet,0 - Self-managed vxnet.",
				ValidateFunc: withinArrayInt(0, 1),
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of vxnet",
			},
			"vpc_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vpc id , vxnet want to join.",
			},
			"ip_network": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNetworkCIDR,
				Description:  "Network segment of Managed vxnet",
			},
			"tag_ids": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "tag ids , vxnet wants to use",
			},
			"tag_names": &schema.Schema{
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "compute by tag ids",
			},
		},
	}
}

func resourceQingcloudVxnetCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet
	input := new(qc.CreateVxNetsInput)
	input.Count = qc.Int(1)
	if d.Get("name").(string) != "" {
		input.VxNetName = qc.String(d.Get("name").(string))
	}
	input.VxNetType = qc.Int(d.Get("type").(int))
	var output *qc.CreateVxNetsOutput
	var err error
	retryServerBusy(func() (s *int, err error) {
		output, err = clt.CreateVxNets(input)
		return output.RetCode, err
	})
	if err != nil {
		return fmt.Errorf("Error create vxnet: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error create vxnet: %s", *output.Message)
	}
	d.SetId(qc.StringValue(output.VxNets[0]))
	return resourceQingcloudVxnetUpdate(d, meta)
}

func resourceQingcloudVxnetRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet
	input := new(qc.DescribeVxNetsInput)
	input.VxNets = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	var output *qc.DescribeVxNetsOutput
	var err error
	retryServerBusy(func() (s *int, err error) {
		output, err = clt.DescribeVxNets(input)
		return output.RetCode, err
	})
	if err != nil {
		return fmt.Errorf("Error describe vxnet: %s", err)
	}
	if output.RetCode == nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error describe vxnet: %s", *output.Message)
	}
	if len(output.VxNetSet) == 0 {
		d.SetId("")
		return nil
	}
	vxnet := output.VxNetSet[0]
	d.Set("name", qc.StringValue(vxnet.VxNetName))
	d.Set("type", qc.IntValue(vxnet.VxNetType))
	d.Set("description", qc.StringValue(vxnet.Description))
	if vxnet.Router != nil {
		d.Set("ip_network", qc.StringValue(vxnet.Router.IPNetwork))
	} else {
		d.Set("ip_network", "")
	}
	d.Set("vpc_id", qc.StringValue(vxnet.VpcRouterID))
	resourceSetTag(d, vxnet.Tags)
	return nil
}

func resourceQingcloudVxnetUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("vpc_id") || d.HasChange("ip_network") {
		vpcID := d.Get("vpc_id").(string)
		IPNetwork := d.Get("ip_network").(string)
		if (vpcID != "" && IPNetwork == "") || (vpcID == "" && IPNetwork != "") {
			return errors.New("vpc_id and ip_network must both be empty or no empty at the same time")
		} else if d.Get("type").(int) == 0 {
			return fmt.Errorf("vpc_id and ip_network can be set in Managed vxnet")
		}
		oldVPC, newVPC := d.GetChange("vpc_id")
		oldVPCID := oldVPC.(string)
		newVPCID := newVPC.(string)
		if oldVPCID == "" {
			// do a join router action
			if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, newVPCID); err != nil {
				return err
			}
			if err := vxnetJoinRouter(d, meta); err != nil {
				return err
			}
		} else if newVPCID == "" {
			// do a leave router action
			if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, oldVPCID); err != nil {
				return err
			}
			if err := vxnetLeaverRouter(d, meta); err != nil {
				return err
			}
		} else {
			// do a leave router then do a  join router action
			if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, oldVPCID); err != nil {
				return err
			}
			if err := vxnetLeaverRouter(d, meta); err != nil {
				return err
			}
			if err := vxnetJoinRouter(d, meta); err != nil {
				return err
			}

		}
	}
	if err := modifyVxnetAttributes(d, meta); err != nil {
		return err
	}
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeVxNet); err != nil {
		return err
	}
	return resourceQingcloudVxnetRead(d, meta)
}

func resourceQingcloudVxnetDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet
	if _, err := VxnetTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	vpcID := d.Get("vpc_id").(string)
	// vxnet leave router
	if vpcID != "" {
		if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, vpcID); err != nil {
			return err
		}
		if err := vxnetLeaverRouter(d, meta); err != nil {
			return err
		}
	}
	input := new(qc.DeleteVxNetsInput)
	input.VxNets = []*string{qc.String(d.Id())}
	var output *qc.DeleteVxNetsOutput
	var err error
	retryServerBusy(func() (s *int, err error) {
		output, err = clt.DeleteVxNets(input)
		return output.RetCode, err
	})
	if err != nil {
		return fmt.Errorf("Error delete vxnet: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error delete vxnet: %s", *output.Message)
	}
	if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, vpcID); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
