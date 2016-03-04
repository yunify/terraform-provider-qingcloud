package qingcloud

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/loadbalancer"
)

func resourceQingcloudLoadbalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudLoadbalancerCreate,
		Read:   resourceQingcloudLoadbalancerRead,
		Update: resourceQingcloudLoadbalancerUpdate,
		Delete: nil,
		Schema: map[string]*schema.Schema{
			"eip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vxnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceQingcloudLoadbalancerCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	params := loadbalancer.CreateLoadBalancerRequest{}
	params.EipsN.Add(d.Get("eip").(string))
	params.Vxnet.Set(d.Get("vxnet").(string))
	params.PrivateIp.Set(d.Get("private_ip").(string))
	params.LoadbalancerType.Set(d.Get("type").(int))
	params.LoadbalancerName.Set(d.Get("name").(string))
	params.SecurityGroup.Set(d.Get("security_group").(string))
	resp, err := clt.CreateLoadBalancer(params)
	if err != nil {
		return err
	}
	d.SetId(resp.LoadbalancerId)
	_, err = LoadbalancerTransitionStateRefresh(clt, d.Id())
	if err != nil {
		return err
	}
	return resourceQingcloudLoadbalancerRead(d, meta)
}
func resourceQingcloudLoadbalancerRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	params := loadbalancer.DescribeLoadBalancersRequest{}
	params.LoadbalancersN.Add(d.Id())
	resp, err := clt.DescribeLoadBalancers(params)
	if err != nil {
		return err
	}
	if len(resp.LoadbalancerSet) == 0 {
		return errors.New("no load balancer")
	}
	lb := resp.LoadbalancerSet[0]
	d.Set("private_ip", lb.Vxnet.PrivateIP)
	return nil
}

func resourceQingcloudLoadbalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceQingcloudLoadbalancerDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	params := loadbalancer.StopLoadBalancersRequest{}
	params.LoadbalancersN.Add(d.Id())
	_, err := clt.StopLoadBalancers(params)
	if err != nil {
		return err
	}
	_, err = LoadbalancerTransitionStateRefresh(clt, d.Id())
	return err
}
