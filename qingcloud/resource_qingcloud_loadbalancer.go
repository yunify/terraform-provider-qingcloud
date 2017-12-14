package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
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
		Read:   resourceQingcloudVpcRead,
		Update: resourceQingcloudVpcUpdate,
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
