package qingcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magicshui/qingcloud-go/securitygroup"
)

func resourceQingcloudSecuritygroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudSecuritygroupRuleCreate,
		Read:   resourceQingcloudSecuritygroupRuleRead,
		Update: resourceQingcloudSecuritygroupRuleUpdate,
		Delete: resourceQingcloudSecuritygroupRuleDelete,
		Schema: map[string]*schema.Schema{
			"securitygroup": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"action": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"direction": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"val1": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"val2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"val3": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceQingcloudSecuritygroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	params := securitygroup.AddSecurityGroupRulesRequest{}
	params.SecurityGroup.Set(d.Get("securitygroup").(string))
	params.RulesNProtocol.Add(d.Get("protocol").(string))
	params.RulesNPriority.Add(int64(d.Get("priority").(int)))
	params.RuleNSecurityGroupRuleName.Add(d.Get("name").(string))
	params.RulesNAction.Add(d.Get("action").(string))
	params.RulesNDirection.Add(int64(d.Get("direction").(int)))
	params.RulesNVal1.Add(d.Get("val1").(string))
	params.RulesNVal2.Add(d.Get("val2").(string))
	params.RulesNVal3.Add(d.Get("val3").(string))
	resp, err := clt.AddSecurityGroupRules(params)
	if err != nil {
		return err
	}
	d.SetId(resp.SecurityGroupRules[0])
	return nil
}
func resourceQingcloudSecuritygroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	params := securitygroup.DescribeSecurityGroupRulesRequest{}
	params.SecurityGroupRuleN.Add(d.Id())
	resp, err := clt.DescribeSecurityGroupRules(params)
	if err != nil {
		return err
	}
	sr := resp.SecurityGroupRuleSet[0]
	d.Set("protocol", sr.Protocol)
	d.Set("priority", sr.Priority)
	d.Set("action", sr.Action)
	d.Set("direction", sr.Direction)
	d.Set("val1", sr.Val1)
	d.Set("val2", sr.Val2)
	d.Set("val3", sr.Val3)
	return nil
}

func resourceQingcloudSecuritygroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceQingcloudSecuritygroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	params := securitygroup.DeleteSecurityGroupRulesRequest{}
	params.SecurityGroupRulesN.Add(d.Id())
	_, err := clt.DeleteSecurityGroupRules(params)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
