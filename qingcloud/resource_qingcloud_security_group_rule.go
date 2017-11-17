package qingcloud

import (
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
				Description:  "protocol",
				ValidateFunc: withinArrayString("tcp", "udp", "icmp", "gre", "esp", "ah", "ipip"),
			},
			"priority": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: withinArrayIntRange(0, 100),
				Default:      0,
				Description:  "priority,From high to low 0 - 100",
			},
			"action": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: withinArrayString("accept", "drop"),
			},
			"direction": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "direction,0 express down ,1 express up.default 0 .",
				ValidateFunc: withinArrayInt(0, 1),
				Default:      0,
			},
			"from_port": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "if protocol is tcp or udp,this value is start port. else if protocol is icmp,this value is the type of ICMP. the others protocol don't need this value.",
				ValidateFunc: validatePortString,
			},
			"to_port": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validatePortString,
				Description:  "if protocol is tcp or udp,this value is end port. else if protocol is icmp,this value is the code of ICMP. the others protocol don't need this value.",
			},
			"cidr_block": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNetworkCIDR,
				Description:  "target IP,the Security Group Rule only affect to those IPs .",
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
	rule.Direction = qc.Int(d.Get("direction").(int))
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
	var output *qc.AddSecurityGroupRulesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.AddSecurityGroupRules(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	d.SetId(qc.StringValue(output.SecurityGroupRules[0]))
	// Lock security group resource for apply security group
	if err := applySecurityGroupRule(d, meta); err != nil {
		return nil
	}
	return resourceQingcloudSecurityGroupRuleRead(d, meta)
}

func resourceQingcloudSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.DescribeSecurityGroupRulesInput)
	input.SecurityGroupRules = []*string{qc.String(d.Id())}
	var output *qc.DescribeSecurityGroupRulesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DescribeSecurityGroupRules(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if len(output.SecurityGroupRuleSet) == 0 {
		d.SetId("")
		return nil
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
	if err := ModifySecurityGroupRuleAttributes(d, meta); err != nil {
		return err
	}
	if err := applySecurityGroupRule(d, meta); err != nil {
		return err
	}
	return resourceQingcloudSecurityGroupRuleRead(d, meta)
}

func resourceQingcloudSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.DeleteSecurityGroupRulesInput)
	input.SecurityGroupRules = []*string{qc.String(d.Id())}
	var output *qc.DeleteSecurityGroupRulesOutput
	var err error
	simpleRetry(func() error {
		output, err = clt.DeleteSecurityGroupRules(input)
		return isServerBusy(err)
	})
	if err != nil {
		return err
	}
	if err := applySecurityGroupRule(d, meta); err != nil {
		return nil
	}
	d.SetId("")
	return nil
}
