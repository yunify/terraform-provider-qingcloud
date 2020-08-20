package qingcloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	resourceVxnetType         = "type"
	resourceVxnetVpcID        = "vpc_id"
	resourceVxnetVpcIPNetwork = "ip_network"
)

func resourceQingcloudVxnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudVxnetCreate,
		Read:   resourceQingcloudVxnetRead,
		Update: resourceQingcloudVxnetUpdate,
		Delete: resourceQingcloudVxnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			resourceName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceVxnetType: {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: withinArrayInt(0, 1),
			},
			resourceDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceVxnetVpcID: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceVxnetVpcIPNetwork: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNetworkCIDR,
			},
			resourceTagIds:   tagIdsSchema(),
			resourceTagNames: tagNamesSchema(),
		},
	}
}

func resourceQingcloudVxnetCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet
	input := new(qc.CreateVxNetsInput)
	input.Count = qc.Int(1)
	input.VxNetName, _ = getNamePointer(d)
	input.VxNetType = qc.Int(d.Get(resourceVxnetType).(int))
	var output *qc.CreateVxNetsOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.CreateVxNets(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
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
	simpleRetry(func() error {
		output, err = clt.DescribeVxNets(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if len(output.VxNetSet) == 0 {
		d.SetId("")
		return nil
	}
	vxnet := output.VxNetSet[0]
	d.Set(resourceName, qc.StringValue(vxnet.VxNetName))
	d.Set(resourceVxnetType, qc.IntValue(vxnet.VxNetType))
	d.Set(resourceDescription, qc.StringValue(vxnet.Description))
	if vxnet.Router != nil {
		d.Set(resourceVxnetVpcIPNetwork, qc.StringValue(vxnet.Router.IPNetwork))
	} else {
		d.Set(resourceVxnetVpcIPNetwork, "")
	}
	d.Set(resourceVxnetVpcID, qc.StringValue(vxnet.VpcRouterID))
	if err := resourceSetTag(d, vxnet.Tags); err != nil {
		return err
	}
	return nil
}

func resourceQingcloudVxnetUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if d.HasChange(resourceVxnetVpcID) || d.HasChange(resourceVxnetVpcIPNetwork) {
		vpcID := d.Get(resourceVxnetVpcID).(string)
		IPNetwork := d.Get(resourceVxnetVpcIPNetwork).(string)
		if (vpcID != "" && IPNetwork == "") || (vpcID == "" && IPNetwork != "") {
			return errors.New("vpc_id and ip_network must both be empty or no empty at the same time")
		} else if d.Get(resourceVxnetType).(int) == 0 {
			return fmt.Errorf("vpc_id and ip_network can be set in Managed vxnet")
		}
		oldVPC, newVPC := d.GetChange(resourceVxnetVpcID)
		oldVPCID := oldVPC.(string)
		newVPCID := newVPC.(string)
		if oldVPCID == "" {
			// do a join router action
			if err := vxnetJoinRouter(d, meta); err != nil {
				return err
			}
		} else if newVPCID == "" {
			// do a leave router action
			if err := vxnetLeaverRouter(d, meta); err != nil {
				return err
			}
		} else {
			// do a leave router then do a  join router action
			if err := vxnetLeaverRouter(d, meta); err != nil {
				return err
			}
			if err := vxnetJoinRouter(d, meta); err != nil {
				return err
			}

		}
	}
	d.SetPartial(resourceVxnetVpcID)
	d.SetPartial(resourceVxnetVpcIPNetwork)
	if err := modifyVxnetAttributes(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceName)
	d.SetPartial(resourceDescription)
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeVxNet); err != nil {
		return err
	}
	d.SetPartial(resourceTagIds)
	d.Partial(false)
	return resourceQingcloudVxnetRead(d, meta)
}

func resourceQingcloudVxnetDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).vxnet
	vpcID := d.Get(resourceVxnetVpcID).(string)
	// vxnet leave router
	if vpcID != "" {
		if err := vxnetLeaverRouter(d, meta); err != nil {
			return err
		}
	}
	input := new(qc.DeleteVxNetsInput)
	input.VxNets = []*string{qc.String(d.Id())}
	var err error
	simpleRetry(func() error {
		_, err = clt.DeleteVxNets(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if _, err := RouterTransitionStateRefresh(meta.(*QingCloudClient).router, vpcID); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
