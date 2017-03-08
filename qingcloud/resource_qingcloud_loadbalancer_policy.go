package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/loadbalancer"
	"log"
)

func resourceQingcloudLoadbalancerPloicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudLoadbalancerPloicyCreate,
		Read:   resourceQingcloudLoadbalancerPloicyRead,
		Update: resourceQingcloudLoadbalancerPloicyUpdate,
		Delete: resourceQingcloudLoadbalancerPloicyDelete,
		Schema: map[string]*schema.Schema{
			"operator": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "转发策略规则间的逻辑关系：”and” 是『与』，”or” 是『或』，默认是 “or”",
				ValidateFunc: withinArrayString("and", "or"),
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
	log.Printf("Finish loadbalancer policy %s", resp.LoadbalancerPolicyId)
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
