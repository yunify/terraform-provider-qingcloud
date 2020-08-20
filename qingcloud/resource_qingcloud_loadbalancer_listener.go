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
	resourceLoadBalancerListenerTimeOut             = "timeout"
)

func resourceQingcloudLoadBalancerListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudLoadBalancerListnerCreate,
		Read:   resourceQingcloudLoadBalancerListenerRead,
		Update: resourceQingcloudLoadBalancerListenerUpdate,
		Delete: resourceQingcloudLoadBalancerListnerDestroy,
		Schema: map[string]*schema.Schema{
			resourceName: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resourceLoadBalancerListenerLBId: {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			resourceLoadBalancerListenerPort: {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: withinArrayIntRange(1, 65536),
			},
			resourceLoadBalancerListenerProtocol: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: withinArrayString("http", "https", "tcp", "ssl"),
			},
			resourceLoadBalancerListenerBalancerMode: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "roundrobin",
				ValidateFunc: withinArrayString("roundrobin", "leastconn", "source"),
			},
			resourceLoadBalancerListenerServerCertificateId: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			resourceLoadBalancerListenerSessionSticky: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			resourceLoadBalancerListenerForwardfor: {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: withinArrayIntRange(0, 7),
			},
			resourceLoadBalancerListenerHealthCheckMethod: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tcp",
			},
			resourceLoadBalancerListenerHealthCheckOption: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10|5|2|5",
			},
			resourceLoadBalancerListenerOption: {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: withinArrayIntRange(0, 1023),
			},
			resourceLoadBalancerListenerTimeOut: {
				Type:     schema.TypeInt,
				Default:  50,
				Optional: true,
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
	if len(d.Get(resourceLoadBalancerListenerServerCertificateId).(*schema.Set).List()) > 0 {
		for _, value := range d.Get(resourceLoadBalancerListenerServerCertificateId).(*schema.Set).List() {
			listener.ServerCertificateID = append(listener.ServerCertificateID, qc.String(value.(string)))
		}
	}
	listener.SessionSticky = getSetStringPointer(d, resourceLoadBalancerListenerSessionSticky)
	listener.Forwardfor = qc.Int(d.Get(resourceLoadBalancerListenerForwardfor).(int))
	listener.HealthyCheckMethod = getSetStringPointer(d, resourceLoadBalancerListenerHealthCheckMethod)
	listener.HealthyCheckOption = getSetStringPointer(d, resourceLoadBalancerListenerHealthCheckOption)
	listener.ListenerOption = qc.Int(d.Get(resourceLoadBalancerListenerOption).(int))
	listener.Timeout = qc.Int(d.Get(resourceLoadBalancerListenerTimeOut).(int))

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
	if err := d.Set(resourceLoadBalancerListenerServerCertificateId, qc.StringValueSlice(output.LoadBalancerListenerSet[0].ServerCertificateID)); err != nil {
		return err
	}
	d.Set(resourceLoadBalancerListenerSessionSticky, qc.StringValue(output.LoadBalancerListenerSet[0].SessionSticky))
	d.Set(resourceLoadBalancerListenerForwardfor, qc.IntValue(output.LoadBalancerListenerSet[0].Forwardfor))
	d.Set(resourceLoadBalancerListenerHealthCheckMethod, qc.StringValue(output.LoadBalancerListenerSet[0].HealthyCheckMethod))
	d.Set(resourceLoadBalancerListenerHealthCheckOption, qc.StringValue(output.LoadBalancerListenerSet[0].HealthyCheckOption))
	d.Set(resourceLoadBalancerListenerOption, qc.IntValue(output.LoadBalancerListenerSet[0].ListenerOption))
	d.Set(resourceLoadBalancerListenerTimeOut, qc.IntValue(output.LoadBalancerListenerSet[0].Timeout))
	return nil
}

func resourceQingcloudLoadBalancerListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.ModifyLoadBalancerListenerAttributesInput)
	input.LoadBalancerListener = qc.String(d.Id())
	input.LoadBalancerListenerName = getSetStringPointer(d, resourceName)
	input.BalanceMode = getSetStringPointer(d, resourceLoadBalancerListenerBalancerMode)
	if d.HasChange(resourceLoadBalancerListenerServerCertificateId) {
		if len(d.Get(resourceLoadBalancerListenerServerCertificateId).(*schema.Set).List()) == 0 {
			input.ServerCertificateID = append(input.ServerCertificateID, qc.String(" "))
		}
		for _, value := range d.Get(resourceLoadBalancerListenerServerCertificateId).(*schema.Set).List() {
			input.ServerCertificateID = append(input.ServerCertificateID, qc.String(value.(string)))
		}
	}
	input.SessionSticky = getSetStringPointer(d, resourceLoadBalancerListenerSessionSticky)
	input.Forwardfor = qc.Int(d.Get(resourceLoadBalancerListenerForwardfor).(int))
	input.HealthyCheckMethod = getSetStringPointer(d, resourceLoadBalancerListenerHealthCheckMethod)
	input.HealthyCheckOption = getSetStringPointer(d, resourceLoadBalancerListenerHealthCheckOption)
	input.ListenerOption = qc.Int(d.Get(resourceLoadBalancerListenerOption).(int))
	input.Timeout = qc.Int(d.Get(resourceLoadBalancerListenerTimeOut).(int))
	var err error
	simpleRetry(func() error {
		_, err = clt.ModifyLoadBalancerListenerAttributes(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if err := updateLoadBalancer(qc.String(d.Get(resourceLoadBalancerListenerLBId).(string)), meta); err != nil {
		return nil
	}
	return resourceQingcloudLoadBalancerListenerRead(d, meta)
}

func resourceQingcloudLoadBalancerListnerDestroy(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	input := new(qc.DeleteLoadBalancerListenersInput)
	input.LoadBalancerListeners = []*string{qc.String(d.Id())}
	var err error
	simpleRetry(func() error {
		_, err = clt.DeleteLoadBalancerListeners(input)
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
