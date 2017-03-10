package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingcloudSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudSecurityGroupCreate,
		Read:   resourceQingcloudSecurityGroupRead,
		Update: resourceQingcloudSecurityGroupUpdate,
		Delete: resourceQingcloudSecurityGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "防火墙名称",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "防火墙介绍",
			},
		},
	}
}

func resourceQingcloudSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.CreateSecurityGroupInput)
	input.SecurityGroupName = qc.String(d.Get("name").(string))
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error create security group input validate: %s", err)
	}
	output, err := clt.CreateSecurityGroup(input)
	if err != nil {
		return fmt.Errorf("Error create security group: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error create security group: %s", *output.Message)
	}

	d.SetId(qc.StringValue(output.SecurityGroupID))
	err = modifySecurityGroupAttributes(d, meta, true)
	if err != nil {
		return err
	}
	return resourceQingcloudSecurityGroupRead(d, meta)
}

func resourceQingcloudSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.DescribeSecurityGroupsInput)
	input.SecurityGroups = []*string{qc.String(d.Id())}
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error describe security group input validate: %s", err)
	}
	output, err := clt.DescribeSecurityGroups(input)
	if err != nil {
		return fmt.Errorf("Error describe security group: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error describe security group: %s", *output.Message)
	}
	sg := output.SecurityGroupSet[0]
	d.Set("name", sg.SecurityGroupName)
	d.Set("description", sg.Description)
	return nil
}
func resourceQingcloudSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	err := modifySecurityGroupAttributes(d, meta, false)
	if err != nil {
		return err
	}
	return resourceQingcloudSecurityGroupRead(d, meta)
}

func resourceQingcloudSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	describeSecurityGroupInput := new(qc.DescribeSecurityGroupsInput)
	describeSecurityGroupInput.SecurityGroups = []*string{qc.String(d.Id())}
	describeSecurityGroupInput.Verbose = qc.Int(1)
	err := describeSecurityGroupInput.Validate()
	if err != nil {
		return fmt.Errorf("Error describe security group input validate: %s", err)
	}
	describeSecurityGroupOutput, err := clt.DescribeSecurityGroups(describeSecurityGroupInput)
	if err != nil {
		return fmt.Errorf("Error describe security group: %s", err)
	}
	if describeSecurityGroupOutput.RetCode != nil && qc.IntValue(describeSecurityGroupOutput.RetCode) != 0 {
		return fmt.Errorf("Error describe security group: %s", err)
	}
	if len(describeSecurityGroupOutput.SecurityGroupSet[0].Resources) > 0 {
		return fmt.Errorf("Error security group %s is using, can't delete", d.Id())
	}
	input := new(qc.DeleteSecurityGroupsInput)
	input.SecurityGroups = []*string{qc.String(d.Id())}
	err = input.Validate()
	if err != nil {
		return fmt.Errorf("Error describe security group input validate: %s", err)
	}
	output, err := clt.DeleteSecurityGroups(input)
	if err != nil {
		return fmt.Errorf("Error delete security group: %s", err)
	}
	if output.RetCode != nil && qc.IntValue(output.RetCode) != 0 {
		return fmt.Errorf("Error delete security group: %s", *output.Message)
	}
	d.SetId("")
	return nil
}
