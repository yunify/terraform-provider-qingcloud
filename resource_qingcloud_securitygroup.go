package qingcloud

import (
	"fmt"
	"log"

	// "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/magicshui/qingcloud-go/securitygroup"
)

func resourceQingcloudSecuritygroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceQingcloudSecuritygroupCreate,
		Read:   resourceQingcloudSecuritygroupRead,
		Update: resourceQingcloudSecuritygroupUpdate,
		Delete: nil,
		Schema: resourceQingcloudSecuritygroupSchema(),
	}
}

func resourceQingcloudSecuritygroupCreate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup

	// TODO: 这个地方以后需要判断错误
	securityGroupName := d.Get("name").(string)

	// 开始创建 ssh 密钥
	params := securitygroup.CreateSecurityGroupRequest{}
	params.SecurityGroupName.Set(securityGroupName)

	resp, err := clt.CreateSecurityGroup(params)
	if err != nil {
		return fmt.Errorf("Error create security group", err)
	}

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
	return nil
}

func resourceQingcloudSecuritygroupRead(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup

	// 设置请求参数
	params := securitygroup.DescribeSecurityGroupsRequest{}
	params.SecurityGroupsN.Add(d.Id())
	params.Verbose.Set(1)

	resp, err := clt.DescribeSecurityGroups(params)
	if err != nil {
		return fmt.Errorf("Error retrieving Keypair: %s", err)
	}
	for _, sg := range resp.SecurityGroupSet {
		log.Printf("Current Security Group is :%s", sg.SecurityGroupID)
		if sg.SecurityGroupID == d.Id() {
			log.Printf("Get security Group %#v", sg)
			d.Set("name", sg.SecurityGroupName)
			d.Set("description", sg.Description)
			return nil
		}
	}
	log.Printf("Unable to find security group %#v within: %#v", d.Id(), resp.SecurityGroupSet)
	d.SetId("")
	return nil
}

func resourceQingcloudSecuritygroupDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceQingcloudSecuritygroupUpdate(d *schema.ResourceData, meta interface{}) error {
	clt := meta.(*QingCloudClient).securitygroup

	if !d.HasChange("name") && !d.HasChange("description") {
		return nil
	}

	params := securitygroup.ModifySecurityGroupAttributesRequest{}
	if d.HasChange("description") {
		params.Description.Set(d.Get("description").(string))
	}
	if d.HasChange("name") {
		params.SecurityGroupName.Set(d.Get("name").(string))
	}
	params.SecurityGroup.Set(d.Id())
	log.Println("--------------  s1")
	_, err := clt.ModifySecurityGroupAttributes(params)
	if err != nil {
		return fmt.Errorf("Error modify security group %s", d.Id())
	}
	return nil
}

func resourceQingcloudSecuritygroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"id": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}
