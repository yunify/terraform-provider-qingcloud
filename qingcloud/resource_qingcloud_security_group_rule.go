package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingcloudSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudSecurityGroupRuleCreate,
		Read:   resourceQingcloudSecurityGroupRuleRead,
		Update: resourceQingcloudSecurityGroupRuleUpdate,
		Delete: resourceQingcloudSecurityGroupRuleDelete,
		Schema: map[string]*schema.Schema{
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "协议",
				ValidateFunc: withinArrayString("tcp", "udp", "icmp", "gre", "esp", "ah", "ipip"),
			},
			"priority": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: withinArrayIntRange(0, 100),
				Default:      0,
				Description:  "优先级，由高到低为 0 - 100",
			},
			"action": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: withinArrayString("accept", "drop"),
			},
			"direction": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "方向，0 表示下行，1 表示上行。默认为 0。",
				ValidateFunc: withinArrayInt(0, 1),
				Default:      0,
			},
			"from_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "如果协议为 tcp 或 udp，此值表示起始端口。 如果协议为 icmp，此值表示 ICMP 类型。 其他协议无需此值。	",
			},
			"to_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "如果协议为 tcp 或 udp，此值表示结束端口。 如果协议为 icmp，此值表示 ICMP 代码。 其他协议无需此值。	",
			},
			"cidr_block": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNetworkCIDR,
				Description: "目标 IP，如果填写，则这条防火墙规则只对此IP（或IP段）有效。	",
			},
		},
	}
}

func resourceQingcloudSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.AddSecurityGroupRulesInput)
	sgID := d.Get("security_group_id").(string)
	input.SecurityGroup = qc.String(sgID)
	rule := new(qc.SecurityGroupRule)
	rule.Priority = qc.Int(d.Get("priority").(int))
	rule.Protocol = qc.String(d.Get("protocol").(string))
	rule.Action = qc.String(d.Get("action").(string))
	if d.Get("name").(string) != "" {
		rule.SecurityGroupRuleName = qc.String(d.Get("name").(string))
	}
	if d.Get("from_port").(string) != "" {
		rule.Val1 = qc.String(d.Get("from_port").(string))
	}
	if d.Get("to_port").(string) != "" {
		rule.Val2 = qc.String(d.Get("to_port").(string))
	}
	if d.Get("cidr_block").(string) != "" {
		rule.Val3 = qc.String(d.Get("cidr_block").(string))
	}
	input.Rules = []*qc.SecurityGroupRule{rule}
	output, err := clt.AddSecurityGroupRules(input)
	if err != nil {
		return fmt.Errorf("Error add security group rule: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error add security group rule: %s", err)
	}
	d.SetId(qc.StringValue(output.SecurityGroupRules[0]))
	// Lock security group resource for apply security group
	err = applySecurityGroupRule(d, meta)
	if err != nil {
		return err
	}
	return resourceQingcloudSecurityGroupRuleRead(d, meta)
}

func resourceQingcloudSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.DescribeSecurityGroupRulesInput)
	input.SecurityGroup = qc.String(d.Get("security_group_id").(string))
	input.SecurityGroupRules = []*string{qc.String(d.Id())}
	output, err := clt.DescribeSecurityGroupRules(input)
	if err != nil {
		return err
	}
	sgRule := output.SecurityGroupRuleSet[0]
	d.Set("security_group_id", qc.StringValue(sgRule.SecurityGroupID))
	d.Set("protocol", qc.StringValue(sgRule.Protocol))
	d.Set("priority", qc.IntValue(sgRule.Priority))
	d.Set("action", qc.StringValue(sgRule.Action))
	d.Set("from_port", qc.StringValue(sgRule.Val1))
	d.Set("to_port", qc.StringValue(sgRule.Val2))
	d.Set("cidr_block", qc.StringValue(sgRule.Val3))
	d.Set("name", qc.StringValue(sgRule.SecurityGroupRuleName))
	return nil
}

func resourceQingcloudSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {

	err := ModifySecurityGroupRuleAttributes(d, meta, false)
	if err != nil {
		return err
	}
	err = applySecurityGroupRule(d, meta)
	if err != nil {
		return err
	}
	return resourceQingcloudSecurityGroupRuleRead(d, meta)
}

func resourceQingcloudSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.DeleteSecurityGroupRulesInput)
	input.SecurityGroupRules = []*string{qc.String(d.Id())}
	_, err := clt.DeleteSecurityGroupRules(input)
	if err != nil {
		return err
	}
	err = applySecurityGroupRule(d, meta)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
