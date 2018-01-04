package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

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
		Create: resourceQingcloudLoadBalancerListnerCreate,
		Read:   resourceQingcloudLoadBalancerListenerRead,
		Update: resourceQingcloudLoadBalancerListenerUpdate,
		Delete: resourceQingcloudLoadBalancerListnerDestroy,
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
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: withinArrayIntRange(0, 7),
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
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: withinArrayIntRange(0, 15),
			},
		},
	}
}

func resourceQingcloudLoadBalancerListnerCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.AddLoadBalancerListenersInput)
	listener := new(qc.LoadBalancerListener)
	input.LoadBalancer = getSetStringPointer(d, resourceLoadBalancerListenerLBId)
	listener.LoadBalancerListenerName = getSetStringPointer(d, resourceName)
	listener.ListenerPort = qc.Int(d.Get(resourceLoadBalancerListenerPort).(int))
	listener.ListenerProtocol = getSetStringPointer(d, resourceLoadBalancerListenerProtocol)
	listener.BackendProtocol = getSetStringPointer(d, resourceLoadBalancerListenerProtocol)
	listener.BalanceMode = getSetStringPointer(d, resourceLoadBalancerListenerBalancerMode)
	listener.ServerCertificateID = getSetStringPointer(d, resourceLoadBalancerListenerServerCertificateId)
	listener.SessionSticky = getSetStringPointer(d, resourceLoadBalancerListenerSessionSticky)
	listener.Forwardfor = qc.Int(d.Get(resourceLoadBalancerListenerForwardfor).(int))
	listener.HealthyCheckMethod = getSetStringPointer(d, resourceLoadBalancerListenerHealthCheckMethod)
	listener.HealthyCheckOption = getSetStringPointer(d, resourceLoadBalancerListenerHealthCheckOption)
	listener.ListenerOption = qc.Int(d.Get(resourceLoadBalancerListenerOption).(int))

	input.Listeners = []*qc.LoadBalancerListener{listener}
	var output *qc.AddLoadBalancerListenersOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.AddLoadBalancerListeners(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if err := updateLoadBalancer(qc.String(d.Get(resourceLoadBalancerListenerLBId).(string)), meta); err != nil {
		return nil
	}
	d.SetId(qc.StringValue(output.LoadBalancerListeners[0]))
	return resourceQingcloudLoadBalancerListenerRead(d, meta)
}

func resourceQingcloudLoadBalancerListenerRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.DescribeLoadBalancerListenersInput)
	input.LoadBalancerListeners = []*string{qc.String(d.Id())}
	input.LoadBalancer = getSetStringPointer(d, resourceLoadBalancerListenerLBId)
	var output *qc.DescribeLoadBalancerListenersOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeLoadBalancerListeners(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if len(output.LoadBalancerListenerSet) == 0 {
		d.SetId("")
		return nil
	}
	d.Set(resourceName, qc.StringValue(output.LoadBalancerListenerSet[0].LoadBalancerListenerName))
	d.Set(resourceLoadBalancerListenerPort, qc.IntValue(output.LoadBalancerListenerSet[0].ListenerPort))
	d.Set(resourceLoadBalancerListenerProtocol, qc.StringValue(output.LoadBalancerListenerSet[0].ListenerProtocol))
	d.Set(resourceLoadBalancerListenerBalancerMode, qc.StringValue(output.LoadBalancerListenerSet[0].BalanceMode))
	d.Set(resourceLoadBalancerListenerServerCertificateId, qc.StringValue(output.LoadBalancerListenerSet[0].ServerCertificateID))
	d.Set(resourceLoadBalancerListenerSessionSticky, qc.StringValue(output.LoadBalancerListenerSet[0].SessionSticky))
	d.Set(resourceLoadBalancerListenerForwardfor, qc.IntValue(output.LoadBalancerListenerSet[0].Forwardfor))
	d.Set(resourceLoadBalancerListenerHealthCheckMethod, qc.StringValue(output.LoadBalancerListenerSet[0].HealthyCheckMethod))
	d.Set(resourceLoadBalancerListenerHealthCheckOption, qc.StringValue(output.LoadBalancerListenerSet[0].HealthyCheckOption))
	d.Set(resourceLoadBalancerListenerOption, qc.IntValue(output.LoadBalancerListenerSet[0].ListenerOption))
	return nil
}

func resourceQingcloudLoadBalancerListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.ModifyLoadBalancerListenerAttributesInput)
	input.LoadBalancerListenerName = getSetStringPointer(d, resourceName)
	input.BalanceMode = getSetStringPointer(d, resourceLoadBalancerListenerBalancerMode)
	input.ServerCertificateID = getSetStringPointer(d, resourceLoadBalancerListenerServerCertificateId)
	input.SessionSticky = getSetStringPointer(d, resourceLoadBalancerListenerSessionSticky)
	input.Forwardfor = qc.Int(d.Get(resourceLoadBalancerListenerForwardfor).(int))
	input.HealthyCheckMethod = getSetStringPointer(d, resourceLoadBalancerListenerHealthCheckMethod)
	input.HealthyCheckOption = getSetStringPointer(d, resourceLoadBalancerListenerHealthCheckOption)
	//TODO input.ListenerOption =  qc.Int(d.Get(resourceLoadBalancerListenerOption).(int))
	var output *qc.ModifyLoadBalancerListenerAttributesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.ModifyLoadBalancerListenerAttributes(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if err := updateLoadBalancer(qc.String(d.Get(resourceLoadBalancerListenerLBId).(string)), meta); err != nil {
		return nil
	}
	return resourceQingcloudVpcStaticRead(d, meta)
}

func resourceQingcloudLoadBalancerListnerDestroy(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.DeleteLoadBalancerListenersInput)
	input.LoadBalancerListeners = []*string{qc.String(d.Id())}
	var output *qc.DeleteLoadBalancerListenersOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DeleteLoadBalancerListeners(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if err := updateLoadBalancer(qc.String(d.Get(resourceLoadBalancerListenerLBId).(string)), meta); err != nil {
		return nil
	}
	d.SetId("")
	return nil
}