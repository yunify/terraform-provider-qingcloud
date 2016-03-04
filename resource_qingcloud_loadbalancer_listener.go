package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/loadbalancer"
)

func resourceQingcloudLoadbalancerListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudLoadbalancerListenerCreate,
		Read:   resourceQingcloudLoadbalancerListenerRead,
		Update: resourceQingcloudLoadbalancerListenerUpdate,
		Delete: resourceQingcloudLoadbalancerListenerDelete,
		Schema: map[string]*schema.Schema{
			"loadbalancer": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"protocol": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"certificate": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"session_sticky": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"forwardfor": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"health_check_method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_option": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_option": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceQingcloudLoadbalancerListenerCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	params := loadbalancer.AddLoadBalancerListenersRequest{}
	params.Loadbalancer.Set(d.Get("loadbalancer").(string))
	params.ListenersNListenerPort.Add(int64(d.Get("port").(int)))
	params.ListenersNListenerProtocol.Add(d.Get("protocol").(string))
	params.ListenersNServerCertificateId.Add(d.Get("certificate").(string))
	params.ListenersNBackendProtocol.Add(d.Get("protocol").(string))
	params.ListenersNLoadbalancerListenerName.Add(d.Get("name").(string))
	params.ListenersNBalanceMode.Add(d.Get("mode").(string))
	params.ListenersNSessionSticky.Add(d.Get("session_sticky").(string))
	params.ListenersNForwardfor.Add(int64(d.Get("forwardfor").(int)))
	params.ListenersNHealthyCheckMethod.Add(d.Get("health_check_method").(string))
	params.ListenersNHealthyCheckOption.Add(d.Get("health_check_option").(string))
	resp, err := clt.AddLoadBalancerListeners(params)
	if err != nil {
		return err
	}
	lb := resp.LoadbalancerListeners[0]
	d.SetId(lb)
	return nil
}

func resourceQingcloudLoadbalancerListenerRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	params := loadbalancer.DescribeLoadBalancerListenersRequest{}
	params.LoadbalancerListenersN.Add(d.Id())
	resp, err := clt.DescribeLoadBalancerListeners(params)
	if err != nil {
		return err
	}
	lb := resp.LoadbalancerListenerSet[0]
	d.Set("loadbalancer", lb.LoadbalancerID)
	d.Set("port", lb.ListenerPort)
	d.Set("protocol", lb.ListenerProtocol)
	d.Set("name", lb.LoadbalancerListenerName)
	d.Set("mode", lb.BalanceMode)
	d.Set("session_sticky", lb.SessionSticky)
	d.Set("forwardfor", lb.Forwardfor)
	d.Set("health_check_method", lb.HealthyCheckMethod)
	d.Set("health_check_option", lb.HealthyCheckOption)
	return nil
}

func resourceQingcloudLoadbalancerListenerDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	params := loadbalancer.DeleteLoadBalancerListenersRequest{}
	params.LoadbalancerListenersN.Add(d.Id())
	_, err := clt.DeleteLoadBalancerListeners(params)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceQingcloudLoadbalancerListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
