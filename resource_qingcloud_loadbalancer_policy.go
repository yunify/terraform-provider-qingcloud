package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/loadbalancer"
)

func resourceQingcloudLoadbalancerPloicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudLoadbalancerPloicyCreate,
		Read:   resourceQingcloudLoadbalancerPloicyRead,
		Update: nil,
		Delete: resourceQingcloudLoadbalancerPloicyDelete,
		Schema: map[string]*schema.Schema{
			"operator": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceQingcloudLoadbalancerPloicyCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	params := loadbalancer.CreateLoadBalancerPolicyRequest{}
	params.Operator.Set(d.Get("operator").(string))
	resp, err := clt.CreateLoadBalancerPolicy(params)
	if err != nil {
		return err
	}
	d.SetId(resp.LoadbalancerPolicyId)
	return nil
}

func resourceQingcloudLoadbalancerPloicyRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	params := loadbalancer.DescribeLoadBalancerPoliciesRequest{}
	params.LoadbalancerPoliciesN.Add(d.Id())
	_, err := clt.DescribeLoadBalancerPolicies(params)
	if err != nil {
		return err
	}
	return nil
}

func resourceQingcloudLoadbalancerPloicyUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceQingcloudLoadbalancerPloicyDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).loadbalancer
	params := loadbalancer.DeleteLoadBalancerPoliciesRequest{}
	params.LoadbalancerPoliciesN.Add(d.Id())
	_, err := clt.DeleteLoadBalancerPolicies(params)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
