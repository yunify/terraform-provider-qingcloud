package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	resourceLoadBalancerType            = "type"
	resourceLoadBalancerPrivateIPs      = "private_ips"
	resourceLoadBalancerEipIDs          = "eip_ids"
	resourceLoadBalancerNodeCount       = "node_count"
	resourceLoadBalancerSecurityGroupID = "security_group_id"
	resourceLoadBalancerVxnetID         = "vxnet_id"
	resourceLoadBalancerHttpHeaderSize  = "http_header_size"
)

func resourceQingcloudLoadBalancer() *schema.Resource {

	return &schema.Resource{
		Create: resourceQingcloudLoadBalancerCreate,
		Read:   resourceQingcloudLoadBalancerRead,
		Update: resourceQingcloudLoadBalancerUpdate,
		Delete: resourceQingcloudLoadBalancerDelete,
		Schema: map[string]*schema.Schema{
			resourceName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceLoadBalancerType: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: withinArrayInt(0, 1, 2, 3, 4, 5),
			},
			resourceLoadBalancerPrivateIPs: {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			resourceLoadBalancerEipIDs: {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			resourceLoadBalancerNodeCount: {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			resourceLoadBalancerSecurityGroupID: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			resourceLoadBalancerVxnetID: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  BasicNetworkID,
				ForceNew: true,
			},
			resourceLoadBalancerHttpHeaderSize: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      15,
				ValidateFunc: withinArrayIntRange(1, 127),
			},
			resourceTagIds:   tagIdsSchema(),
			resourceTagNames: tagNamesSchema(),
		},
	}
}
func resourceQingcloudLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := waitLoadBalancerLease(d, meta); err != nil {
		return err
	}
	d.Partial(true)
	if err := modifyLoadBalancerAttributes(d, meta); err != nil {
		return err
	}
	d.SetPartial(resourceLoadBalancerPrivateIPs)
	d.SetPartial(resourceLoadBalancerHttpHeaderSize)
	d.SetPartial(resourceLoadBalancerSecurityGroupID)
	d.SetPartial(resourceLoadBalancerNodeCount)
	d.SetPartial(resourceName)
	d.SetPartial(resourceDescription)
	if d.HasChange(resourceLoadBalancerEipIDs) {
		if err := updateLoadbalancerEips(d, meta); err != nil {
			return err
		}
	}
	d.SetPartial(resourceLoadBalancerEipIDs)
	if d.HasChange(resourceLoadBalancerType) && !d.IsNewResource() {
		if err := resizeLoadBalancer(qc.String(d.Id()), qc.Int(d.Get(resourceLoadBalancerType).(int)), meta); err != nil {
			return err
		}
	}
	d.SetPartial(resourceLoadBalancerType)
	if err := resourceUpdateTag(d, meta, qingcloudResourceTypeLoadBalancer); err != nil {
		return err
	}
	d.Partial(false)
	return resourceQingcloudLoadBalancerRead(d, meta)
}

func resourceQingcloudLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.CreateLoadBalancerInput)
	input.LoadBalancerName, _ = getNamePointer(d)
	input.VxNet = getSetStringPointer(d, resourceLoadBalancerVxnetID)
	input.SecurityGroup = getSetStringPointer(d, resourceLoadBalancerSecurityGroupID)
	input.HTTPHeaderSize = qc.Int(d.Get(resourceLoadBalancerHttpHeaderSize).(int))
	if d.Get(resourceLoadBalancerNodeCount).(int) != 0 && qc.StringValue(input.VxNet) == BasicNetworkID {
		input.NodeCount = qc.Int(d.Get(resourceLoadBalancerNodeCount).(int))
	}
	input.LoadBalancerType = qc.Int(d.Get(resourceLoadBalancerType).(int))
	if _, ok := d.GetOk(resourceLoadBalancerPrivateIPs); ok {
		privateIPs := d.Get(resourceLoadBalancerPrivateIPs).(*schema.Set).List()
		if len(privateIPs) != 1 || d.Get(resourceLoadBalancerVxnetID).(string) == BasicNetworkID {
			return fmt.Errorf("error private_ips info")
		}
		input.PrivateIP = qc.String(privateIPs[0].(string))
	}
	if qc.StringValue(input.VxNet) == BasicNetworkID {
		var eips []*string
		for _, value := range d.Get(resourceLoadBalancerEipIDs).(*schema.Set).List() {
			eips = append(eips, qc.String(value.(string)))
		}
		input.EIPs = eips
	}
	var output *qc.CreateLoadBalancerOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.CreateLoadBalancer(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.LoadBalancerID))
	if _, err = LoadBalancerTransitionStateRefresh(clt, qc.String(d.Id())); err != nil {
		return err
	}
	return resourceQingcloudLoadBalancerUpdate(d, meta)
}

func resourceQingcloudLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.DescribeLoadBalancersInput)
	input.LoadBalancers = []*string{qc.String(d.Id())}
	input.Verbose = qc.Int(1)
	var output *qc.DescribeLoadBalancersOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeLoadBalancers(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if isLoadBalancerDeleted(output.LoadBalancerSet) {
		d.SetId("")
		return nil
	}
	lb := output.LoadBalancerSet[0]
	d.Set(resourceName, qc.StringValue(lb.LoadBalancerName))
	d.Set(resourceDescription, qc.StringValue(lb.Description))
	d.Set(resourceLoadBalancerType, qc.IntValue(lb.LoadBalancerType))
	d.Set(resourceLoadBalancerVxnetID, qc.StringValue(lb.VxNetID))
	if err := d.Set(resourceLoadBalancerPrivateIPs, qc.StringValueSlice(lb.PrivateIPs)); err != nil {
		return err
	}
	d.Set(resourceLoadBalancerSecurityGroupID, qc.StringValue(lb.SecurityGroupID))
	d.Set(resourceLoadBalancerNodeCount, qc.IntValue(lb.NodeCount))
	var eipIDs []string
	if d.Get(resourceLoadBalancerVxnetID) == BasicNetworkID {
		for _, eip := range lb.Cluster {
			eipIDs = append(eipIDs, qc.StringValue(eip.EIPID))
		}
	} else {
		for _, eip := range lb.EIPs {
			eipIDs = append(eipIDs, qc.StringValue(eip.EIPID))
		}
	}
	d.Set(resourceLoadBalancerEipIDs, eipIDs)
	if err := resourceSetTag(d, lb.Tags); err != nil {
		return err
	}
	return nil
}

func resourceQingcloudLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	if _, err := LoadBalancerTransitionStateRefresh(clt, qc.String(d.Id())); err != nil {
		return err
	}
	if err := waitLoadBalancerLease(d, meta); err != nil {
		return err
	}
	input := new(qc.DeleteLoadBalancersInput)
	input.LoadBalancers = []*string{qc.String(d.Id())}
	var err error
	simpleRetry(func() error {
		_, err = clt.DeleteLoadBalancers(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if _, err := LoadBalancerTransitionStateRefresh(clt, qc.String(d.Id())); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
