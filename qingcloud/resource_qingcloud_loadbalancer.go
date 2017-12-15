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
		Delete: resourceQingcloudVpcDelete,
		Schema: map[string]*schema.Schema{
			resourceName: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceDescription: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceLoadBalancerType: &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: withinArrayInt(0, 1, 2, 3, 4, 5),
			},
			resourceLoadBalancerPrivateIPs: &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			resourceLoadBalancerEipIDs: &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			resourceLoadBalancerNodeCount: &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
			resourceLoadBalancerSecurityGroupID: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			resourceLoadBalancerVxnetID: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "vxnet-0",
				ForceNew: true,
			},
			resourceLoadBalancerHttpHeaderSize: &schema.Schema{
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
	if err := modifyLoadBalancerAttributes(d, meta); err != nil {
		return err
	}
	if d.HasChange(resourceLoadBalancerEipIDs) && !d.IsNewResource() {
		if err := updateLoadbalancerEips(d, meta); err != nil {
			return err
		}
	}
	if d.HasChange(resourceLoadBalancerType) && !d.IsNewResource() {
		if err := resizeLoadBalancer(qc.String(d.Id()), qc.Int(d.Get(resourceLoadBalancerType).(int)), meta); err != nil {
			return err
		}
	}
	return resourceQingcloudLoadBalancerRead(d, meta)
}

func resourceQingcloudLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.CreateLoadBalancerInput)
	input.LoadBalancerName, _ = getNamePointer(d)
	input.VxNet = getSetStringPointer(d, resourceLoadBalancerVxnetID)
	input.SecurityGroup = getSetStringPointer(d, resourceLoadBalancerSecurityGroupID)
	input.HTTPHeaderSize = qc.Int(d.Get(resourceLoadBalancerHttpHeaderSize).(int))
	input.NodeCount = qc.Int(d.Get(resourceLoadBalancerNodeCount).(int))
	input.LoadBalancerType = qc.Int(d.Get(resourceLoadBalancerType).(int))
	if _, ok := d.GetOk(resourceLoadBalancerPrivateIPs); ok {
		privateIPs := d.Get(resourceLoadBalancerPrivateIPs).(*schema.Set).List()
		if len(privateIPs) != 1 || d.Get(resourceLoadBalancerVxnetID).(string) == "vxnet-0" {
			return fmt.Errorf("error private_ips info")
		}
		input.PrivateIP = qc.String(privateIPs[0].(string))
	}
	var eips []*string
	for _, value := range d.Get(resourceLoadBalancerEipIDs).(*schema.Set).List() {
		eips = append(eips, qc.String(value.(string)))
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
	return resourceQingcloudVpcUpdate(d, meta)
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
	d.Set(resourceLoadBalancerPrivateIPs, qc.StringValueSlice(lb.PrivateIPs))
	d.Set(resourceLoadBalancerSecurityGroupID, qc.StringValue(lb.SecurityGroupID))
	d.Set(resourceLoadBalancerNodeCount, qc.IntValue(lb.NodeCount))
	var eipIDs []string
	for _, eip := range lb.Cluster {
		eipIDs = append(eipIDs, qc.StringValue(eip.EIPID))
	}
	d.Set(resourceLoadBalancerEipIDs, eipIDs)
	resourceSetTag(d, lb.Tags)
	return nil
}
