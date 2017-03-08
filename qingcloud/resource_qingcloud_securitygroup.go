package qingcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/magicshui/qingcloud-go/securitygroup"
	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func resourceQingcloudSecuritygroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudSecuritygroupCreate,
		Read:   resourceQingcloudSecuritygroupRead,
		Update: resourceQingcloudSecuritygroupUpdate,
		Delete: resourceQingcloudSecuritygroupDelete,
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

func resourceQingcloudSecuritygroupCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.CreateSecurityGroupInput)
	input.SecurityGroupName = qc.String(d.Get("name").(string))
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error create securitygroup input validate: %s", err)
	}
	output, err := clt.CreateSecurityGroup(input)
	if err != nil {
		return fmt.Errorf("Error create securitygroup: %s", err)
	}
	if output.RetCode != 0 {
		return fmt.Errorf("Error create securitygroup: %s", output.Message)
	}

	d.SetId(qc.StringValue(output.SecurityGroupID))

	description := d.Get("description").(string)
	if description != "" {
		modifyAtrributes := securitygroup.ModifySecurityGroupAttributesRequest{}
		modifyAtrributes.SecurityGroup.Set(resp.SecurityGroupId)
		modifyAtrributes.Description.Set(description)
		_, err := clt.ModifySecurityGroupAttributes(modifyAtrributes)
		if err != nil {
			// 这里可以不用返回错误
			return fmt.Errorf("Error modify security description: %s", err)
		}
	}
	d.SetId(resp.SecurityGroupId)
	return resourceQingcloudSecuritygroupRead(d, meta)
}

func resourceQingcloudSecuritygroupRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	input := new(qc.DescribeSecurityGroupsInput)
	input.SecurityGroups = []*string{qc.String(d.Id())}
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("Error describe securitygroup input validate: %s", err)
	}
	output, err := clt.DescribeSecurityGroups(input)
	if err != nil {
		return fmt.Errorf("Error describe securitygroup: %s", err)
	}
	if output.RetCode != 0 {
		return fmt.Errorf("Error describe securitygroup: %s", output.Message)
	}
	sg := output.SecurityGroupSet[0]
	d.Set("name", sg.SecurityGroupName)
	d.Set("description", sg.Description)
	return nil
}
func resourceQingcloudSecuritygroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return modifySecurityGroupAttributes(d, meta, false)
}

func resourceQingcloudSecuritygroupDelete(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup
	describeSecurityGroupInput := new(qc.DescribeSecurityGroupsInput)
	describeSecurityGroupInput.SecurityGroups = []*string{qc.String(d.Id())}
	describeSecurityGroupInput.Verbose = 1
	err := describeSecurityGroupInput.Validate()
	if err != nil {
		return fmt.Errorf("Error describe securitygroup input validate: %s", err)
	}
	describeSecurityGroupOutput, err := clt.DescribeSecurityGroups(describeSecurityGroupInput)
	if err != nil {
		return fmt.Errorf("Error describe securitygroup: %s", err)
	}
	if describeSecurityGroupOutput.RetCode != 0 {
		return fmt.Errorf("Error describe securitygroup: %s", err)
	}
	if len(describeSecurityGroupOutput.SecurityGroupSet[0].Resources) > 0 {
		return fmt.Errorf("Error securitygroup %s is using, can't delete", d.Id())
	}
	input := new(qc.DeleteSecurityGroupsInput)
	input.SecurityGroups = []*string{qc.String(d.Id())}
	err = input.Validate()
	if err != nil {
		return fmt.Errorf("Error describe securitygroup input validate: %s", err)
	}
	output, err := clt.DeleteSecurityGroups(input)
	if err != nil {
		return fmt.Errorf("Error delete securitygroup: %s", err)
	}
	if output.RetCode != 0 {
		return fmt.Errorf("Error delete securitygroup: %s", output.Message)
	}
	d.SetId("")
	return nil
}
