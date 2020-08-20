package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	resourceVpcType            = "type"
	resourceVpcNetwork         = "vpc_network"
	resourceVpcEipID           = "eip_id"
	resourceVpcSecurityGroupID = "security_group_id"
	resourceVpcPrivateIP       = "private_ip"
	resourceVpcPublicIP        = "public_ip"
)

func resourceQingcloudVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudVpcCreate,
		Read:   resourceQingcloudVpcRead,
		Update: resourceQingcloudVpcUpdate,
		Delete: resourceQingcloudVpcDelete,
		Schema: map[string]*schema.Schema{
			resourceName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceVpcType: {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: withinArrayInt(0, 1, 2, 3),
			},
			resourceVpcNetwork: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: withinArrayString("192.168.0.0/16", "172.16.0.0/16", "172.17.0.0/16",
					"172.18.0.0/16", "172.19.0.0/16", "172.20.0.0/16", "172.21.0.0/16", "172.22.0.0/16",
					"172.23.0.0/16", "172.24.0.0/16", "172.25.0.0/16", "172.26.0.0/16", "172.27.0.0/16",
					"172.28.0.0/16", "172.29.0.0/16", "172.30.0.0/16", "172.31.0.0/16"),
			},
			resourceTagIds:   tagIdsSchema(),
			resourceTagNames: tagNamesSchema(),
			resourceVpcEipID: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceVpcSecurityGroupID: {
				Type:     schema.TypeString,
				Required: true,
			},
			resourceDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceVpcPrivateIP: {
				Type:     schema.TypeString,
				Computed: true,
			},
			resourceVpcPublicIP: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceQingcloudRouterCreate
func resourceQingcloudVpcCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	input := new(qc.CreateRoutersInput)
	input.RouterName, _ = getNamePointer(d)
	input.VpcNetwork = getSetStringPointer(d, resourceVpcNetwork)
	input.SecurityGroup = getSetStringPointer(d, resourceVpcSecurityGroupID)
	input.RouterType = qc.Int(d.Get(resourceVpcType).(int))
	input.Count = qc.Int(1)
	var output *qc.CreateRoutersOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.CreateRouters(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.Routers[0]))
	if _, err = RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return fmt.Errorf("Error waiting for router (%s) to start: %s", d.Id(), err.Error())
	}
	return resourceQingcloudVpcUpdate(d, meta)
}

func resourceQingcloudVpcRead(d *schema.ResourceData, meta interface{}) error {
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
	if isRouterDeleted(output.RouterSet) {
		d.SetId("")
		return nil
	}
	rtr := output.RouterSet[0]
	d.Set(resourceName, qc.StringValue(rtr.RouterName))
	d.Set(resourceVpcType, qc.IntValue(rtr.RouterType))
	d.Set(resourceVpcSecurityGroupID, qc.StringValue(rtr.SecurityGroupID))
	d.Set(resourceDescription, qc.StringValue(rtr.Description))
	d.Set(resourceVpcPrivateIP, qc.StringValue(rtr.PrivateIP))
	d.Set(resourceVpcEipID, qc.StringValue(rtr.EIP.EIPID))
	d.Set(resourceVpcPublicIP, qc.StringValue(rtr.EIP.EIPAddr))
	if err := resourceSetTag(d, rtr.Tags); err != nil {
		return err
	}
	return nil
}

func resourceQingcloudVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	if _, err := RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	if err := waitRouterLease(d, meta); err != nil {
		return err
	}
	d.Partial(true)
	if err := modifyRouterAttributes(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceName)
	d.SetPartial(resourceDescription)
	if d.HasChange(resourceVpcEipID) {
		if err := applyRouterUpdate(qc.String(d.Id()), meta); err != nil {
			return err
		}
	}
	d.SetPartial(resourceVpcEipID)
	if d.HasChange(resourceVpcSecurityGroupID) && !d.IsNewResource() {
		if err := applySecurityGroupRule(qc.String(resourceVpcSecurityGroupID), meta); err != nil {
			return err
		}
	}
	d.SetPartial(resourceVpcSecurityGroupID)
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeRouter); err != nil {
		return err
	}
	d.SetPartial(resourceTagIds)
	d.Partial(false)
	return resourceQingcloudVpcRead(d, meta)
}

func resourceQingcloudVpcDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).router
	if _, err := RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	if err := waitRouterLease(d, meta); err != nil {
		return err
	}
	input := new(qc.DeleteRoutersInput)
	input.Routers = []*string{qc.String(d.Id())}
	var err error
	simpleRetry(func() error {
		_, err = clt.DeleteRouters(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if _, err := RouterTransitionStateRefresh(clt, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
