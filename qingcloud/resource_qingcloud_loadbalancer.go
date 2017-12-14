package qingcloud

import (
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
		Create: resourceQingcloudVpcCreate,
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
	if d.HasChange(resourceLoadBalancerEipIDs) {
		if err := updateLoadbalancerEips(d, meta); err != nil {
			return err
		}
	}
	if d.HasChange(resourceLoadBalancerType) {
		if err := resizeLoadBalancer(qc.String(d.Id()), qc.Int(d.Get(resourceLoadBalancerType).(int)), meta); err != nil {
			return err
		}
	}
	return resourceQingcloudLoadBalancerRead(d, meta)
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
