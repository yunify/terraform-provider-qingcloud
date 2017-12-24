package qingcloud

import "github.com/hashicorp/terraform/helper/schema"

const (
	resourceLoadBalancerListenerLBId                = "load_balancer_id"
	resourceLoadBalancerListenerPort                = "listener_port"
	resourceLoadBalancerListenerProtocol            = "listener_protocol"
	resourceLoadBalancerListenerServerCertificateId = "server_certificate_id"
	resourceLoadBalancerListenerBalancerMode        = "balance_mode"
	resourceLoadBalancerListenerSessionSticky       = "session_sticky"
	resourceLoadBalancerListenerForwardfor          = "forwardfor"
	resourceLoadBalancerListenerHealthCheckMethod   = "healthy_check_method"
	resourceLoadBalancerListenerHealthCheckOption   = "healthy_check_option"
	resourceLoadBalancerListenerOption              = "listener_option"
)

func resourceQingcloudLoadBalancerListener() *schema.Resource {
	return &schema.Resource{
		Create: ,
		Read:   ,
		Update: ,
		Delete: ,
		Schema: map[string]*schema.Schema{
			resourceName: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceLoadBalancerListenerLBId: &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			resourceLoadBalancerListenerPort: &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: withinArrayIntRange(1, 65536),
			},
			resourceLoadBalancerListenerProtocol: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: withinArrayString("http", "https", "tcp", "ssl"),
			},
			resourceLoadBalancerListenerBalancerMode: &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "roundrobin",
				ValidateFunc: withinArrayString("roundrobin", "leastconn", "source"),
			},
			resourceLoadBalancerListenerServerCertificateId: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceLoadBalancerListenerSessionSticky: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			resourceLoadBalancerListenerForwardfor: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			resourceLoadBalancerListenerHealthCheckMethod: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			resourceLoadBalancerListenerHealthCheckOption: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			resourceLoadBalancerListenerOption: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}
